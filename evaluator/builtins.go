package evaluator

import (
	"fmt"
	"wordbuilder/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},

	"max": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			max, ok := args[0].(*object.Integer)
			if !ok {
				return newError("arguments to `max` not supported, got %s", args[0].Type())
			}

			for _, arg := range args {
				n, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `max` not supported, got %s", arg.Type())
				}
				if n.Value > max.Value {
					max = n
				}
			}
			return &object.Integer{Value: max.Value}
		},
	},

	"exists": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {

			if len(args) != 1 || args[0].Type() != object.StringObj {
				return newError("argument to `first` must be STRING, got %s", args[0].Type())
			}

			str := args[0].(*object.String)
			_, ok := env.Get(str.Value)
			return nativeBoolToBooleanIObject(ok)
		},
	},

	"first": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.ArrayObj {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},

	"last": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.ArrayObj {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},

	"rest": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.ArrayObj {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
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
	"push": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 && args[0].Type() != object.ArrayObj {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},

	"puts": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},

	"printwords": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			for k, v := range env.Store() {
				if v == nil {
					fmt.Printf("%q: \"\"\n", k)
				} else {
					fmt.Printf("%q: %q\n", k, v.Inspect())
				}
			}
			return NULL
		},
	},

	"wc": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			return &object.Integer{Value: int64(len(env.Store()))}
		},
	},

	"wordcount": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			c := 0
			for _, v := range env.Store() {
				if _, ok := v.(*object.Word); ok {
					c++
				}
			}
			return &object.Integer{Value: int64(c)}
		},
	},

	"refcount": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			c := 0
			for _, v := range env.Store() {
				if _, ok := v.(*object.Reference); ok {
					c++
				}
			}
			return &object.Integer{Value: int64(c)}
		},
	},

	"trcount": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			c := 0
			for _, v := range env.Store() {
				if _, ok := v.(*object.Translation); ok {
					c++
				}
			}
			return &object.Integer{Value: int64(c)}
		},
	},

	"cptcount": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			c := 0
			for _, v := range env.Store() {
				if _, ok := v.(*object.Concept); ok {
					c++
				}
			}
			return &object.Integer{Value: int64(c)}
		},
	},

	"mecount": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			c := 0
			for _, v := range env.Store() {
				if _, ok := v.(*object.MeThought); ok {
					c++
				}
			}
			return &object.Integer{Value: int64(c)}
		},
	},

	"thoughts": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			elements := []object.Object{}
			for _, th := range env.Thoughts() {
				elements = append(elements, &object.String{Value: th})
			}
			return &object.Array{Elements: elements}
		},
	},
}
