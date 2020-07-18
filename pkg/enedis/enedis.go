package enedis

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

// App of package
type App interface {
	Start()
}

// Config of package
type Config struct {
	file *string
	name *string
}

type app struct {
	file string
	name string

	db *sql.DB
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		file: flags.New(prefix, "enedis").Name("File").Default("").Label("CSV export to load").ToString(fs),
		name: flags.New(prefix, "enedis").Name("Name").Default("home").Label("Name").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, db *sql.DB) App {
	return &app{
		file: strings.TrimSpace(*config.file),
		name: strings.TrimSpace(*config.name),
		db:   db,
	}
}

// Start the package
func (a app) Start() {
	logger.Fatal(a.handleFile(a.file))
}

func (a app) handleFile(filename string) error {
	if len(filename) == 0 {
		return errors.New("no filename provided")
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error while closing file: %s", err)
		}
	}()

	lastInsert, err := a.getLastFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get last fetch: %s", err)
	}

	scanner := bufio.NewScanner(file)

	var datas []Value
	for scanner.Scan() {
		datas = handleLine(datas, lastInsert, strings.TrimSpace(scanner.Text()))
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

func handleLine(datas []Value, lastInsert time.Time, line string) []Value {
	parts := strings.Split(line, ";")
	if len(parts) != 2 {
		logger.Warn("ignoring line `%s`: invalid format", line)
		return datas
	}

	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		logger.Warn("ignoring line `%s`: invalid date format", line)
		return datas
	}

	if timestamp.Before(lastInsert) || timestamp.Equal(lastInsert) {
		logger.Warn("ignoring line `%s`: timestamp already inserted", line)
		return datas
	}

	value, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		logger.Warn("ignoring line `%s`: invalid value format", line)
		return datas
	}

	return append(datas, Value{
		Valeur:    value / 1000,
		Timestamp: parts[0],
	})
}
