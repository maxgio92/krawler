package matrix

type Points interface{}

type Column struct {
	CurrentOrdinateIndex int
	Points               Points
}
