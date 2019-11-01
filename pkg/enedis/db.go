package enedis

import (
	"database/sql"
	"errors"
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

func (a *app) getLastFetch() (lastTimestamp time.Time, err error) {
	err = a.db.QueryRow(lastFetch).Scan(&lastTimestamp)

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
		err = fmt.Errorf("unable to save %#v: %w", o, err)
		return
	}

	return
}
