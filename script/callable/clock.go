package callable

import (
	"github.com/itsert/ofin/script/environment"
	"time"
)

type Clock struct{}

func NewClock() Clock {
	return Clock{}
}
func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(env *environment.Environment, arguments []interface{}) interface{} {
	return float64(time.Now().UnixMilli()) / 1000.0
}
