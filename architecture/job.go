package architecture

import (
	"errors"
	"github.com/satori/go.uuid"
	"time"
	"strconv"
	"log"
)

type State int

const ( // iota is reset to 0
	READY    State = iota // = 0
	DELAYED               // = 1
	RESERVED              // = 2
	BURIED                // = 3
)

const NANO = 1000000000

type Job struct {
	id             string
	Pri            int64
	Delay          int64 // time set as delay in seconds
	StartedDelayAt int64 // timestamp of when it was set to delayed
	StartedTTRAt   int64 // timestamp of when it was reserved
	TTR            int64 // time set as ttr in seconds
	Bytes          int64
	Data           string

	// states
	state State

	timestamp int64
}

func NewJob(id string, pri, delay, ttr, bytes int64, data string) *Job {
	j := new(Job)
	j.id = id
	j.Pri = pri
	j.Delay = delay
	j.TTR = ttr
	j.Bytes = bytes
	j.Data = data
	if j.Delay <= 0 {
		j.state = READY
		// add to the ready queue
	} else {
		j.state = DELAYED
		j.StartedDelayAt = time.Now().UnixNano()
		// add to the delayed queue
	}
	j.timestamp = time.Now().UnixNano()
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
func (j *Job) SetState(state State) error {
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
			j.StartedDelayAt = time.Now().UnixNano()
		} else {
			return errors.New("Invalid state transition to RESERVED")
		}
	case RESERVED:
		if j.state == READY {
			j.state = state
			j.StartedTTRAt = time.Now().UnixNano()
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

func (j *Job) State() State {
	return j.state
}

// Return proper key according to the present job state
func (j *Job) Key() int64 {
	switch j.state {
	case READY:
		return j.Pri
	case DELAYED:
		// time remaining from Delay till it gets ready becomes priority
		// log.Println(j.Delay * NANO, time.Now().UnixNano(), j.StartedDelayAt)
		return j.Delay * NANO - (time.Now().UnixNano() - j.StartedDelayAt)
	case RESERVED:
		// time remaining from TTR till it gets ready becomes the priority
		return j.TTR * NANO - (time.Now().UnixNano() - j.StartedTTRAt)
	}
	return 0
}

func (j *Job) Id() string {
	return j.id
}

func (j *Job) Timestamp() int64 {
	return j.timestamp
}

func (j *Job) Enqueued() {
	j.timestamp = time.Now().UnixNano()
}

func (j *Job) Dequeued() {
	j.timestamp = time.Now().UnixNano()
}

// AwaitingClient stores an awaiting client send channel for a tube
type AwaitingClient struct {
	id          string
	SendChannel chan Command
	Request     Command
	QueuedAt    int64
	Timeout	    int64
}

func NewAwaitingClient(request Command, sendChannel chan Command) *AwaitingClient {
	a := new(AwaitingClient)
	a.id = uuid.NewV1().String()
	a.Request = request
	a.SendChannel = sendChannel
	a.QueuedAt = time.Now().UnixNano()
	a.Timeout = -1
	if _, ok := request.Params["timeout"]; ok {
		timeout, err := strconv.ParseInt(request.Params["timeout"], 10, 0)
		if err == nil {
			a.Timeout = timeout * NANO // convert to nano seconds
		} else {
			log.Println(err)
		}
	}
	// log.Println(a)
	return a
}

func (w *AwaitingClient) Key() int64 {
	return w.QueuedAt
}

func (w *AwaitingClient) Timeleft() int64 {
	timeleft := w.Timeout - (time.Now().UnixNano() - w.QueuedAt)
	// log.Println(timeleft, w.Timeout, time.Now().UnixNano(), w.QueuedAt)
	return timeleft
}

func (w *AwaitingClient) Id() string {
	return w.id
}

func (w *AwaitingClient) Timestamp() int64 {
	return w.QueuedAt
}

func (w *AwaitingClient) Enqueued() {
	w.QueuedAt = time.Now().UnixNano()
}

func (w *AwaitingClient) Dequeued() {
}
