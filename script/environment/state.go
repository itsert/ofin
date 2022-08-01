package environment

import (
	"fmt"
)

type ProgramState int

const (
	GLOBAL ProgramState = iota
	SCENARIO
	GIVEN
	WHEN
	THEN
)

type state struct {
	currentState ProgramState
}

func NewState() *state {
	return &state{
		currentState: GLOBAL,
	}
}

func (s *state) Transition(newState ProgramState) error {
	message := "Invalid state transition from %+v to %+v"
	switch newState {
	case GLOBAL:
		if s.currentState == THEN {
			s.currentState = newState
			return nil
		} else {
			return fmt.Errorf(message, s.currentState, newState)
		}
		break
	case SCENARIO:
		break
	case GIVEN:
		break
	case WHEN:
		break
	case THEN:
		break
	}

	return nil
}
