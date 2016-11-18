package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"github.com/vimukthi-git/beanstalkg/server"
)

func main() {
	service := ":11300"
	me := server.Beanstalkg{}
	me.Init()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		fmt.Println("Waiting..")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, &me)
	}
}

func handleClient(conn net.Conn, me *server.Beanstalkg) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	fmt.Println("Scanning.. ")
	// this command object will be replaced each time the client sends a new one
	c := server.Command{}
	// scan continuously for client commands
	for scanner.Scan() {
		rawCommand := scanner.Text()
		parsed, err := c.Parse(rawCommand)
		if err != nil {
			return
		}
		if parsed {
			_, err2 := conn.Write([]byte(me.ExecCommand(c) + "\r\n"))
			if err2 != nil {
				return
			}
			// fmt.Println(c)
			// we replace previous command once its parsing is finished
			c = server.Command{}
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
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}