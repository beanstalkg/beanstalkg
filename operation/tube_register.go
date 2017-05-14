package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
)

const DEFAULT_TUBE string = "default"

func NewTubeRegister(
	commands chan architecture.Command,
	useTubeConnectionReceiver chan chan architecture.Command,
	watchedTubeConnectionsReceiver chan chan architecture.Command,
	stop chan bool,
) {
	go func() {
		tubeStopChannels := make(map[string]chan bool)
		tubeChannels := make(map[string]chan architecture.Command)
		tubeChannels[DEFAULT_TUBE], tubeStopChannels[DEFAULT_TUBE] = createTubeHandler(DEFAULT_TUBE, watchedTubeConnectionsReceiver)
		for {
			select {
			case c := <-commands:
				switch c.Name {
				case architecture.USE:
					if _, ok := tubeChannels[c.Params["tube"]]; !ok {
						tubeChannels[c.Params["tube"]], tubeStopChannels[c.Params["tube"]] =
							createTubeHandler(c.Params["tube"], watchedTubeConnectionsReceiver)
					}
					useTubeConnectionReceiver <- tubeChannels[c.Params["tube"]]
					log.Debug("TUBE_REGISTER sent tube for use: ", c.Params["tube"])
				case architecture.WATCH:
					if _, ok := tubeChannels[c.Params["tube"]]; !ok {
						tubeChannels[c.Params["tube"]], tubeStopChannels[c.Params["tube"]] =
							createTubeHandler(c.Params["tube"], watchedTubeConnectionsReceiver)
					}
					useTubeConnectionReceiver <- tubeChannels[c.Params["tube"]]
					log.Debug("TUBE_REGISTER sent tube for watch: ", c.Params["tube"])
				}
			// TODO handle commands and send tubeChannels to clients if required
			case <-stop:
				// TODO send stop signal to all tube channels
				return
			}
		}
	}()

}

func createTubeHandler(
	name string,
	watchedTubeConnectionsReceiver chan chan architecture.Command,
) (
	chan architecture.Command,
	chan bool,
) {
	tubeChannel := make(chan architecture.Command)
	stop := make(chan bool)
	NewTubeHandler(name, tubeChannel, watchedTubeConnectionsReceiver, stop)
	return tubeChannel, stop
}
