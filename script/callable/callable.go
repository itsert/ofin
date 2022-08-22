package callable

import "github.com/itsert/ofin/script/environment"

type Callable interface {
	Arity() int
	Call(env *environment.Environment, arguments []interface{}) interface{}
}
