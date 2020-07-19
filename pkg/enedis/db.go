package enedis

import (
	"context"
	"database/sql"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/db"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
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

func (a app) save(ctx context.Context, datas []Value) error {
	logger.Info("Saving %d records", len(datas))

	return db.DoAtomic(ctx, a.db, func(ctx context.Context) error {
		var index int
		feeder := func(stmt *sql.Stmt) error {
			if index == len(datas) {
				return db.ErrBulkEnded
			}

			data := datas[index]
			index++

			_, err := stmt.Exec(a.name, data.Timestamp, data.Valeur)
			return err
		}

		return db.Bulk(ctx, feeder, "spurf", "enedis_value", "name", "ts", "value")
	})
}
