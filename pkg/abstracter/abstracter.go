package abstracter

import (
	"fmt"

	"github.com/alexjwhite-cb/crystalang/pkg/tokeniser"
)

type AbstractSyntaxTree struct {
}

type Declaration struct {
	// The name of what's being iterated
	Identity string

	// The type of declaration taking place (Variable/Object/Function)
	Type string

	// The names of the input into the declaration
	Args []string

	// What occurs during declaration (only for functions)
	// Contents interface{}

	// The values being exported
	// Outputs interface{}
}

type Comparison struct {
	Identity string
}

type Iteration struct {
	Identity string
	Value    int
}

type FunctionCall struct {
	// Function string
	// Arguments
}

type Value struct {
}

func NewAbstractSyntaxTree() *AbstractSyntaxTree {
	return &AbstractSyntaxTree{}
}

func (a *AbstractSyntaxTree) Abstract(tokens tokeniser.TokenMap) *AbstractSyntaxTree {
	for i := 0; i < len(tokens); i++ {
		fmt.Printf("%v, %+v\n", i, tokens[i])
	}
	return a
}
