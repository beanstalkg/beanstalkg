package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"github.com/vimukthi-git/beanstalkg/backend"
	"log"
	"time"
)

func NewTubeHandler(
	name string,
	commands chan architecture.Command,
	jobConnections chan chan architecture.Job,
	stop chan bool,
) {
	// commands := make(chan architecture.Command)
	go func() {
		tube := createTube(name)
		ticker := time.NewTicker(1 * time.Second)
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
						tube.Ready.Enqueue(c.Job)
					} else {
						log.Println("TUBE_HANDLER put job to delayed queue: ", c, name)
						tube.Delayed.Enqueue(c.Job)
					}
					c.Err = nil
					c.Params["id"] = c.Job.Id()
					commands <- c
				case architecture.RESERVE:

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
