package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/spurf/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	check := fs.Bool("c", false, "Healthcheck (check and exit)")

	dbConfig := db.Flags(fs, "db")
	enedisConfig := enedis.Flags(fs, "enedis")

	logger.Fatal(fs.Parse(os.Args[1:]))

	if *check {
		return
	}

	spurfDb, err := db.New(dbConfig)
	logger.Fatal(err)

	enedisApp, err := enedis.New(enedisConfig, spurfDb)
	logger.Fatal(err)

	enedisApp.Start()
}