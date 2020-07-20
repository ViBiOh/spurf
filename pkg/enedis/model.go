package enedis

var (
	emptyValue = value{}
)

type value struct {
	Valeur    float64
	Timestamp string
}
