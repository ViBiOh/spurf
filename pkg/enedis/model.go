package enedis

import "time"

var emptyValue = value{}

type value struct {
	Timestamp time.Time
	Valeur    float64
}
