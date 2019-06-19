package enedis

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/scheduler"
	"github.com/ViBiOh/httputils/pkg/tools"
)

const (
	loginURL   = "https://espace-client-connexion.enedis.fr/auth/UI/Login"
	consumeURL = "https://espace-client-particuliers.enedis.fr/group/espace-particuliers/suivi-de-consommation?"

	frenchDateFormat = "02/01/2006"
)

var _ scheduler.Task = &App{}

// Config of package
type Config struct {
	email    *string
	password *string
	timezone *string
}

// App of package
type App struct {
	email    string
	password string
	cookie   string

	location *time.Location
	db       *sql.DB
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	docPrefix := prefix
	if prefix == "" {
		docPrefix = "enedis"
	}

	return Config{
		email:    fs.String(tools.ToCamel(fmt.Sprintf("%sEmail", prefix)), "", fmt.Sprintf("[%s]  Email", docPrefix)),
		password: fs.String(tools.ToCamel(fmt.Sprintf("%sPassword", prefix)), "", fmt.Sprintf("[%s]  Password", docPrefix)),
		timezone: fs.String(tools.ToCamel(fmt.Sprintf("%sTimezone", prefix)), "Europe/Paris", fmt.Sprintf("[%s] Timezone", docPrefix)),
	}
}

// New creates new App from Config
func New(config Config, db *sql.DB) (*App, error) {
	location, err := time.LoadLocation(strings.TrimSpace(*config.timezone))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	email := strings.TrimSpace(*config.email)
	password := strings.TrimSpace(*config.password)

	if email == "" || password == "" {
		return nil, errors.New("no credentials provided")
	}

	return &App{
		email:    email,
		password: password,
		location: location,
		db:       db,
	}, nil
}

// Start the package
func (a *App) Start() error {
	return a.Do(context.Background(), time.Now().In(a.location))
}

// Do enedis fetch
func (a *App) Do(ctx context.Context, currentTime time.Time) error {
	if err := a.Login(); err != nil {
		return err
	}

	lastTimestamp, err := a.getLastFetch()
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	lastSync := lastTimestamp.In(a.location).Truncate(oneDay).Add(oneDay)

	for lastSync.Before(currentTime) {
		logger.Info("Fetching data for %s", lastSync.Format(frenchDateFormat))
		if err := a.fetchAndSave(context.Background(), lastSync); err != nil {
			return err
		}

		lastSync = lastSync.Add(oneDay)
	}

	return nil
}

func (a *App) fetchAndSave(ctx context.Context, currentTime time.Time) (err error) {
	var data *Consumption

	data, err = a.GetData(ctx, currentTime.Format(frenchDateFormat), true)
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
