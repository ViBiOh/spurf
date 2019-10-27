package enedis

import (
	"context"
	"database/sql"
	"flag"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v2/pkg/cron"
	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/errors"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/tools"
)

const (
	loginURL   = "https://espace-client-connexion.enedis.fr/auth/UI/Login"
	consumeURL = "https://espace-client-particuliers.enedis.fr/group/espace-particuliers/suivi-de-consommation?"

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
	timezone *string
}

type app struct {
	email    string
	password string
	cookie   string

	location *time.Location
	db       *sql.DB
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		email:    tools.NewFlag(prefix, "enedis").Name("Email").Default("").Label("Email").ToString(fs),
		password: tools.NewFlag(prefix, "enedis").Name("Password").Default("").Label("Password").ToString(fs),
		timezone: tools.NewFlag(prefix, "enedis").Name("Timezone").Default("Europe/Paris").Label("Timezone").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, db *sql.DB) (App, error) {
	location, err := time.LoadLocation(strings.TrimSpace(*config.timezone))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	email := strings.TrimSpace(*config.email)
	password := strings.TrimSpace(*config.password)

	if email == "" || password == "" {
		return nil, errors.New("no credentials provided")
	}

	return &app{
		email:    email,
		password: password,
		location: location,
		db:       db,
	}, nil
}

// Start the package
func (a app) Start() {
	cron.New().Days().At("08:00").In(a.location.String()).Retry(time.Hour).MaxRetry(5).Now().Start(a.fetch, func(err error) {
		logger.Error("%+v", err)
	})
}

func (a app) fetch(currentTime time.Time) error {
	if err := a.login(); err != nil {
		return err
	}

	lastTimestamp, err := a.getLastFetch()
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	currentDate := currentTime.Format(isoDateFormat)
	lastSync := lastTimestamp.In(a.location).Truncate(oneDay).Add(oneDay)

	for lastSync.Format(isoDateFormat) != currentDate {
		lastSyncFrench := lastSync.Format(frenchDateFormat)

		logger.Info("Fetching data for %s", lastSyncFrench)
		if err := a.fetchAndSave(context.Background(), lastSyncFrench); err != nil {
			return err
		}

		lastSync = lastSync.Add(oneDay)
	}

	return nil
}

func (a app) fetchAndSave(ctx context.Context, date string) (err error) {
	var data *Consumption

	data, err = a.getData(ctx, date, true)
	if err != nil {
		return
	}

	var tx *sql.Tx
	if tx, err = db.GetTx(a.db, nil); err != nil {
		return
	}

	defer func() {
		err = db.EndTx(tx, err)
	}()

	for _, value := range data.Graphe.Data {
		if err = a.saveValue(value, tx); err != nil {
			return
		}
	}

	return
}
