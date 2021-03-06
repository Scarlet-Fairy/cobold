package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	dockerBuild "github.com/Scarlet-Fairy/cobold/pkg/build/docker"
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	"github.com/Scarlet-Fairy/cobold/pkg/clone/git"
	"github.com/Scarlet-Fairy/cobold/pkg/log"
	"github.com/Scarlet-Fairy/cobold/pkg/notify"
	amqpNotify "github.com/Scarlet-Fairy/cobold/pkg/notify/amqp"
	"github.com/Scarlet-Fairy/cobold/pkg/push"
	dockerPush "github.com/Scarlet-Fairy/cobold/pkg/push/docker"
	otelTracing "github.com/Scarlet-Fairy/cobold/pkg/tracing/otel"
	dockerAPI "github.com/fsouza/go-dockerclient"
	goKitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	pushgateway "github.com/prometheus/client_golang/prometheus/push"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"io/ioutil"
	"os"
	"time"
)

var (
	jobID          = flag.String("job-id", "1", "Job's ID that identify the actual job")
	gitRepository  = flag.String("git-repo", "https://github.com/buildkite/nodejs-docker-example", "repository to clone")
	dockerUrl      = flag.String("docker-url", "unix:///var/run/docker.sock", "docker daemon endpoint")
	dockerRegistry = flag.String("docker-registry", "localhost:5000", "docker registry to push image")
	tracingHost    = flag.String("tracing-host", "localhost", "host where send traces")
	tracingPort    = flag.String("tracing-port", "6831", "port of the host where send traces")
	rabbitMQUrl    = flag.String("rabbitmq-url", "amqp://guest:guest@localhost:5672/", "rabbitmq url where publish complete message")
	pushGatewayUrl = flag.String("pushgateway-url", "http://localhost:9091", "pushgateway url where metrics are shipped")
)

var ctx = context.Background()

func main() {
	flag.Parse()
	imageName := fmt.Sprintf("%s/cobold/%s", *dockerRegistry, *jobID)

	logger, cloneLogger, buildLogger, pushLogger, notifyLogger := log.InitLogger(*jobID)

	if _, ok := os.LookupEnv("DEV"); ok {
		time.Sleep(time.Hour * 10)
	}

	flush, err := otelTracing.InitTraceProvider(false, "cobold", *jobID, *tracingHost, *tracingPort)
	if err != nil {
		level.Error(logger).Log("msg", "Tracing init failed", "error", err.Error())
		os.Exit(1)
	}
	defer flush()

	var cloneDuration, buildDuration, pushDuration, notifyDuration metrics.Histogram
	{
		cloneDuration = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: "scarlet_fairy",
			Subsystem: "cobold",
			Name:      "clone_duration",
			Help:      "",
		}, []string{"jobID"})
		buildDuration = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: "scarlet_fairy",
			Subsystem: "cobold",
			Name:      "build_duration",
			Help:      "",
		}, []string{"jobID"})
		pushDuration = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: "scarlet_fairy",
			Subsystem: "cobold",
			Name:      "push_duration",
			Help:      "",
		}, []string{"jobID"})
		notifyDuration = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: "scarlet_fairy",
			Subsystem: "cobold",
			Name:      "notify_duration",
			Help:      "",
		}, []string{"jobID"})
	}

	tr := otel.Tracer("cobold")
	ctx, span := tr.Start(ctx, "job")
	logger.Log("traceID", span.SpanContext().TraceID)
	defer span.End()

	tmpDir, err := ioutil.TempDir("", "clone-")
	if err != nil {
		level.Error(logger).Log("msg", "temp dir creation failed", "error", err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	level.Debug(logger).Log("dir", tmpDir)

	dockerClient, err := dockerAPI.NewClient(*dockerUrl)
	if err != nil {
		level.Error(logger).Log("docker-endpoint", *dockerUrl, "msg", "docker client cannot be created", "error", err)
		os.Exit(1)
	}
	defer dockerClient.RemoveImage(imageName) // debug

	rabbitMQClient, err := newRabbitMQClient(*rabbitMQUrl)
	if err != nil {
		level.Error(logger).Log("rabbitmq-endpoint", *rabbitMQUrl, "msg", "rabbitmq client cannot be created", "error", err)
		os.Exit(1)
	}

	rabbitMQChannel, err := rabbitMQClient.Channel()
	if err != nil {
		level.Error(logger).Log("rabbitmq-endpoint", *rabbitMQUrl, "msg", "rabbitmq channel cannot be created", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := rabbitMQChannel.Close(); err != nil {
			panic(err)
		}
	}()

	cloneInstance := git.MakeClone(*jobID, cloneDuration, cloneLogger, tr)
	cloneOptions := clone.Options{
		Url:  *gitRepository,
		Path: tmpDir + "/",
	}

	buildInstance := dockerBuild.MakeBuild(*jobID, dockerClient, buildDuration, buildLogger, tr)
	buildOptions := build.Options{
		Name:      imageName,
		Directory: tmpDir,
	}

	pushInstance := dockerPush.MakePush(*jobID, dockerClient, pushDuration, pushLogger, tr)
	pushOptions := push.Options{
		Name:     imageName,
		Tag:      "latest",
		Registry: *dockerRegistry,
	}

	notifyInstance := amqpNotify.MakeNotify(rabbitMQChannel, *jobID, notifyDuration, notifyLogger, tr)

	if err := notifyInstance.Init(ctx); err != nil {
		level.Error(logger).Log("msg", "Notify Init failed", "err", err)
		os.Exit(1)
	}

	err = cloneInstance.Clone(ctx, cloneOptions)
	handleStepError(ctx, notifyInstance, notifyLogger, cloneLogger, err, clone.StepName)

	buildOutputStream, err := buildInstance.Build(ctx, buildOptions)
	handleStepError(ctx, notifyInstance, notifyLogger, cloneLogger, err, build.StepName)
	buildLogger.Log("stream", buildOutputStream)

	err = pushInstance.Push(ctx, pushOptions)
	handleStepError(ctx, notifyInstance, notifyLogger, cloneLogger, err, push.StepName)

	if err = pushgateway.New(*pushGatewayUrl, "cobold").Gatherer(stdprometheus.DefaultGatherer).Push(); err != nil {
		level.Error(logger).Log("pushgateway-url", *pushGatewayUrl, "err", err)
	}
}

func newRabbitMQClient(uri string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func handleStepError(
	ctx context.Context,
	notifyInstance notify.Notify,
	notifyLogger goKitLog.Logger,
	stepLogger goKitLog.Logger,
	err error,
	step string,
) {
	if err != nil {
		level.Error(stepLogger).Log("msg", "failed", "error", err.Error())

		if err := notifyInstance.NotifyCompletion(ctx, notify.Options{
			Err:    err,
			Reason: step,
			JobID:  *jobID,
		}); err != nil {
			level.Error(notifyLogger).Log("msg", "failed", "error", err.Error())
		}

		os.Exit(1)
	}

	if err := notifyInstance.NotifyCompletion(ctx, notify.Options{
		Err:    nil,
		Reason: step,
		JobID:  *jobID,
	}); err != nil {
		level.Error(notifyLogger).Log("msg", "failed", "error", err.Error())

		os.Exit(1)
	}

}
