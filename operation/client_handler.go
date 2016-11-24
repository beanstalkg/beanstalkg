package operation

import (
	"bufio"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
	"net"
	"errors"
	"strconv"
	"fmt"
)

type clientHandler struct {
	conn                           net.Conn
	registerConnection             chan architecture.Command
	tubeConnectionReceiver         chan chan architecture.Command
	usedTubeConnection             chan architecture.Command
	watchedTubeConnectionsReceiver chan chan architecture.Command
	watchedTubeConnections         map[string]chan architecture.Command
	stop                           chan bool
}

func NewClientHandler(
	conn net.Conn,
	registerConnection chan architecture.Command,
	useTubeConnectionReceiver chan chan architecture.Command,
	watchedTubeConnectionsReceiver chan chan architecture.Command,
	stop chan bool,
) {
	go func() {
		defer conn.Close()
		client := clientHandler{
			conn,
			registerConnection,
			useTubeConnectionReceiver,
			nil,
			watchedTubeConnectionsReceiver,
			nil,
			stop,
		}
		client.startSession()
	}()
}

func (client *clientHandler) handleReply(c architecture.Command) error {
	for {
		more, reply := c.Reply()
		_, err := client.conn.Write([]byte(reply + "\r\n"))
		if err != nil {
			log.Print(err)
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
	client.registerConnection <- c
	client.usedTubeConnection = <-client.tubeConnectionReceiver
	client.watchedTubeConnections = map[string]chan architecture.Command{
		"default": client.usedTubeConnection,
	}
	// convert scan to a selectable
	scan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(client.conn)
		for scanner.Scan() {
			scan <- scanner.Text()
		}
	}()

	for {
		select {
		case rawCommand := <-scan:
			parsed, err := c.Parse(rawCommand)
			if err != nil { // check if parse error
				err = client.handleReply(c)
				c = architecture.Command{}
				if err != nil {
					return
				}
			} else if parsed { // check if the command has been parsed completely
				c = client.handleCommand(c)
				err = client.handleReply(c)
				if err != nil {
					return
				}
				// we replace previous command once its parsing is finished
				c = architecture.Command{}
			}
		case <-client.stop:
			return
		}
	}
}

func (client *clientHandler) handleCommand(command architecture.Command) architecture.Command {
	switch command.Name {
	case architecture.USE:
		// send command to tube register
		client.registerConnection <- command
		client.usedTubeConnection = <- client.tubeConnectionReceiver
		log.Println("CLIENT_HANDLER started using tube: ", command.Params["tube"])
	case architecture.PUT:
		client.usedTubeConnection <- command  // send the command to tube
		command = <-client.usedTubeConnection // get the response
	case architecture.WATCH:
		client.registerConnection <- command
		client.watchedTubeConnections[command.Params["tube"]] = <- client.tubeConnectionReceiver
		command.Params["count"] = strconv.FormatInt(int64(len(client.watchedTubeConnections)), 10)
	case architecture.IGNORE:
		if _, ok := client.watchedTubeConnections[command.Params["tube"]]; ok && len(client.watchedTubeConnections) > 1 {
			delete(client.watchedTubeConnections, command.Params["tube"])
			command.Params["count"] = strconv.FormatInt(int64(len(client.watchedTubeConnections)), 10)
		} else {
			command.Err = errors.New(architecture.NOT_IGNORED)
		}
	case architecture.RESERVE:
		recv := make(chan architecture.Command)
		go func() {
			receiveConnections := []chan architecture.Command{}
			for _, connection := range client.watchedTubeConnections {
				connection <- command
				append(receiveConnections, <-client.watchedTubeConnectionsReceiver)
			}


		}()
	case architecture.RESERVE_WITH_TIMEOUT:
	case architecture.RELEASE:
	case architecture.BURY:
	case architecture.TOUCH:

	}

	return command
}
