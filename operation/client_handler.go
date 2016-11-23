package operation

import (
	"bufio"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
	"net"
)

func NewClientHandler(
	conn net.Conn,
	registerConnection chan architecture.Command,
	tubeConnections chan chan architecture.Command,
	jobConnections chan chan architecture.Job,
	stop chan bool,
) {
	go func() {
		defer conn.Close()

		client := clientHandler{
			conn,
			registerConnection,
			tubeConnections,
			nil,
			jobConnections,
			stop,
		}
		client.startSession()
	}()
}

func handleReply(conn net.Conn, c architecture.Command) error {
	_, err := conn.Write([]byte(c.Reply() + "\r\n"))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

type clientHandler struct {
	conn net.Conn
	registerConnection chan architecture.Command
	tubeConnections chan chan architecture.Command
	currentTubeConnection chan architecture.Command
	jobConnections chan chan architecture.Job
	stop chan bool
}

func (client *clientHandler) startSession() {
	// this command object will be replaced each time the client sends a new one
	c := architecture.NewDefaultCommand()
	// selects default tube first up
	client.registerConnection <- c
	client.currentTubeConnection = <-client.tubeConnections

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
				err = handleReply(client.conn, c)
				c = architecture.Command{}
				if err != nil {
					return
				}
			} else if parsed { // check if the command has been parsed completely
				c = client.handleCommand(c)
				err = handleReply(client.conn, c)
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
		client.currentTubeConnection = <-client.tubeConnections
		log.Println("CLIENT_HANDLER started using tube: ", command.Params["tube"])
	case architecture.PUT:
		client.currentTubeConnection <- command  // send the command to tube
		command = <-client.currentTubeConnection // get the response
	}
	return command
}
