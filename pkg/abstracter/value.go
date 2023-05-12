package abstracter

// Value represents constant and variable declarations
type Value struct {
}

func (v *Value) Pos() uint {
	return 0
}

func (v *Value) End() uint {
	return 0
}
