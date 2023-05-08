package abstracter

// Represents constant and variable declarations
type Value struct {
	Identity string
	Value    any
	Start    uint
}

func (v Value) Pos() uint {
	return 0
}
