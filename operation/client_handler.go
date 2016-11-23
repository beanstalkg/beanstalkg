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
		scanner := bufio.NewScanner(conn)
		log.Println("Scanning.. ")
		// this command object will be replaced each time the client sends a new one
		c := architecture.NewDefaultCommand()
		var tubeConnection chan architecture.Command
		// selects default tube first up
		registerConnection <- c
		tubeConnection = <-tubeConnections

		// convert scan to a selectable
		scan := make(chan string)
		go func() {
			for scanner.Scan() {
				scan <- scanner.Text()
			}
		}()

		for {
			select {
			case rawCommand := <-scan:
				parsed, err := c.Parse(rawCommand)
				if err != nil { // check if parse error
					err = handleReply(conn, c)
					c = architecture.Command{}
					if err != nil {
						return
					}
				} else if parsed { // check if the command has been parsed completely
					c = handleCommand(
						c,
						registerConnection,
						tubeConnections,
						&tubeConnection,
					)
					err = handleReply(conn, c)
					if err != nil {
						return
					}
					// we replace previous command once its parsing is finished
					c = architecture.Command{}
				}
				log.Println("Scanning.. ")
			case <-stop:
				return
			}
		}
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

func handleCommand(
	command architecture.Command,
	registerConnection chan architecture.Command,
	tubeConnections chan chan architecture.Command,
	tubeConnection *chan architecture.Command,
) architecture.Command {
	switch command.Name {
	case architecture.USE:
		// send command to tube register
		registerConnection <- command
		tubeConnectionTemp := <-tubeConnections
		tubeConnection = &tubeConnectionTemp
		log.Println("CLIENT_HANDLER started using tube: ", command.Params["tube"])
	case architecture.PUT:
		*tubeConnection <- command  // send the command to tube
		command = <-*tubeConnection // get the response
	}
	return command
}
