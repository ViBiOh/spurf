package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/scheduler"
	"github.com/ViBiOh/spurf/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	opentracingConfig := opentracing.Flags(fs, "tracing")
	dbConfig := db.Flags(fs, "db")
	schedulerConfig := scheduler.Flags(fs, "scheduler")
	enedisConfig := enedis.Flags(fs, "enedis")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatal("%#v", err)
	}

	if *check {
		return
	}

	opentracing.New(opentracingConfig)

	spurfDb, err := db.New(dbConfig)
	if err != nil {
		logger.Fatal("%#v", err)
	}

	enedisApp, err := enedis.New(enedisConfig, spurfDb)
	if err != nil {
		logger.Fatal("%#v", err)
	}

	scheduler, err := scheduler.New(schedulerConfig, enedisApp)
	if err != nil {
		logger.Fatal("%#v", err)
	}

	if err := enedisApp.Start(); err != nil {
		logger.Error("%#v", err)
	}

	scheduler.Start()
}
