package operation

import (
	"github.com/beanstalkg/beanstalkg/architecture"
)

const DEFAULT_TUBE string = "default"

type TubeRegister struct {
	tubeStopChannels               map[string]chan bool
	tubeChannels                   map[string]chan architecture.Command
	commands                       chan architecture.Command
	useTubeConnectionReceiver      chan chan architecture.Command
	watchedTubeConnectionsReceiver chan chan architecture.Command
	stop                           chan bool
}

func (tr *TubeRegister) init() {
	tr.tubeChannels[DEFAULT_TUBE], tr.tubeStopChannels[DEFAULT_TUBE] = tr.createTubeHandler(DEFAULT_TUBE,
		tr.watchedTubeConnectionsReceiver)
	for {
		select {
		case c := <-tr.commands:
			switch c.Name {
			case architecture.USE:
				fallthrough
			case architecture.WATCH:
				tr.createTubeIfNotExists(c.Tube)
				// send the tube connection to the client
				tr.useTubeConnectionReceiver <- tr.tubeChannels[c.Tube]
				log.Debugf("TUBE_REGISTER sent tube for %s: %s", c.Name, c.Tube)
			}
		case <-tr.stop:
			// TODO send stop signal to all tube channels
			return
		}
	}
}

func (tr *TubeRegister) createTubeIfNotExists(name string) {
	if _, ok := tr.tubeChannels[name]; !ok {
		// create tube_handler if does not exist
		tr.tubeChannels[name], tr.tubeStopChannels[name] =
			tr.createTubeHandler(name, tr.watchedTubeConnectionsReceiver)
	}
}

// createTubeHandler creates a new tube_handler with required command channel and stop channel
func (tr *TubeRegister) createTubeHandler(
	name string, watchedTubeConnectionsReceiver chan chan architecture.Command) (
	chan architecture.Command, chan bool) {
	tubeChannel := make(chan architecture.Command)
	stop := make(chan bool)
	NewTubeHandler(name, tubeChannel, watchedTubeConnectionsReceiver, stop)
	return tubeChannel, stop
}

func NewTubeRegister(
	commands chan architecture.Command,
	useTubeConnectionReceiver chan chan architecture.Command,
	watchedTubeConnectionsReceiver chan chan architecture.Command,
	stop chan bool,
) {
	// store the tube stop signalling channels
	tubeStopChannels := make(map[string]chan bool)
	// store the tube command sending channels
	tubeChannels := make(map[string]chan architecture.Command)
	tubeRegister := TubeRegister{
		tubeStopChannels,
		tubeChannels,
		commands,
		useTubeConnectionReceiver,
		watchedTubeConnectionsReceiver,
		stop,
	}
	go tubeRegister.init()
}
