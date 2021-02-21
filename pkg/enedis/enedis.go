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

	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
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
	db *sql.DB

	file string
	name string
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

	return a.handleLines(bufio.NewScanner(file))
}

func (a app) handleLines(scanner *bufio.Scanner) error {
	lastInsert, err := a.getLastFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get last fetch: %s", err)
	}

	count := 0

	feedLine := func(stmt *sql.Stmt) error {
		if !scanner.Scan() {
			return db.ErrBulkEnded
		}

		value := handleLine(lastInsert, strings.TrimSpace(scanner.Text()))
		if value == emptyValue {
			return nil
		}

		count++

		_, err := stmt.Exec(a.name, value.Timestamp, value.Valeur)
		return err
	}

	if err := a.save(context.Background(), feedLine); err != nil {
		return fmt.Errorf("unable to save datas: %s", err)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading line-by-line: %s", err)
	}

	logger.Info("%d lines inserted", count)

	return nil
}

func handleLine(lastInsert time.Time, line string) value {
	parts := strings.Split(line, ";")
	if len(parts) != 2 {
		logger.Warn("ignoring line `%s`: invalid format", line)
		return emptyValue
	}

	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		logger.Warn("ignoring line `%s`: invalid date format", line)
		return emptyValue
	}

	if timestamp.Before(lastInsert) || timestamp.Equal(lastInsert) {
		logger.Warn("ignoring line `%s`: timestamp already inserted", line)
		return emptyValue
	}

	valeur, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		logger.Warn("ignoring line `%s`: invalid value format", line)
		return emptyValue
	}

	return value{
		Valeur:    valeur / 1000,
		Timestamp: parts[0],
	}
}
