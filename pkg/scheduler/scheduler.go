package scheduler

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/tools"
)

// Config of package
type Config struct {
	hour     *int
	minute   *int
	interval *string
	retry    *string
	timezone *string
}

// App of package
type App struct {
	hour     int
	minute   int
	location *time.Location

	interval time.Duration
	retry    time.Duration

	task Task
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	docPrefix := prefix
	if prefix == "" {
		docPrefix = "scheduler"
	}

	return Config{
		hour:     fs.Int(tools.ToCamel(fmt.Sprintf("%sHour", prefix)), 8, fmt.Sprintf("[%s] Hour of running", docPrefix)),
		minute:   fs.Int(tools.ToCamel(fmt.Sprintf("%sMinute", prefix)), 0, fmt.Sprintf("[%s] Minute of running", docPrefix)),
		timezone: fs.String(tools.ToCamel(fmt.Sprintf("%sTimezone", prefix)), "Europe/Paris", fmt.Sprintf("[%s] Timezone of running", docPrefix)),
		interval: fs.String(tools.ToCamel(fmt.Sprintf("%sInterval", prefix)), "24h", fmt.Sprintf("[%s] Duration between two runs", docPrefix)),
		retry:    fs.String(tools.ToCamel(fmt.Sprintf("%sRetry", prefix)), "10m", fmt.Sprintf("[%s] Duration between two retries", docPrefix)),
	}
}

// New creates new App from Config
func New(config Config, task Task) (*App, error) {
	location, err := time.LoadLocation(strings.TrimSpace(*config.timezone))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	interval, err := time.ParseDuration(strings.TrimSpace(*config.interval))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	retry, err := time.ParseDuration(strings.TrimSpace(*config.retry))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &App{
		hour:     *config.hour,
		minute:   *config.minute,
		interval: interval,
		retry:    retry,
		location: location,
		task:     task,
	}, nil
}

// Start scheduler
func (a App) Start() {
	for {
		a.scheduler()
	}
}

func (a App) getNextTick() (time.Time, time.Time) {
	currentTime := time.Now().In(a.location)
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), a.hour, a.minute, 0, 0, a.location), currentTime
}

func (a App) getTimer() *time.Timer {
	nextTime, currentTime := a.getNextTick()
	if !nextTime.After(currentTime) {
		nextTime = nextTime.Add(a.interval)
	}

	logger.Info("Next run at %v", nextTime)

	return time.NewTimer(time.Until(nextTime))
}

func (a App) scheduler() {
	timer := a.getTimer()

	for {
		select {
		case currentTime := <-timer.C:
			ctx := context.Background()

			if err := a.task.Do(ctx, currentTime); err != nil {
				logger.Error(`%+v`, err)

				timer.Reset(a.retry)
				logger.Warn("Retrying in 10 minutes")
			} else {
				return
			}
		}
	}
}
