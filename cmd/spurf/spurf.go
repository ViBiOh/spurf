package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/server"
	"github.com/ViBiOh/spurf/v2/pkg/datahub"
	"github.com/ViBiOh/spurf/v2/pkg/enedis"
)

func main() {
	fs := flag.NewFlagSet("spurf", flag.ExitOnError)

	promServerConfig := server.Flags(fs, "prometheus", flags.NewOverride("Port", 9090), flags.NewOverride("IdleTimeout", "10s"), flags.NewOverride("ShutdownTimeout", "5s"))
	healthConfig := health.Flags(fs, "")

	alcotestConfig := alcotest.Flags(fs, "")
	loggerConfig := logger.Flags(fs, "logger")
	prometheusConfig := prometheus.Flags(fs, "prometheus")

	dbConfig := db.Flags(fs, "db")
	datahubConfig := datahub.Flags(fs, "datahub")
	enedisConfig := enedis.Flags(fs, "enedis")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)
	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	promServer := server.New(promServerConfig)
	prometheusApp := prometheus.New(prometheusConfig)

	spurfDb, err := db.New(dbConfig)
	logger.Fatal(err)

	healthApp := health.New(healthConfig, spurfDb.Ping)

	enedisApp, err := enedis.New(enedisConfig, spurfDb, datahub.New(datahubConfig))
	logger.Fatal(err)

	enedisApp.Start(healthApp.Done())

	go promServer.Start("prometheus", healthApp.End(), prometheusApp.Handler())

	healthApp.WaitForTermination(promServer.Done())
	server.GracefulWait(promServer.Done())
}
