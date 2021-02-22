package main

import (
	"context"
	"flag"
	"github.com/Scarlet-Fairy/cobold/pkg/clone/git"
	"github.com/Scarlet-Fairy/cobold/pkg/log"
	"github.com/Scarlet-Fairy/cobold/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"io/ioutil"
	"os"
)

var (
	buildId       = flag.String("build-id", "", "build ID that identify the actual job")
	gitRepository = flag.String("git-repo", "https://github.com/ArcaneDiver/my-dynamical-ip", "repository to clone")
	// dockerRegistry = flag.String("docker-registry", "", "docker registry to push image")
	// redisUrl      = flag.String("redis-url", "localhost", "redis url where publish complete events")
	// redisUser     = flag.String("redis-user", "", "user used to authenticate to redis")
	// redisPassword = flag.String("redis-password", "", "password used to authenticate to redis")
)

var ctx = context.Background()

func main() {
	flag.Parse()

	tracer, closer := tracing.Init("cobold")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, "cobold")
	defer span.Finish()

	tmpDir, err := ioutil.TempDir("", "clone")
	if err != nil {
		log.Logger.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	log.Logger.Debug(tmpDir)

	clone := git.MakeClone(tracer)

	if err := clone.Clone(ctx, *gitRepository, tmpDir); err != nil {
		log.Logger.Fatal(err)
	}
}
