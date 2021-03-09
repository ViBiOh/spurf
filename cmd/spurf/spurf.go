package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/spurf/v2/pkg/datahub"
	"github.com/ViBiOh/spurf/v2/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")

	dbConfig := db.Flags(fs, "db")
	datahubConfig := datahub.Flags(fs, "datahub")
	enedisConfig := enedis.Flags(fs, "enedis")

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	spurfDb, err := db.New(dbConfig)
	logger.Fatal(err)

	enedisApp, err := enedis.New(enedisConfig, spurfDb, datahub.New(datahubConfig))
	logger.Fatal(err)

	enedisApp.Start()
}
