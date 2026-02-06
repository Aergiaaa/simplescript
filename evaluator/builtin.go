package evaluator

import (
	"fmt"

	"github.com/Aergiaaa/simplescript/object"
)

var builtins = map[string]*object.Builtin{
	// return len of a variable
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			case *object.Array:
				return &object.Integer{
					Value: int64(len(arg.Elems)),
				}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	// return the first elements of array
	"head": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments.")
			}

			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elems) != 0 {
				return arr.Elems[0]
			}

			return NULL
		},
	},

	// return last piece of array
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments.")
			}

			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elems) != 0 {
				return arr.Elems[len(arr.Elems)-1]
			}

			return NULL
		},
	},

	// returning its array without the first val
	"killHead": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments.")
			}

			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elems)
			if length != 0 {
				newElems := make([]object.Object, length-1)
				copy(newElems, arr.Elems[1:length])

				return &object.Array{
					Elems: newElems,
				}
			}

			return NULL
		},
	},

	// return array with added piece at the end
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments.")
			}

			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elems)

			newElems := make([]object.Object, length+1)
			copy(newElems, arr.Elems)
			newElems[length] = args[1]

			return &object.Array{
				Elems: newElems,
			}
		},
	},

	// print
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
