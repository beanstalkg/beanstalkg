package operation

import (
	"errors"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"github.com/vimukthi-git/beanstalkg/backend"
	// "log"
	"time"
	"strconv"
)

func NewTubeHandler(
	name string,
	commands chan architecture.Command,
	watchedTubeConnectionsReceiver chan chan architecture.Command,
	stop chan bool,
) {
	// commands := make(chan architecture.Command)
	go func() {
		tube := createTube(name)
		ticker := time.NewTicker(architecture.QUEUE_FREQUENCY)
		// TODO make sure all logic that can be moved to Tube struct is moved there
		for {
			select {
			case <-ticker.C:
				// log.Println("House Keeping Started for: ", name)
				tube.Process()
				tube.ProcessTimedClients()
			case c := <-commands:
				log.Info("tube handling ", c.Name)
				switch c.Name {
				case architecture.PUT:
					if c.Job.State() == architecture.READY {
						// log.Println("TUBE_HANDLER put job to ready queue: ", c, name)
						tube.Ready.Enqueue(&c.Job)
					} else {
						// log.Println("TUBE_HANDLER put job to delayed queue: ", c, name)
						tube.Delayed.Enqueue(&c.Job)
					}
					c.Err = nil
					c.Params["id"] = c.Job.Id()
					commands <- c.Copy()
				case architecture.RESERVE:
					sendChan := make(chan architecture.Command)
					watchedTubeConnectionsReceiver <- sendChan
					tube.AwaitingClients.Enqueue(architecture.NewAwaitingClient(c, sendChan))
				case architecture.RESERVE_WITH_TIMEOUT:
					log.Info("reserve-with-timeout", c)
					sendChan := make(chan architecture.Command)
					watchedTubeConnectionsReceiver <- sendChan
					client := architecture.NewAwaitingClient(c, sendChan)
					tube.AwaitingClients.Enqueue(client)
					tube.AwaitingTimedClients[client.Id()] = client
					tube.ProcessTimedClients()
				case architecture.DELETE:
					if tube.Buried.Delete(c.Params["id"]) != nil || tube.Reserved.Delete(c.Params["id"]) != nil {
						// log.Println("TUBE_HANDLER deleted job: ", c, name)
						c.Err = nil
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c.Copy()
				case architecture.RELEASE:
					item := tube.Reserved.Delete(c.Params["id"])
					if item != nil {
						job := item.(*architecture.Job)
						// log.Println("TUBE_HANDLER released job: ", c, name)
						job.SetState(architecture.READY)
						tube.Ready.Enqueue(job)
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c.Copy()
				case architecture.BURY:
					log.Info("buried ", c.Params["id"])
					item := tube.Reserved.Delete(c.Params["id"])
					if item != nil {
						job := item.(*architecture.Job)
						// log.Println("TUBE_HANDLER buried job: ", c, name)
						job.SetState(architecture.BURIED)
						tube.Buried.Enqueue(job)
						size := tube.Buried.Size()
						log.Info("buried size ", size)
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c.Copy()
				case architecture.KICK:
					amount := 0
					bound, err := strconv.Atoi(c.Params["bound"])
					log.Info("kicking laa", bound)
					if err != nil {
						// handle non-integer number of jobs to kick
					}
					size := tube.Buried.Size()
					log.Info("kicking buried size ", size)
					if size < bound {
						bound = size
					}
					for amount < bound {
						log.Info("kicking laa ", amount)
						item := tube.Buried.Dequeue()
						job := item.(*architecture.Job)
						job.SetState(architecture.READY)
						tube.Reserved.Enqueue(job)
						amount += 1
					}
					commands <- c.Copy()
				case architecture.KICK_JOB:
					item := tube.Buried.Delete(c.Params["id"])
					if item != nil {
						job := item.(*architecture.Job)
						job.SetState(architecture.READY)
						tube.Ready.Enqueue(job)
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c.Copy()
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

// createTube ensures that we can change the implementation data structure of the priority queue easily
// by changing only here
func createTube(name string) *architecture.Tube {
	t := &architecture.Tube{
		Name:                 name,
		Ready:                &backend.MinHeap{},
		Reserved:             &backend.MinHeap{},
		Delayed:              &backend.MinHeap{},
		Buried:               &backend.MinHeap{},
		AwaitingClients:      &backend.MinHeap{},
		AwaitingTimedClients: make(map[string]*architecture.AwaitingClient),
	}
	t.Ready.Init()
	t.Delayed.Init()
	t.Reserved.Init()
	t.Buried.Init()
	t.AwaitingClients.Init()
	return t
}
