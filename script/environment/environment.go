package environment

import (
	"fmt"

	"github.com/itsert/ofin/merror"
	"github.com/itsert/ofin/script/token"
)

type Environment struct {
	Global map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		Global: map[string]interface{}{},
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Global[name] = value
}
func (e *Environment) Assign(name token.Token, value interface{}) {
	if _, ok := e.Global[name.Lexeme]; ok {
		e.Global[name.Lexeme] = value
	} else {
		merror.RuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}
}
func (e *Environment) Get(name token.Token) (interface{}, error) {
	if v, ok := e.Global[name.Lexeme]; ok {
		return v, nil
	}
	// merror.RuntimeError(name, fmt.Sprintf("Variable %s is undefined.", name.Lexeme))
	return nil, fmt.Errorf("variable %s is undefined", name.Lexeme)
}
