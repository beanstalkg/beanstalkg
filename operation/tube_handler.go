package operation

import (
	"github.com/beanstalkg/beanstalkg/architecture"
	"github.com/beanstalkg/beanstalkg/backend"
	"time"
)

func NewTubeHandler(
	name string,
	commands chan architecture.Command,
	watchedTubeConnectionsReceiver chan chan architecture.Command,
	stop chan bool,
) {
	go func() {
		// create the tube
		tube := architecture.NewTube(name, func() architecture.PriorityQueue {
			// plug backends here
			return &backend.MinHeap{}
		})
		ticker := time.NewTicker(architecture.QUEUE_FREQUENCY)
		for {
			select {
			case <-ticker.C:
				log.Debug("House Keeping Started for: ", name)
				tube.Process()
				tube.ProcessTimedClients()
			case incomingCommand := <-commands:
				switch incomingCommand.Name {
				case architecture.PUT:
					tube.Put(&incomingCommand)
					commands <- incomingCommand.Copy()
				case architecture.RESERVE:
					sendChan := make(chan architecture.Command)
					watchedTubeConnectionsReceiver <- sendChan
					tube.Reserve(&incomingCommand, sendChan)
				case architecture.RESERVE_WITH_TIMEOUT:
					log.Info("reserve-with-timeout", incomingCommand)
					sendChan := make(chan architecture.Command)
					watchedTubeConnectionsReceiver <- sendChan
					tube.ReserveWithTimeout(&incomingCommand, sendChan)
				case architecture.DELETE:
					tube.Delete(&incomingCommand)
					commands <- incomingCommand.Copy()
				case architecture.RELEASE:
					tube.Release(&incomingCommand)
					commands <- incomingCommand.Copy()
				case architecture.BURY:
					tube.Bury(&incomingCommand)
					commands <- incomingCommand.Copy()
				case architecture.KICK:
					tube.Kick(&incomingCommand)
					commands <- incomingCommand.Copy()
				case architecture.KICK_JOB:
					tube.KickJob(&incomingCommand)
					commands <- incomingCommand.Copy()
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}
