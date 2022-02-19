package main

import (
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

	tracerApp, err := tracer.New(tracerConfig)
	logger.Fatal(err)
	defer tracerApp.Close()

	spurfDb, err := db.New(dbConfig, tracerApp)
	logger.Fatal(err)
	defer spurfDb.Close()

	enedis.New(enedisConfig, spurfDb).Start()
}
