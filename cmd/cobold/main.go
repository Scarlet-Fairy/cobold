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
	"github.com/Scarlet-Fairy/cobold/pkg/push"
	dockerPush "github.com/Scarlet-Fairy/cobold/pkg/push/docker"
	"github.com/Scarlet-Fairy/cobold/pkg/tracing"
	dockerAPI "github.com/fsouza/go-dockerclient"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/xid"
	"io/ioutil"
	"os"
)

var (
	jobID          = flag.String("job-id", xid.New().String(), "Job's ID that identify the actual job")
	gitRepository  = flag.String("git-repo", "https://github.com/buildkite/nodejs-docker-example", "repository to clone")
	dockerUrl      = flag.String("docker-url", "localhost:2375", "docker daemon endpoint")
	dockerRegistry = flag.String("docker-registry", "localhost:5000", "docker registry to push image")

// redisUrl      = flag.String("redis-url", "localhost", "redis url where publish complete events")
)

var ctx = context.Background()

func main() {
	flag.Parse()
	imageName := fmt.Sprintf("%s/cobold/%s", *dockerRegistry, *jobID)

	logger, cloneLogger, buildLogger, pushLogger := log.InitLogger(*jobID)

	tracer, closer, err := tracing.Init("cobold")
	if err != nil {
		level.Error(logger).Log("msg", "Tracing init failed", "error", err.Error())
		os.Exit(1)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, "cobold")
	defer span.Finish()

	tmpDir, err := ioutil.TempDir("", "clone-")
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	level.Debug(logger).Log("dir", tmpDir)

	dockerClient, err := dockerAPI.NewClient(*dockerUrl)
	if err != nil {
		level.Error(logger).Log("docker-endpoint", dockerUrl, "msg", "docker client cannot be created", "error", err)
		os.Exit(1)
	}
	defer dockerClient.RemoveImage(imageName)

	cloneInstance := git.MakeClone(*jobID, cloneLogger, tracer)
	cloneOptions := clone.Options{
		Url:  *gitRepository,
		Path: tmpDir + "/",
	}

	buildInstance := dockerBuild.MakeBuild(*jobID, dockerClient, buildLogger, tracer)
	buildOptions := build.Options{
		Name:      imageName,
		Directory: tmpDir,
	}

	pushInstance := dockerPush.MakePush(*jobID, dockerClient, pushLogger, tracer)
	pushOptions := push.Options{
		Name:     imageName,
		Tag:      "latest",
		Registry: *dockerRegistry,
	}

	if err := cloneInstance.Clone(ctx, cloneOptions); err != nil {
		level.Error(logger).Log("msg", "Clone failed", "error", err.Error())
		os.Exit(1)
	}

	buildOutputStream, err := buildInstance.Build(ctx, buildOptions)
	if err != nil {
		level.Error(logger).Log("msg", "Build failed", "error", err.Error())
		os.Exit(1)
	}
	buildLogger.Log("stream", buildOutputStream)

	if err := pushInstance.Push(ctx, pushOptions); err != nil {
		level.Error(logger).Log("msg", "Push failed", "error", err.Error())
		os.Exit(1)
	}
}
