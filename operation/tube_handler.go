package operation

import (
	"errors"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"github.com/vimukthi-git/beanstalkg/backend"
	// "log"
	"time"
)

func NewTubeHandler(
	name string,
	commands chan *architecture.Command,
	watchedTubeConnectionsReceiver chan chan *architecture.Command,
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
				// log.Debug("TUBE_HANDLER received: ", c)
				switch c.Name {
				case architecture.PUT:
					if c.Job.State() == architecture.READY {
						// log.Println("TUBE_HANDLER put job to ready queue: ", c, name)
						v := architecture.PriorityQueueItem(&c.Job)
						tube.Ready.Enqueue(&v)
					} else {
						// log.Println("TUBE_HANDLER put job to delayed queue: ", c, name)
						v := architecture.PriorityQueueItem(&c.Job)
						tube.Delayed.Enqueue(&v)
					}
					c.Err = nil
					c.Params["id"] = c.Job.Id()
					commands <- c
				case architecture.RESERVE:
					sendChan := make(chan *architecture.Command, 1)
					watchedTubeConnectionsReceiver <- sendChan
					v := architecture.PriorityQueueItem(architecture.NewAwaitingClient(c.Params["client_id"], *c, sendChan))
					tube.AwaitingClients.Enqueue(&v)
				case architecture.RESERVE_WITH_TIMEOUT:
					sendChan := make(chan *architecture.Command, 1)
					watchedTubeConnectionsReceiver <- sendChan
					client := architecture.NewAwaitingClient(c.Params["client_id"], *c, sendChan)
					v := architecture.PriorityQueueItem(client)
					tube.AwaitingClients.Enqueue(&v)
					tube.AwaitingTimedClients[client.Id()] = client
					tube.ProcessTimedClients()
				case architecture.DELETE:
					if tube.Buried.Delete(c.Params["id"]) != nil || tube.Reserved.Delete(c.Params["id"]) != nil {
						// log.Println("TUBE_HANDLER deleted job: ", c, name)
						c.Err = nil
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c
				case architecture.RELEASE:
					item := tube.Reserved.Delete(c.Params["id"])
					if item != nil {
						job := item.(*architecture.Job)
						// log.Println("TUBE_HANDLER released job: ", c, name)
						job.SetState(architecture.READY)
						v := architecture.PriorityQueueItem(job)
						tube.Ready.Enqueue(&v)
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c
				case architecture.BURY:
					item := tube.Reserved.Delete(c.Params["id"])
					if item != nil {
						job := item.(*architecture.Job)
						// log.Println("TUBE_HANDLER buried job: ", c, name)
						job.SetState(architecture.BURIED)
						v := architecture.PriorityQueueItem(job)
						tube.Buried.Enqueue(&v)
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c
				case architecture.INTERNAL_CLOSE_CLIENT:
					item := tube.AwaitingClients.Find(c.Params["client_id"])
					// log.Info("TUBE_HANDLER close ", c.Params["client_id"])
					if item != nil {
						client := item.(*architecture.AwaitingClient)
						close(client.SendChannel)
					}
					tube.AwaitingClients.Delete(c.Params["client_id"])
					delete(tube.AwaitingTimedClients, c.Params["client_id"])
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
		Name: name,
		Ready: &backend.MinHeap{},
		Reserved: &backend.MinHeap{},
		Delayed: &backend.MinHeap{},
		Buried: &backend.MinHeap{},
		AwaitingClients: &backend.MinHeap{},
		AwaitingTimedClients: make(map[string]*architecture.AwaitingClient),
	}
	t.Ready.Init()
	t.Delayed.Init()
	t.Reserved.Init()
	t.Buried.Init()
	t.AwaitingClients.Init()
	return t
}