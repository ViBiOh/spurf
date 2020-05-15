package enedis

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

// StartAtomic starts atomic work
func StartAtomic(ctx context.Context, usedDB *sql.DB) (context.Context, error) {
	if db.ReadTx(ctx) != nil {
		return ctx, nil
	}

	tx, err := usedDB.Begin()
	if err != nil {
		return ctx, err
	}

	return db.StoreTx(ctx, tx), nil
}

// EndAtomic ends atomic work
func EndAtomic(ctx context.Context, err error) error {
	tx := db.ReadTx(ctx)
	if tx == nil {
		return err
	}

	return db.EndTx(tx, err)
}

const lastFetch = `
SELECT
  MAX(ts)
FROM
  enedis_value;
`

func (a *app) getLastFetch(ctx context.Context) (lastTimestamp time.Time, err error) {
	scanner := func(row db.RowScanner) error {
		return row.Scan(&lastTimestamp)
	}
	err = db.GetRow(ctx, a.db, scanner, lastFetch)

	return
}

const insertQuery = `
INSERT INTO
  enedis_value
(
  ts,
  value
) VALUES (
  to_timestamp($1),
  $2
);
`

func (a *app) saveValue(ctx context.Context, o Value) (err error) {
	err = db.Exec(ctx, a.db, insertQuery, o.Timestamp, o.Valeur)
	if err != nil {
		err = fmt.Errorf("unable to save %#v: %w", o, err)
		return
	}

	return
}
