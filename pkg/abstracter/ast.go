package abstracter

import (
	"fmt"

	"github.com/alexjwhite-cb/crystalang/pkg/tokeniser"
)

type NodeFunc int

const (
	NewDeclaration NodeFunc = iota
	NewStatement
	NewValue
	NewSpoon
	NewPassthrough
	NewError
)

type Node interface {
	Pos() uint // key of opening node
}

type AbstractSyntaxTree struct {
	Nodes []Node
}

// type Iteration struct {
// 	Identity string
// 	Value    int
// }

// type FunctionCall struct {
// 	// Function string
// 	// Arguments
// }

// type Passthrough struct {
// }

func NewAbstractSyntaxTree() *AbstractSyntaxTree {
	return &AbstractSyntaxTree{}
}

func (a *AbstractSyntaxTree) Abstract(tokens []tokeniser.Token) *AbstractSyntaxTree {
	for i := 0; i < len(tokens); i++ {
		fmt.Printf("%v, %+v\n", i, tokens[i])

	}
	return a
}

//func inReservedValues(x map[tokeniser.Token]string) NodeFunc {
//	for t, v := range x {
//
//	}
//}
