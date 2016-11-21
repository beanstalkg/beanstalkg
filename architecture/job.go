package architecture

import (
	"errors"
	"time"
)

type State int

const ( // iota is reset to 0
	READY State = iota  // = 0
	DELAYED  // = 1
	RESERVED  // = 2
	BURIED  // = 3
)

type job struct {
	id string
	Pri int64
	Delay int64 // time set as delay in seconds
	StartedDelayAt int64 // timestamp of when it was set to delayed
	StartedTTRAt int64 // timestamp of when it was reserved
	TTR int64 // time set as ttr in seconds
	Bytes int64
	Data string

	// states
	state State
}

func NewJob(id string, pri, delay, ttr, bytes int64, data string) *job {
	j := new(job)
	j.id = id
	j.Pri = pri
	j.Delay = delay
	j.TTR = ttr
	j.Bytes = bytes
	j.Data = data
	if j.Delay <= 0 {
		j.state = READY
		// add to the delayed queue
	} else {
		j.state = DELAYED
		// add to the ready queue
	}
	return j
}

/**

   put with delay               release with delay
  ----------------> [DELAYED] <------------.
                        |                   |
                        | (time passes)     |
                        |                   |
   put                  v     reserve       |       delete
  -----------------> [READY] ---------> [RESERVED] --------> *poof*
                       ^  ^                |  |
                       |   \  release      |  |
                       |    `-------------'   |
                       |                      |
                       | kick                 |
                       |                      |
                       |       bury           |
                    [BURIED] <---------------'
                       |
                       |  delete
                        `--------> *poof*
 */
func (j *job) SetState(state State) error {
	switch state {
	case READY:
		if j.state == RESERVED || j.state == DELAYED || j.state == BURIED {
			j.state = state
		} else {
			return errors.New("Invalid state transition to READY")
		}
	case DELAYED:
		if j.state == RESERVED {
			j.state = state
			j.StartedDelayAt = time.Now().Unix()
		} else {
			return errors.New("Invalid state transition to RESERVED")
		}
	case RESERVED:
		if j.state == READY {
			j.state = state
			j.StartedTTRAt = time.Now().Unix()
		} else {
			return errors.New("Invalid state transition to RESERVED")
		}
	case BURIED:
		if j.state == RESERVED {
			j.state = state
		} else {
			return errors.New("Invalid state transition to BURIED")
		}
	}
	return nil
}

// Return proper key according to the present job state
func (j job) Key() int64 {
	switch j.state {
	case READY:
		return j.Pri
	case DELAYED:
		// time remaining from Delay till it gets ready becomes priority
		return j.Delay - (time.Now().Unix() - j.StartedDelayAt)
	case RESERVED:
		// time remaining from TTR till it gets ready becomes the priority
		return j.TTR - (time.Now().Unix() - j.StartedTTRAt)
	}
	return 0
}

func (j job) Id() string {
	return j.id
}
