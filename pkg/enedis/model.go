package enedis

var (
	emptyConsumption = Consumption{}
)

// Consumption describes consumption response
type Consumption struct {
	Etat   Etat
	Graphe Graphe
}

// Etat describes status of output
type Etat struct {
	Valeur     string
	ErreurText string
}

// Graphe describes graphical data point
type Graphe struct {
	Data []Value
}

// Value describes data point
type Value struct {
	Valeur    float64
	Ordre     int64
	Timestamp int64
}
