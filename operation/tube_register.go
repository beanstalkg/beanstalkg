package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
)

const DEFAULT_TUBE string = "default"

func NewTubeRegister(commands chan architecture.Command, handlerConnections chan chan architecture.Command, stop chan bool) {
	go func() {
		tubeStopChannels := make(map[string]chan bool)
		tubeChannels := make(map[string]chan architecture.Command)
		tubeChannels[DEFAULT_TUBE], tubeStopChannels[DEFAULT_TUBE] = createTubeHandler(DEFAULT_TUBE)
		for {
			select {
			case c := <-commands:
				switch c.Name {
				case architecture.USE:
					if _, ok := tubeChannels[c.Params["tube"]]; !ok {
						tubeChannels[c.Params["tube"]], tubeStopChannels[c.Params["tube"]] =
							createTubeHandler(c.Params["tube"])
					}
					handlerConnections<-tubeChannels[c.Params["tube"]]
					log.Println("TUBE_REGISTER sent tube: ", c.Params["tube"])
				}
			// TODO handle commands and send tubeChannels to clients if required
			case <-stop:
				return
			}
		}
	}()

}

func createTubeHandler(name string) (chan architecture.Command, chan bool) {
	tubeChannel := make(chan architecture.Command)
	stop := make(chan bool)
	NewTubeHandler(name, tubeChannel, stop)
	return tubeChannel, stop
}

