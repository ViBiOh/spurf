package enedis

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/cron"
	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

const (
	frenchDateFormat = "02/01/2006"
	isoDateFormat    = "2006-01-02"
)

// App of package
type App interface {
	Start()
}

// Config of package
type Config struct {
	email    *string
	password *string
	file     *string
	name     *string
	timezone *string
	cron     *bool
}

type app struct {
	email    string
	password string
	file     string
	name     string
	cron     bool
	cookies  []*http.Cookie

	location *time.Location
	db       *sql.DB
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		email:    flags.New(prefix, "enedis").Name("Email").Default("").Label("Email").ToString(fs),
		password: flags.New(prefix, "enedis").Name("Password").Default("").Label("Password").ToString(fs),
		file:     flags.New(prefix, "enedis").Name("File").Default("").Label("CSV export to load").ToString(fs),
		name:     flags.New(prefix, "enedis").Name("Name").Default("home").Label("Name").ToString(fs),
		timezone: flags.New(prefix, "enedis").Name("Timezone").Default("Europe/Paris").Label("Timezone").ToString(fs),
		cron:     flags.New(prefix, "enedis").Name("Cron").Default(false).Label("Start enedis as a cron").ToBool(fs),
	}
}

// New creates new App from Config
func New(config Config, db *sql.DB) (App, error) {
	location, err := time.LoadLocation(strings.TrimSpace(*config.timezone))
	if err != nil {
		return nil, err
	}

	email := strings.TrimSpace(*config.email)
	password := strings.TrimSpace(*config.password)
	file := strings.TrimSpace(*config.file)

	if len(file) == 0 && (len(email) == 0 || len(password) == 0) {
		return nil, errors.New("no credentials provided")
	}

	return &app{
		email:    email,
		password: password,
		file:     file,
		name:     strings.TrimSpace(*config.name),
		location: location,
		cron:     *config.cron,
		db:       db,
	}, nil
}

// Start the package
func (a *app) Start() {
	if len(a.file) != 0 {
		logger.Fatal(a.loadFile(a.file))
		return
	}

	if !a.cron {
		logger.Fatal(a.run(time.Now()))
		return
	}

	cron.New().Days().At("08:00").In(a.location.String()).Retry(time.Hour).MaxRetry(5).Now().Start(a.run, func(err error) {
		logger.Error("%s", err)
	})
}

func (a *app) loadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error while closing file: %s", err)
		}
	}()

	scanner := bufio.NewScanner(file)

	var datas []Value
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			logger.Warn("ignoring line `%s`: invalid format", line)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			logger.Warn("ignoring line `%s`: invalid date format", line)
			continue
		}

		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			logger.Warn("ignoring line `%s`: invalid value format", line)
			continue
		}

		datas = append(datas, Value{
			Valeur:    value / 1000,
			Timestamp: timestamp.Unix(),
		})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while checking scanner: %s", err)
	}

	logger.Info("Saving data for %s", filename)
	if err := a.save(context.Background(), datas); err != nil {
		return fmt.Errorf("error while saving datas: %s", err)
	}

	return nil
}

func (a *app) run(currentTime time.Time) error {
	if err := a.login(); err != nil {
		return err
	}

	ctx := context.Background()

	lastTimestamp, err := a.getLastFetch(ctx)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	logger.Info("Last fetch was %s", lastTimestamp)

	currentDate := currentTime.Format(isoDateFormat)
	lastSync := lastTimestamp.In(a.location).Truncate(oneDay).Add(oneDay)

	for lastSync.Format(isoDateFormat) < currentDate {
		date := lastSync.Format(isoDateFormat)

		logger.Info("Fetching data for %s", date)
		data, err := a.fetch(ctx, date)
		if err != nil {
			return err
		}

		logger.Info("Saving data for %s", date)
		if err := a.save(ctx, data.Graphe.Data); err != nil {
			return err
		}

		lastSync = lastSync.AddDate(0, 0, 1)
	}

	return nil
}

func (a *app) fetch(ctx context.Context, date string) (Consumption, error) {
	return a.getDataFromLegacy(ctx, date, true)
}

func (a *app) save(ctx context.Context, datas []Value) error {
	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		for _, value := range datas {
			if err := a.saveValue(ctx, value); err != nil {
				return err
			}
		}

		return nil
	})
}
