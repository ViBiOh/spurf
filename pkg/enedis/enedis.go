package enedis

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/cron"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/spurf/v2/pkg/datahub"
)

// App of package
type App interface {
	Start(<-chan struct{})
}

// Config of package
type Config struct {
	name *string
	tz   *string
}

type app struct {
	db         *sql.DB
	datahubApp datahub.App
	tz         *time.Location

	name string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		name: flags.New(prefix, "enedis").Name("Name").Default("home").Label("Name").ToString(fs),
		tz:   flags.New(prefix, "enedis").Name("Timezone").Default("Europe/Paris").Label("Timezone").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, db *sql.DB, datahubApp datahub.App) (App, error) {
	timezone, err := time.LoadLocation(strings.TrimSpace(*config.tz))
	if err != nil {
		return nil, fmt.Errorf("unable to load location: %s", err)
	}

	return &app{
		db:   db,
		name: strings.TrimSpace(*config.name),
		tz:   timezone,

		datahubApp: datahubApp,
	}, nil
}

// Start the package
func (a app) Start(done <-chan struct{}) {
	feeder := cron.New().Days().In(a.tz.String()).OnError(func(err error) {
		logger.Error("unable to feed spurf: %s", err)
	}).OnSignal(syscall.SIGUSR1)

	logger.Info("Feeding spurf database %s", feeder)

	feeder.Start(func(now time.Time) error {
		ctx := context.Background()

		lastInsert, err := a.getLastFetch(ctx)
		if err != nil {
			return fmt.Errorf("unable to get last fetch: %s", err)
		}

		consumption, err := a.datahubApp.GetConsumption(ctx, lastInsert, now)
		if err != nil {
			return fmt.Errorf("unable to get enedis consumption: %s", err)
		}

		return a.update(ctx, lastInsert, consumption)
	}, done)
}

func (a app) update(ctx context.Context, lastInsert time.Time, consumption datahub.Consumption) error {
	index := 0
	count := 0

	feedLine := func(stmt *sql.Stmt) error {
		if index == len(consumption.Master.Readings) {
			return db.ErrBulkEnded
		}

		value := a.handleReading(lastInsert, consumption.Master.Readings[index])
		index++

		if value == emptyValue {
			return nil
		}

		count++

		_, err := stmt.Exec(a.name, value.Timestamp, value.Valeur)
		return err
	}

	if err := a.save(ctx, feedLine); err != nil {
		return fmt.Errorf("unable to save datas: %s", err)
	}

	logger.Info("%d lines inserted", count)

	return nil
}

func (a app) handleReading(lastInsert time.Time, reading datahub.Reading) value {
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", reading.Date, a.tz)
	if err != nil {
		logger.Warn("ignoring `%+v`: invalid date format", reading)
		return emptyValue
	}

	if timestamp.Before(lastInsert) || timestamp.Equal(lastInsert) {
		logger.Warn("ignoring `%+v`: timestamp already inserted", reading)
	}

	valeur, err := strconv.ParseFloat(reading.Value, 64)
	if err != nil {
		logger.Warn("ignoring `%+v`: invalid value format", reading)
		return emptyValue
	}

	return value{
		Valeur:    valeur / 1000,
		Timestamp: reading.Date,
	}
}
