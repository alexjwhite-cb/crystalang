package evaluator

import (
	"github.com/alexjwhite-cb/jet/pkg/object"
)

var builtins = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Method: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("len: incorrect argument count; want 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
}
