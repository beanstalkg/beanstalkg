package architecture

import (
	"strconv"
	"time"

	"github.com/satori/go.uuid"
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
	job := &Job{
		id:    id,
		Pri:   pri,
		Delay: delay,
		TTR:   ttr,
		Bytes: bytes,
		Data:  data,
	}
	if job.Delay <= 0 {
		job.state = READY
		// add to the ready queue
	} else {
		job.state = DELAYED
		job.StartedDelayAt = time.Now().UnixNano()
		// add to the delayed queue
	}
	job.timestamp = time.Now().UnixNano()
	return job
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
func (job *Job) SetState(state State) error {
	// Ensure the desired to state is valid.
	validFrom, ok := validTransitionsTo[state]
	if !ok {
		return nil
	}

	// The to state is valid.  Now ensure that we are coming from
	// a valid state.
	if validp, _ := validFrom[job.state]; !validp {
		return transitionErrors[state]
	}

	// All is well here.
	job.state = state

	// If a timer needs updating, do so.
	if f, ok := updateFuncs[state]; ok {
		f(job, time.Now().UnixNano())
	}

	return nil
}

func (job *Job) State() State {
	return job.state
}

// Return proper key according to the present job state
func (job *Job) Key() int64 {
	switch job.state {
	case READY:
		return job.Pri
	case DELAYED:
		// time remaining from Delay till it gets ready becomes priority
		// log.Println(job.Delay * NANO, time.Now().UnixNano(), job.StartedDelayAt)
		return job.Delay*NANO - (time.Now().UnixNano() - job.StartedDelayAt)
	case RESERVED:
		// time remaining from TTR till it gets ready becomes the priority
		return job.TTR*NANO - (time.Now().UnixNano() - job.StartedTTRAt)
	}
	return 0
}

func (job *Job) Id() string {
	return job.id
}

func (job *Job) Timestamp() int64 {
	return job.timestamp
}

func (job *Job) Enqueued() {
	job.timestamp = time.Now().UnixNano()
}

func (job *Job) Dequeued() {
	job.timestamp = time.Now().UnixNano()
}

// AwaitingClient stores an awaiting client send channel for a tube
type AwaitingClient struct {
	id          string
	SendChannel chan Command
	Request     Command
	QueuedAt    int64
	Timeout     int64
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
			log.Error(err)
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
