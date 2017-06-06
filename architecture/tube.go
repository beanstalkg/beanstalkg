package architecture

import (
	"errors"
	"strconv"
	"time"

	"github.com/op/go-logging"
	"github.com/vimukthi-git/beanstalkg/backend"
)

var log = logging.MustGetLogger("BEANSTALKG")

const QUEUE_FREQUENCY time.Duration = 20 * time.Millisecond // process every 20ms. TODO check why some clients get stuck when this is lower
const MAX_JOBS_PER_ITERATION int = 20                       // maximum number of jobs processed per queue per one cycle

// PriorityQueue is the interface that all backends should implement, See backend/min_heap.go for an example
type PriorityQueue interface {
	Init()
	// queue item
	Enqueue(item PriorityQueueItem)
	// get the highest priority item without removing
	Peek() (item PriorityQueueItem)
	// remove item from begining
	Dequeue() (item PriorityQueueItem)
	// find an item by id in the queue
	Find(id string) (item PriorityQueueItem)
	// delete an item and return it by id
	Delete(id string) PriorityQueueItem
	// size
	Size() int
}

// PriorityQueueItem is a single item in the PriorityQueue.
// This interface helps in isolating details of backend items
type PriorityQueueItem interface {
	Key() int64
	Id() string
	Timestamp() int64
	Enqueued()
	Dequeued()
}

type PriorityQueueCreator func() PriorityQueue

// Tube represents a single tube(queue) in the beanstalkg server
type Tube struct {
	Name                 string
	ready                PriorityQueue
	reserved             PriorityQueue
	delayed              PriorityQueue
	buried               PriorityQueue
	awaitingClients      PriorityQueue
	awaitingTimedClients map[string]*AwaitingClient
	db                   backend.Persister
}

func NewTube(name string, priorityQueueCreator PriorityQueueCreator) *Tube {
	tube := &Tube{
		Name:                 name,
		ready:                priorityQueueCreator(),
		reserved:             priorityQueueCreator(),
		delayed:              priorityQueueCreator(),
		buried:               priorityQueueCreator(),
		awaitingClients:      priorityQueueCreator(),
		awaitingTimedClients: make(map[string]*AwaitingClient),
	}
	tube.ready.Init()
	tube.delayed.Init()
	tube.reserved.Init()
	tube.buried.Init()
	tube.awaitingClients.Init()
	return tube
}

// Process runs all the necessary operations for upkeep of the tube. Just a convenience method.
func (tube *Tube) Process() {
	tube.ProcessDelayedQueue(MAX_JOBS_PER_ITERATION)
	tube.ProcessReservedQueue(MAX_JOBS_PER_ITERATION)
	tube.ProcessReadyQueue(MAX_JOBS_PER_ITERATION)
}

// ProcessDelayedQueue processes the Delayed queue.
// Which involves checking if any delayed jobs are ready to be served to clients
// and enqueuing those jobs that are ready on Ready queue
func (tube *Tube) ProcessDelayedQueue(limit int) {
	// log.Debug("Number of awaiting clients", tube.AwaitingClients.Size())
	counter := 1
	for delayedJob := tube.delayed.Peek(); delayedJob != nil &&
		delayedJob.Key() <= 0; delayedJob = tube.delayed.Peek() {
		log.Debug("QUEUE delayed job got ready: ", delayedJob)
		delayedJob = tube.delayed.Dequeue()
		delayedJob.(*Job).SetState(READY)
		tube.ready.Enqueue(delayedJob)
		if counter >= limit {
			break
		}
		counter++
	}
}

// ProcessReservedQueue processes the Reserved queue.
// Which involves checking if any reserved jobs have timed out while being processed by clients
// and enqueuing those jobs that have timeout on Ready queue, so that other clients can reserved it again.
func (tube *Tube) ProcessReservedQueue(limit int) {
	counter := 1
	// reserved jobs are put to ready
	for reservedJob := tube.reserved.Peek(); reservedJob != nil &&
		reservedJob.Key() <= 0; reservedJob = tube.reserved.Peek() {
		// log.Println("QUEUE found reserved job thats ready: ", reservedJob)
		reservedJob = tube.reserved.Dequeue()
		reservedJob.(*Job).SetState(READY)
		tube.ready.Enqueue(reservedJob)
		if counter >= limit {
			break
		}
		counter++
	}
}

// ProcessReadyQueue processes the Ready queue.
// Which involves pushing ready jobs to awaiting clients and putting the sent jobs
// to Reserved queue
func (tube *Tube) ProcessReadyQueue(limit int) {
	counter := 1
	// ready jobs are sent
	for tube.awaitingClients.Peek() != nil && tube.ready.Peek() != nil {
		availableClientConnection := tube.awaitingClients.Dequeue()
		client := availableClientConnection.(*AwaitingClient)
		// log.Println("QUEUE sending job to client: ", client.id)
		readyJob := tube.ready.Dequeue().(*Job)
		client.Request.Job = *readyJob
		client.SendChannel <- client.Request.Copy()
		readyJob.SetState(RESERVED)
		tube.reserved.Enqueue(readyJob)
		if counter >= limit {
			break
		}
		counter++
	}
}

// ProcessTimedClients reserves jobs for or times out the clients with a timeout
func (tube *Tube) ProcessTimedClients() {
	for id, client := range tube.awaitingTimedClients {
		// log.Println(client)
		if client.Timeleft() <= 0 {
			if tube.ready.Peek() != nil {
				readyJob := tube.ready.Dequeue().(*Job)
				client.Request.Job = *readyJob
				client.SendChannel <- client.Request.Copy()
				readyJob.SetState(RESERVED)
				tube.reserved.Enqueue(readyJob)
			} else {
				client.Request.Err = errors.New(TIMED_OUT)
				client.SendChannel <- client.Request.Copy()
			}
			delete(tube.awaitingTimedClients, id)
			tube.awaitingClients.Delete(id)
		}
	}
}

func (tube *Tube) Put(command *Command) {
	job := command.Job
	if db := tube.db; db != nil {
		db.Put(job, tube)
	}

	if job.State() == READY {
		// log.Println("TUBE_HANDLER put job to ready queue: ", c, name)
		tube.ready.Enqueue(&job)
	} else {
		// log.Println("TUBE_HANDLER put job to delayed queue: ", c, name)
		tube.delayed.Enqueue(&job)
	}
	command.Err = nil
	command.Params["id"] = job.Id()
}

func (tube *Tube) Reserve(command *Command, sendChannel chan Command) {
	tube.awaitingClients.Enqueue(NewAwaitingClient(*command, sendChannel))
}

func (tube *Tube) ReserveWithTimeout(command *Command, sendChannel chan Command) {
	client := NewAwaitingClient(*command, sendChannel)
	tube.awaitingClients.Enqueue(client)
	tube.awaitingTimedClients[client.Id()] = client
	tube.ProcessTimedClients()
}

func (tube *Tube) Delete(command *Command) {
	if db := tube.db; db != nil {
		db.Delete(command.Job, tube)
	}

	if tube.buried.Delete(command.Params["id"]) != nil ||
		tube.reserved.Delete(command.Params["id"]) != nil {
		// log.Println("TUBE_HANDLER deleted job: ", c, name)
		command.Err = nil
	} else {
		command.Err = errors.New(NOT_FOUND)
	}
}

func (tube *Tube) Release(command *Command) {
	item := tube.reserved.Delete(command.Params["id"])
	if item != nil {
		job := item.(*Job)
		// log.Println("TUBE_HANDLER released job: ", c, name)
		job.SetState(READY)
		tube.ready.Enqueue(job)
	} else {
		command.Err = errors.New(NOT_FOUND)
	}
}

func (tube *Tube) Bury(command *Command) {
	item := tube.reserved.Delete(command.Params["id"])
	if item != nil {
		job := item.(*Job)
		// log.Println("TUBE_HANDLER buried job: ", c, name)
		job.SetState(BURIED)
		tube.buried.Enqueue(job)
	} else {
		command.Err = errors.New(NOT_FOUND)
	}
}

func (tube *Tube) Kick(command *Command) {
	bound, err := strconv.Atoi(command.Params["bound"])
	if err != nil {
		command.Err = errors.New(NOT_FOUND)
	}
	size := tube.buried.Size()
	if size < bound {
		bound = size
	}
	for amount := 0; amount < bound; amount++ {
		item := tube.buried.Dequeue()
		job := item.(*Job)
		job.SetState(READY)
		tube.ready.Enqueue(job)
	}
}

func (tube *Tube) KickJob(command *Command) {
	item := tube.buried.Delete(command.Params["id"])
	if item != nil {
		job := item.(*Job)
		job.SetState(READY)
		tube.ready.Enqueue(job)
	} else {
		command.Err = errors.New(NOT_FOUND)
	}
}

func (tube *Tube) PersistTo(p Persister) {
	tube.db = p
}
