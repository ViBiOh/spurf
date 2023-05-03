package enedis

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/db"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

type App struct {
	db db.App

	file string
	name string
}

type Config struct {
	file *string
	name *string
}

func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		file: flags.New("File", "CSV export to load").Prefix(prefix).DocPrefix("enedis").String(fs, "", nil),
		name: flags.New("Name", "Name").Prefix(prefix).DocPrefix("enedis").String(fs, "home", nil),
	}
}

func New(config Config, db db.App) App {
	return App{
		file: strings.TrimSpace(*config.file),
		name: strings.TrimSpace(*config.name),
		db:   db,
	}
}

func (a App) Start(ctx context.Context) {
	logger.Fatal(a.handleFile(ctx, a.file))
}

func (a App) handleFile(ctx context.Context, filename string) error {
	if len(filename) == 0 {
		return errors.New("no filename provided")
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error while closing file: %s", err)
		}
	}()

	return a.handleLines(ctx, bufio.NewScanner(file))
}

func (a App) handleLines(ctx context.Context, scanner *bufio.Scanner) error {
	lastInsert, err := a.getLastFetch(ctx)
	if err != nil {
		return fmt.Errorf("get last fetch: %w", err)
	}

	count := 0

	feedLine := func() ([]any, error) {
		if !scanner.Scan() {
			return nil, nil
		}

	forward:
		value := handleLine(lastInsert, strings.TrimSpace(scanner.Text()))
		if value == emptyValue {
			if scanner.Scan() {
				goto forward
			}

			return nil, nil
		}

		count++
		lastInsert = value.Timestamp

		return []any{a.name, value.Timestamp, value.Valeur}, nil
	}

	if err := a.save(ctx, feedLine); err != nil {
		return fmt.Errorf("save datas: %w", err)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading line-by-line: %w", err)
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

	if !timestamp.After(lastInsert) {
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
		Timestamp: timestamp,
	}
}
