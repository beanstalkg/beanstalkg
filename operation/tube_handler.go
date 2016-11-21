package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"github.com/vimukthi-git/beanstalkg/backend"
	"time"
)

func NewTubeHandler(name string, commands chan architecture.Command, stop chan bool) {
	// commands := make(chan architecture.Command)
	go func() {
		tube := createTube(name)
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				tube.Delayed.Init()
				// TODO house keeping - check if any delayed jobs are ready, reserved jobs are ready
			case c := <-commands:
				switch c.Name {
				case architecture.PUT:

				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func createTube(name string) *architecture.Tube {
	return architecture.Tube{
		name,
		&backend.MinHeap{},
		&backend.MinHeap{},
		&backend.MinHeap{},
		make([]architecture.PriorityQueueItem, 100),
	}
}