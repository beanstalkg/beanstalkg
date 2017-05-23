package architecture

import (
	"errors"
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("BEANSTALKG")

const QUEUE_FREQUENCY time.Duration = 20 * time.Millisecond // process every 20ms. TODO check why some clients get stuck when this is lower
const MAX_JOBS_PER_ITERATION int = 20 // maximum number of jobs processed per queue per one cycle

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

// Tube represents a single tube(queue) in the beanstalkg server
type Tube struct {
	Name                 string
	Ready                PriorityQueue
	Reserved             PriorityQueue
	Delayed              PriorityQueue
	Buried               PriorityQueue
	AwaitingClients      PriorityQueue
	AwaitingTimedClients map[string]*AwaitingClient
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
	for delayedJob := tube.Delayed.Peek(); delayedJob != nil &&
			delayedJob.Key() <= 0; delayedJob = tube.Delayed.Peek() {
		log.Debug("QUEUE delayed job got ready: ", delayedJob)
		delayedJob = tube.Delayed.Dequeue()
		delayedJob.(*Job).SetState(READY)
		tube.Ready.Enqueue(delayedJob)
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
	for reservedJob := tube.Reserved.Peek(); reservedJob != nil &&
			reservedJob.Key() <= 0; reservedJob = tube.Reserved.Peek() {
		// log.Println("QUEUE found reserved job thats ready: ", reservedJob)
		reservedJob = tube.Reserved.Dequeue()
		reservedJob.(*Job).SetState(READY)
		tube.Ready.Enqueue(reservedJob)
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
	for tube.AwaitingClients.Peek() != nil && tube.Ready.Peek() != nil {
		availableClientConnection := tube.AwaitingClients.Dequeue()
		client := availableClientConnection.(*AwaitingClient)
		// log.Println("QUEUE sending job to client: ", client.id)
		readyJob := tube.Ready.Dequeue().(*Job)
		client.Request.Job = *readyJob
		client.SendChannel <- client.Request.Copy()
		readyJob.SetState(RESERVED)
		tube.Reserved.Enqueue(readyJob)
		if counter >= limit {
			break
		}
		counter++
	}
}

// ProcessTimedClients reserves jobs for or times out the clients with a timeout
func (tube *Tube) ProcessTimedClients() {
	for id, client := range tube.AwaitingTimedClients {
		// log.Println(client)
		if client.Timeleft() <= 0 {
			if tube.Ready.Peek() != nil {
				readyJob := tube.Ready.Dequeue().(*Job)
				client.Request.Job = *readyJob
				client.SendChannel <- client.Request.Copy()
				readyJob.SetState(RESERVED)
				tube.Reserved.Enqueue(readyJob)
			} else {
				client.Request.Err = errors.New(TIMED_OUT)
				client.SendChannel <- client.Request.Copy()
			}
			delete(tube.AwaitingTimedClients, id)
			tube.AwaitingClients.Delete(id)
		}
	}
}
