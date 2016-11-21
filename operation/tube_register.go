package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
)

const DEFAULT_TUBE string = "default"

func NewTubeRegister(commands chan architecture.Command, stop chan bool) {
	go func() {
		tubeStopChannels := make(map[string]chan bool)
		tubeChannels := make(map[string]chan chan architecture.Command)
		tubeChannels[DEFAULT_TUBE], tubeStopChannels[DEFAULT_TUBE] = createTubeHandler(DEFAULT_TUBE)
		for {
			select {
			case c := <-commands:
				switch c.Name {
				case architecture.USE:
					if _, ok := b.tubes[c.Params["tube"]]; !ok {
						b.tubes[c.Params["tube"]] = createTube(c.Params["tube"])
					}
					context["tube"] = c.Params["tube"]
					return c.Reply()
				}
			// TODO handle commands and send tubeChannels to clients if required
			case <-stop:
				return
			}
		}
	}()

}

func createTubeHandler(name string) chan chan architecture.Command {
	tubeChannel := make(chan chan architecture.Command)
	stop := make(chan bool)
	NewTubeHandler(name, tubeChannel, stop)
	return tubeChannel, stop
}

