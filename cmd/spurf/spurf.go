package main

import (
	"context"
	"flag"
	"os"

	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/tracer"
	"github.com/ViBiOh/spurf/v2/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")
	tracerConfig := tracer.Flags(fs, "tracer")

	dbConfig := db.Flags(fs, "db")
	enedisConfig := enedis.Flags(fs, "enedis")

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	ctx := context.Background()

	tracerApp, err := tracer.New(ctx, tracerConfig)
	logger.Fatal(err)
	defer tracerApp.Close(ctx)

	spurfDb, err := db.New(ctx, dbConfig, tracerApp.GetTracer("database"))
	logger.Fatal(err)
	defer spurfDb.Close()

	enedis.New(enedisConfig, spurfDb).Start()
}
