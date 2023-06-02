package evaluator

import (
	"fmt"
	"github.com/alexjwhite-cb/jet/pkg/object"
	"strings"
)

var builtins = map[string]*object.BuiltIn{
	"puts": {
		Method: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return nil
		},
	},

	"print": {
		Method: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				switch a := arg.(type) {
				case *object.String:
					out = append(out, a.Value)

				default:
					return newError("argument to `print` not supported, got %s", arg.Type())
				}
			}
			print(strings.Join(out, ""))
			return nil
		},
	},

	"len": {
		Method: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("len: incorrect argument count; want 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},

	"first": {
		Method: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, want 1, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},

	"tail": {
		Method: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, want 1, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `tail` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},

	"append": {
		Method: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, want 2, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `append` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
}
