package scheduler

import (
	"context"
	"time"
)

// Task are abstraction of scheduled work
type Task interface {
	Do(context.Context, time.Time) error
}
