package enedis

var (
	emptyValue = value{}
)

type value struct {
	Timestamp string
	Valeur    float64
}
