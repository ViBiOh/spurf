package enedis

import (
	"database/sql"
	"time"

	"github.com/ViBiOh/httputils/v2/pkg/db"
	"github.com/ViBiOh/httputils/v2/pkg/errors"
)

const lastFetch = `
SELECT
  MAX(ts)
FROM
  enedis_value;
`

func (a *app) getLastFetch() (lastTimestamp time.Time, err error) {
	if err = a.db.QueryRow(lastFetch).Scan(&lastTimestamp); err != nil {
		err = errors.WithStack(err)
		return
	}

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

func (a *app) saveValue(o *Value, tx *sql.Tx) (err error) {
	if o == nil {
		return errors.New("cannot save nil Value")
	}

	var usedTx *sql.Tx
	if usedTx, err = db.GetTx(a.db, tx); err != nil {
		return
	}

	if usedTx != tx {
		defer func() {
			err = db.EndTx(usedTx, err)
		}()
	}

	if _, err = usedTx.Exec(insertQuery, o.Timestamp, o.Valeur); err != nil {
		err = errors.Wrap(err, "unable to save %#v", o)
		return
	}

	return
}
