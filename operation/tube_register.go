package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
)

const DEFAULT_TUBE string = "default"

func NewTubeRegister(
	commands chan architecture.Command,
	tubeConnections chan chan architecture.Command,
	jobConnections chan chan architecture.Job,
	stop chan bool,
) {
	go func() {
		tubeStopChannels := make(map[string]chan bool)
		tubeChannels := make(map[string]chan architecture.Command)
		tubeChannels[DEFAULT_TUBE], tubeStopChannels[DEFAULT_TUBE] = createTubeHandler(DEFAULT_TUBE, jobConnections)
		for {
			select {
			case c := <-commands:
				switch c.Name {
				case architecture.USE:
					if _, ok := tubeChannels[c.Params["tube"]]; !ok {
						tubeChannels[c.Params["tube"]], tubeStopChannels[c.Params["tube"]] =
							createTubeHandler(c.Params["tube"], jobConnections)
					}
					tubeConnections <- tubeChannels[c.Params["tube"]]
					log.Println("TUBE_REGISTER sent tube: ", c.Params["tube"])
				}
			// TODO handle commands and send tubeChannels to clients if required
			case <-stop:
				return
			}
		}
	}()

}

func createTubeHandler(
	name string,
	jobConnections chan chan architecture.Job,
) (
	chan architecture.Command,
	chan bool,
) {
	tubeChannel := make(chan architecture.Command)
	stop := make(chan bool)
	NewTubeHandler(name, tubeChannel, jobConnections, stop)
	return tubeChannel, stop
}
