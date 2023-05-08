package abstracter

type Expression struct {
	Start uint
}

func (e Expression) Pos() uint {
	return e.Start
}
