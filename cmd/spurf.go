package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/opentracing"
	"github.com/ViBiOh/spurf/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	opentracingConfig := opentracing.Flags(fs, "tracing")
	dbConfig := db.Flags(fs, "db")
	enedisConfig := enedis.Flags(fs, "enedis")

	logger.Fatal(fs.Parse(os.Args[1:]))

	if *check {
		return
	}

	opentracing.New(opentracingConfig)

	spurfDb, err := db.New(dbConfig)
	logger.Fatal(err)

	enedisApp, err := enedis.New(enedisConfig, spurfDb)
	logger.Fatal(err)

	enedisApp.Start()
}
