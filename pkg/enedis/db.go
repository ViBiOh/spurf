package enedis

import (
	"context"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

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
	err = db.Exec(ctx, insertQuery, o.Timestamp, o.Valeur)
	if err != nil {
		err = fmt.Errorf("unable to save %#v: %w", o, err)
		return
	}

	return
}
