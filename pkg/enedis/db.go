package enedis

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
)

const lastFetch = `
SELECT
  MAX(ts)
FROM
  spurf.enedis_value
WHERE
  name = $1;
`

func (a *app) getLastFetch(ctx context.Context) (lastTimestamp time.Time, err error) {
	scanner := func(row *sql.Row) error {
		return row.Scan(&lastTimestamp)
	}
	err = db.Get(ctx, a.db, scanner, lastFetch, a.name)

	return
}

const insertQuery = `
INSERT INTO
  spurf.enedis_value
(
  name,
  ts,
  value
) VALUES (
  $1,
  to_timestamp($2),
  $3
);
`

func (a *app) saveValue(ctx context.Context, o Value) (err error) {
	err = db.Exec(ctx, insertQuery, a.name, o.Timestamp, o.Valeur)
	if err != nil {
		err = fmt.Errorf("unable to save %#v: %w", o, err)
		return
	}

	return
}
