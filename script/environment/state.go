package environment

import (
	"fmt"
	"sync"
)

type State string

const (
	GLOBAL   State = "GLOBAL"
	SCENARIO       = "SCENARIO"
	GIVEN          = "GIVEN"
	WHEN           = "WHEN"
	THEN           = "THEN"
)

type ProgramState struct {
	currentState         State
	stateTransitionTable map[StateTransitionTupple]TransitionFunc
	// mutex ensures that only 1 event is processed by the state machine at any given time.
	mutex sync.Mutex
}

type StateTransitionTupple struct {
	initialState State
	newState     State
}

type TransitionFunc func(state *State, newState State)

func transitionFuncImpl(state *State, newState State) {
	*state = newState
}

func NewState() *ProgramState {
	return &ProgramState{
		currentState: GLOBAL,
		stateTransitionTable: map[StateTransitionTupple]TransitionFunc{
			{GLOBAL, SCENARIO}: transitionFuncImpl,
			{GLOBAL, GIVEN}:    transitionFuncImpl, // strictly for testing purpose
			{SCENARIO, GIVEN}:  transitionFuncImpl,
			{SCENARIO, WHEN}:   transitionFuncImpl,
			{GIVEN, WHEN}:      transitionFuncImpl,
			{GIVEN, THEN}:      transitionFuncImpl,
			{GIVEN, SCENARIO}:  transitionFuncImpl,
			{WHEN, THEN}:       transitionFuncImpl,
			{THEN, GLOBAL}:     transitionFuncImpl,
			{THEN, SCENARIO}:   transitionFuncImpl,
			{THEN, WHEN}:       transitionFuncImpl,
		},
	}
}

func (p *ProgramState) Transition(newState State) (State, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// if p.IsState(newState) {
	// 	return newState, nil
	// }
	tupple := StateTransitionTupple{p.currentState, newState}
	if f, ok := p.stateTransitionTable[tupple]; !ok {
		return p.currentState, fmt.Errorf("invalid state transition from %+v to %+v", p.currentState, newState)
	} else {
		f(&p.currentState, newState)
		return p.currentState, nil
	}
}

func (p *ProgramState) CurrentState() State {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.currentState
}

func (p *ProgramState) IsState(state State) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.currentState == state
}

func (p *ProgramState) NotState(state State) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.currentState != state
}
