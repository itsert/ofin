package environment

import (
	"fmt"

	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/script/token"
)

type Environment struct {
	value     map[string]interface{}
	enclosing *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		value:     map[string]interface{}{},
		enclosing: nil,
	}
}
func NewEnvironmentWithParent(enclosing *Environment) *Environment {
	return &Environment{
		value:     map[string]interface{}{},
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.value[name] = value
}
func (e *Environment) Assign(name token.Token, value interface{}) {
	if _, ok := e.value[name.Lexeme]; ok {
		e.value[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	merror.RuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")

}
func (e *Environment) Get(name token.Token) (interface{}, error) {
	if v, ok := e.value[name.Lexeme]; ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, fmt.Errorf("variable %s is undefined", name.Lexeme)
}
