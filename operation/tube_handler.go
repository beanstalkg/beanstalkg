package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"github.com/vimukthi-git/beanstalkg/backend"
	"log"
	"time"
	"errors"
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
		ticker := time.NewTicker(1 * time.Second)
		// TODO make sure all logic that can be moved to Tube struct is moved there
		for {
			select {
			case <-ticker.C:
				// log.Println("House Keeping Started for: ", name)
				tube.Process()
			case c := <-commands:
				switch c.Name {
				case architecture.PUT:
					if c.Job.State() == architecture.READY {
						log.Println("TUBE_HANDLER put job to ready queue: ", c, name)
						tube.Ready.Enqueue(&c.Job)
					} else {
						log.Println("TUBE_HANDLER put job to delayed queue: ", c, name)
						tube.Delayed.Enqueue(&c.Job)
					}
					c.Err = nil
					c.Params["id"] = c.Job.Id()
					commands <- c.Copy()
				case architecture.RESERVE:
					sendChan := make(chan architecture.Command)
					tube.AwaitingClients.Enqueue(architecture.NewAwaitingClient(c, sendChan))
					watchedTubeConnectionsReceiver <- sendChan
				case architecture.DELETE:
					if tube.Buried.Delete(c.Params["id"]) != nil || tube.Reserved.Delete(c.Params["id"]) != nil {
						log.Println("TUBE_HANDLER deleted job: ", c, name)
						c.Err = nil
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c.Copy()
				case architecture.RELEASE:
					job := tube.Reserved.Delete(c.Params["id"]).(*architecture.Job)
					if job != nil {
						log.Println("TUBE_HANDLER released job: ", c, name)
						job.SetState(architecture.READY)
						tube.Ready.Enqueue(job)
					} else {
						c.Err = errors.New(architecture.NOT_FOUND)
					}
					commands <- c.Copy()
				case architecture.BURY:
					job := tube.Reserved.Delete(c.Params["id"]).(*architecture.Job)
					if job != nil {
						log.Println("TUBE_HANDLER buried job: ", c, name)
						job.SetState(architecture.BURIED)
						tube.Buried.Enqueue(job)
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
		name,
		&backend.MinHeap{},
		&backend.MinHeap{},
		&backend.MinHeap{},
		&backend.MinHeap{},
		&backend.MinHeap{},
	}
	t.Ready.Init()
	t.Delayed.Init()
	t.Reserved.Init()
	t.Buried.Init()
	t.AwaitingClients.Init()
	return t
}
