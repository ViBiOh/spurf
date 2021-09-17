package enedis

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

const lastFetch = `
SELECT
  MAX(ts)
FROM
  spurf.enedis_value
WHERE
  name = $1;
`

func (a App) getLastFetch(ctx context.Context) (time.Time, error) {
	var output time.Time

	scanner := func(row pgx.Row) error {
		return row.Scan(&output)
	}

	return output, a.db.Get(ctx, scanner, lastFetch, a.name)
}

func (a App) save(ctx context.Context, feeder func() ([]interface{}, error)) error {
	return a.db.DoAtomic(ctx, func(ctx context.Context) error {
		return a.db.Bulk(ctx, feeder, "spurf", "enedis_value", "name", "ts", "value")
	})
}
