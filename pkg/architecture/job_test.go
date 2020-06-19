package architecture

import (
	"testing"
)

var (
	validTransitionToTest = map[State][]State{
		READY:    {RESERVED, DELAYED, BURIED},
		DELAYED:  {RESERVED},
		RESERVED: {READY},
		BURIED:   {RESERVED},
	}

	invalidTransitionToTest = map[State][]State{
		READY:    {READY},
		DELAYED:  {READY, DELAYED, BURIED},
		RESERVED: {RESERVED, DELAYED, BURIED},
		BURIED:   {READY, DELAYED, BURIED},
	}
)

// State's stringer implementation is valuable in grokking tests, but
// it's not really needed elsewhere.  I could use map lookups, but
// these tests are fast enough as is.
func (s State) String() string {
	switch s {
	case READY:
		return "READY"
	case DELAYED:
		return "DELAYED"
	case RESERVED:
		return "RESERVED"
	case BURIED:
		return "BURIED"
	}

	return "INVALID"
}

func TestJob_SetState(t *testing.T) {
	for toState, fromStates := range validTransitionToTest {
		for _, fromState := range fromStates {
			testJ := &Job{state: fromState}

			if err := testJ.SetState(toState); err != nil {
				t.Errorf("State transition failed: %s -> %s; %s", testJ.state, toState, err)
			}

			if testJ.state != toState {
				t.Errorf("SetState succeeded, but Job's state wasn't changed: %s -> %s", testJ.state, toState)
			}
		}
	}

	for toState, fromStates := range invalidTransitionToTest {
		for _, fromState := range fromStates {
			testJ := &Job{state: fromState}

			if err := testJ.SetState(toState); err == nil {
				t.Errorf("State transition succeeded; should have failed: %s -> %s", testJ.state, toState)
			}

			if testJ.state != fromState {
				t.Errorf("SetState failed, but Job's state changed: %s -> %s", fromState, testJ.state)
			}
		}
	}
}
