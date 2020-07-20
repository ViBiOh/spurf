package enedis

import (
	"context"
	"database/sql"
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

func (a app) getLastFetch(ctx context.Context) (time.Time, error) {
	var output time.Time

	scanner := func(row *sql.Row) error {
		return row.Scan(&output)
	}

	return output, db.Get(ctx, a.db, scanner, lastFetch, a.name)
}

func (a app) save(ctx context.Context, feeder func(stmt *sql.Stmt) error) error {
	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		return db.Bulk(ctx, feeder, "spurf", "enedis_value", "name", "ts", "value")
	})
}
