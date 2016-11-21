package operation

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"net"
	"bufio"
	"fmt"
	"os"
)

func NewClientHandler(conn net.Conn, registerConnection chan architecture.Command,
	tubeConnections chan chan architecture.Command, stop chan bool) {
	go func() {
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		fmt.Println("Scanning.. ")
		// this command object will be replaced each time the client sends a new one
		c := architecture.Command{}
		// this map contains information regarding client connection
		// selects default tube first up
		context := make(map[string]string)
		context["tube"] = DEFAULT_TUBE
		// scan continuously for client commands
		for scanner.Scan() {
			rawCommand := scanner.Text()
			parsed, err := c.Parse(rawCommand)
			if err != nil {
				return
			}
			if parsed {
				registerConnection <- c
				select {
				case c2 := <-registerConnection:
					// TODO
				case con := <-tubeConnections:
					// TODO write on received channel to interact with the tube
				}
				_, err2 := conn.Write([]byte(ExecCommand(c, context) + "\r\n"))
				if err2 != nil {
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
			fmt.Println("Scanning.. ")
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()
}