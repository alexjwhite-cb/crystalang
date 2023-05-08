package abstracter

// Represents functional declarations
type Declaration struct {
	// The name of what's being iterated
	Identity string

	Start uint

	// The type of declaration taking place (Variable/Object/Function)
	Type string

	// The names of the input into the declaration
	Args []string

	// What occurs during declaration (only for functions)
	// Contents interface{}

	// The values being exported
	// Outputs interface{}
}

func (d Declaration) Pos() uint {
	return d.Start
}
