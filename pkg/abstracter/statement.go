package abstracter

type Statement struct {
	Identity string
	Start    uint
}

func (s Statement) Pos() uint {
	return s.Start
}
