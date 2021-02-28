package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Scarlet-Fairy/cobold/pkg/build"
	"github.com/Scarlet-Fairy/cobold/pkg/build/docker"
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	"github.com/Scarlet-Fairy/cobold/pkg/clone/git"
	"github.com/Scarlet-Fairy/cobold/pkg/log"
	"github.com/Scarlet-Fairy/cobold/pkg/tracing"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/xid"
	"io/ioutil"
	"os"
)

var (
	jobID          = flag.String("build-id", xid.New().String(), "build ID that identify the actual job")
	gitRepository  = flag.String("git-repo", "https://github.com/buildkite/nodejs-docker-example", "repository to clone")
	dockerEndpoint = flag.String("docker-endpoint", "localhost:2375", "docker daemon endpoint")
	dockerRegistry = flag.String("docker-registry", "localhost:5000", "docker registry to push image")
	// redisUrl      = flag.String("redis-url", "localhost", "redis url where publish complete events")
	// redisUser     = flag.String("redis-user", "", "user used to authenticate to redis")
	// redisPassword = flag.String("redis-password", "", "password used to authenticate to redis")
)

var ctx = context.Background()

func main() {
	flag.Parse()

	logger := log.InitLogger(*jobID)

	tracer, closer, err := tracing.Init("cobold")
	if err != nil {
		level.Error(logger).Log("error", err.Error())
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

	cloneInstance := git.MakeClone(*jobID, logger, tracer)
	cloneOptions := clone.Options{
		Url:  *gitRepository,
		Path: tmpDir + "/",
	}

	buildInstance := docker.MakeBuild(*jobID, *dockerEndpoint, logger, tracer)
	buildOptions := build.Options{
		Tag:       fmt.Sprintf("cobold/%s", *jobID),
		Directory: tmpDir,
	}

	if err := cloneInstance.Clone(ctx, cloneOptions); err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}

	buildOutputStream, err := buildInstance.Build(ctx, buildOptions)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}

	logger.Log("stream", buildOutputStream)
}
