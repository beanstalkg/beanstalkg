package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"net"
	"bufio"
	"log"
)

func NewClientHandler(conn net.Conn, registerConnection chan architecture.Command,
	tubeConnections chan chan architecture.Command, stop chan bool) {
	go func() {
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		log.Println("Scanning.. ")
		// this command object will be replaced each time the client sends a new one
		c := architecture.NewDefaultCommand()
		var tubeConnection chan architecture.Command
		// selects default tube first up
		registerConnection <- c
		tubeConnection = <- tubeConnections

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
				if err != nil {
					return
				}
			// check if the command has been parsed completely
				if parsed {
					var err error
					c, err = handleCommand(
						c,
						registerConnection,
						tubeConnections,
						&tubeConnection,
					)
					if err != nil {
						log.Print(err)
						return
					}

					_, err2 := conn.Write([]byte(c.Reply() + "\r\n"))
					if err2 != nil {
						log.Print(err2)
						return
					}
					// fmt.Println(c)
					// we replace previous command once its parsing is finished
					c = architecture.Command{}
				}
			//_, err2 := conn.Write([]byte(rawCommand + "\r\n"))
			//if err2 != nil {
			//	return
			//}
				log.Println("Scanning.. ")
			case <-stop:
				return
			}
		}
	}()
}

func handleCommand(
			command architecture.Command,
			registerConnection chan architecture.Command,
			tubeConnections chan chan architecture.Command,
			tubeConnection *chan architecture.Command,
		) (architecture.Command, error) {
	switch command.Name {
	case architecture.USE:
		// send command to tube register
		registerConnection <- command
		tubeConnectionTemp := <- tubeConnections
		tubeConnection = &tubeConnectionTemp
		log.Println("CLIENT_HANDLER started using tube: ", command.Params["tube"])
	case architecture.PUT:

	}
	return command, nil
}