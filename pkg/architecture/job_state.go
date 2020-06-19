package architecture

import "errors"

type updateFunc func(*Job, int64)

// helper functions for updating timers.
func updateTTR(job *Job, t int64) {
	job.StartedTTRAt = t
}

func updateDelayed(job *Job, t int64) {
	job.StartedDelayAt = t
}

// Describes valid transitions to a given state.
// key = to state, value = valid from states.
var validTransitionsTo = map[State]map[State]bool{
	READY:    map[State]bool{RESERVED: true, DELAYED: true, BURIED: true},
	DELAYED:  map[State]bool{RESERVED: true},
	RESERVED: map[State]bool{READY: true},
	BURIED:   map[State]bool{RESERVED: true},
}

var transitionErrors = map[State]error{
	READY:    errors.New("Invalid state transition to READY"),
	DELAYED:  errors.New("Invalid state transition to DELAYED"),
	RESERVED: errors.New("Invalid state transition to RESERVED"),
	BURIED:   errors.New("Invalid state transition to BURIED"),
}

var updateFuncs = map[State]updateFunc{
	DELAYED:  updateTTR,
	RESERVED: updateDelayed,
}
