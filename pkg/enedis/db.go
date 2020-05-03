package enedis

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
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

	if _, err = tx.Exec(insertQuery, o.Timestamp, o.Valeur); err != nil {
		err = fmt.Errorf("unable to save %#v: %w", o, err)
		return
	}

	return
}
