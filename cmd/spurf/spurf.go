package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/spurf/v2/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	dbConfig := db.Flags(fs, "db")
	enedisConfig := enedis.Flags(fs, "enedis")
	loggerConfig := logger.Flags(fs, "logger")

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	spurfDb, err := db.New(dbConfig)
	logger.Fatal(err)
	defer spurfDb.Close()

	enedis.New(enedisConfig, spurfDb).Start()
}
