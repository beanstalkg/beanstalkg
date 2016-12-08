package operation

import (
	"bufio"
	"errors"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"net"
	"reflect"
	"strconv"
	// "github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

// var log = logging.MustGetLogger("BEANSTALKG")

type clientHandler struct {
	id			       string
	conn                           net.Conn
	registerConnection             chan *architecture.Command
	tubeConnectionReceiver         chan chan *architecture.Command
	usedTubeConnection             chan *architecture.Command
	watchedTubeConnectionsReceiver chan chan *architecture.Command
	watchedTubeConnections         map[string]chan *architecture.Command
	reservedJobs                   map[string]string
	stop                           chan bool
}

func NewClientHandler(
	conn net.Conn,
	registerConnection chan *architecture.Command,
	useTubeConnectionReceiver chan chan *architecture.Command,
	watchedTubeConnectionsReceiver chan chan *architecture.Command,
	stop chan bool,
) {
	go func() {
		defer conn.Close()
		client := clientHandler{
			uuid.NewV1().String(),
			conn,
			registerConnection,
			useTubeConnectionReceiver,
			nil,
			watchedTubeConnectionsReceiver,
			nil,
			map[string]string{},
			stop,
		}
		// log.Info("CLIENT_HANDLER start ", client.id)
		client.startSession()
		// log.Info("CLIENT_HANDLER exit ", client.id)
		return
	}()
}

func (client *clientHandler) handleReply(c architecture.Command) error {
	for {
		more, reply := c.Reply()
		_, err := client.conn.Write([]byte(reply + "\r\n"))
		if err != nil {
			// log.Error(err)
			return err
		}
		if !more {
			break
		}
	}
	return nil
}

func (client *clientHandler) startSession() {
	// this command object will be replaced each time the client sends a new one
	c := architecture.NewDefaultCommand()
	// selects default tube first up
	cCopy := c.Copy()
	client.registerConnection <- &cCopy
	client.usedTubeConnection = <-client.tubeConnectionReceiver
	client.watchedTubeConnections = map[string]chan *architecture.Command{
		"default": client.usedTubeConnection,
	}
	// convert scan to a selectable
	scan := make(chan string)
	exit := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(client.conn)
		for scanner.Scan() {
			scan <- scanner.Text()
		}
		// log.Info("******** DROPPED **** ", client.id)
		exit <- true
	}()

	for {
		select {
		case rawCommand := <-scan:
			parsed, err := c.Parse(rawCommand)
			if err != nil { // check if parse error
				err = client.handleReply(c)
				c = architecture.NewCommand()
				if err != nil {
					return
				}
			} else if parsed { // check if the command has been parsed completely
				c.Params["client_id"] = client.id
				c = client.handleBasicCommand(c)
				err = client.handleReply(c)
				if err != nil {
					return
				}
				// we replace previous command once its parsing is finished
				c = architecture.NewCommand()
			}
		case <-client.stop:
			client.close()
			return
		case <-exit:
			client.close()
			return
		}
	}
}

func (client *clientHandler) handleBasicCommand(command architecture.Command) architecture.Command {
	// log.Debug("CLIENT_HANDLER handling command ", command.Params["client_id"], command)
	cCopy := command.Copy()
	switch command.Name {
	case architecture.USE:
		// send command to tube register
		client.registerConnection <- &cCopy
		client.usedTubeConnection = <-client.tubeConnectionReceiver
	case architecture.PUT:
		client.usedTubeConnection <- &cCopy // send the command to tube
		command = *(<-client.usedTubeConnection)       // get the response
	case architecture.WATCH:
		client.registerConnection <- &cCopy
		client.watchedTubeConnections[command.Params["tube"]] = <-client.tubeConnectionReceiver
		command.Params["count"] = strconv.FormatInt(int64(len(client.watchedTubeConnections)), 10)
	case architecture.IGNORE:
		if _, ok := client.watchedTubeConnections[command.Params["tube"]]; ok && len(client.watchedTubeConnections) > 1 {
			delete(client.watchedTubeConnections, command.Params["tube"])
			command.Params["count"] = strconv.FormatInt(int64(len(client.watchedTubeConnections)), 10)
		} else {
			command.Err = errors.New(architecture.NOT_IGNORED)
		}
	case architecture.RESERVE:
		command = client.reserve(command)
	case architecture.RESERVE_WITH_TIMEOUT: // TODO add a new queue to each tube to track timeout client ids to handle this
		command = client.reserve(command)
	case architecture.DELETE:
		if tube, ok := client.reservedJobs[command.Params["id"]]; ok {
			if con, ok := client.watchedTubeConnections[tube]; ok {
				con <- &cCopy
				command = *(<-con)
			}
		} else {
			command.Err = errors.New(architecture.NOT_FOUND)
		}
	case architecture.RELEASE:
		if tube, ok := client.reservedJobs[command.Params["id"]]; ok {
			if con, ok := client.watchedTubeConnections[tube]; ok {
				con <- &cCopy
				command = *(<-con)
			}
		} else {
			command.Err = errors.New(architecture.NOT_FOUND)
		}
	case architecture.BURY:
		if tube, ok := client.reservedJobs[command.Params["id"]]; ok {
			if con, ok := client.watchedTubeConnections[tube]; ok {
				con <- &cCopy
				command = *(<-con)
			}
		} else {
			command.Err = errors.New(architecture.NOT_FOUND)
		}
	case architecture.TOUCH:

	}
	return command
}

func (client *clientHandler) reserve(command architecture.Command) architecture.Command {
	recv := make(chan *architecture.Command)
	go func() {
		// iterate and create a list of watched connections to receive from
		receiveConnections := []chan *architecture.Command{}
		receiveConnectionNames := []string{}
		for name, connection := range client.watchedTubeConnections {
			connection <- &command
			receiveConnections = append(receiveConnections, <-client.watchedTubeConnectionsReceiver)
			receiveConnectionNames = append(receiveConnectionNames, name)
		}
		// receive from one of the channels
		cases := make([]reflect.SelectCase, len(receiveConnections))
		for i, ch := range receiveConnections {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}
		// log.Debug("CLIENT_HANDLER waiting on reserve client_id: ", c.Params["client_id"])
		chosen, value, _ := reflect.Select(cases)
		resultCommand := value.Interface().(*architecture.Command)
		resultCommand.Params["tube"] = receiveConnectionNames[chosen]
		// log.Debug("CLIENT_HANDLER selected on reserve client_id: ", c.Params["client_id"], resultCommand)
		recv <- resultCommand
		// log.Debug("CLIENT_HANDLER sent reserved client_id: ", c.Params["client_id"], resultCommand)
		return
	}()
	// log.Debug("CLIENT_HANDLER **waiting on receive reserved client_id: ", command.Params["client_id"])
	newCommand := <-recv
	close(recv)
	if (*newCommand).Err == nil {
		client.reservedJobs[(*newCommand).Job.Id()] = (*newCommand).Params["tube"]
	}
	return *newCommand
}

func (client *clientHandler) close() {
	for _, connection := range client.watchedTubeConnections {
		c := architecture.Command{Name:architecture.INTERNAL_CLOSE_CLIENT, Params:map[string]string{"client_id": client.id}}
		connection <- &c
	}
}
