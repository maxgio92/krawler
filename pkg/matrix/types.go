package matrix

type Points interface{}

type Column struct {
	OrdinateIndex int
	Points        Points
}
