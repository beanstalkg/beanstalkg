package operation

import (
	"bufio"
	"errors"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
	"net"
	"reflect"
	"strconv"
	"net/url"
)

type clientProxiedHandler struct {
	conn                           net.Conn
	proxiedServers 			[]chan string
	stop                           chan bool
}

func NewProxiedClientHandler(
	conn net.Conn,
	proxiedServerNames []string,
	stop chan bool,
) {
	go func() {
		defer conn.Close()
		proxiedServers := []chan string{}
		for _, server := range proxiedServerNames {
			proxiedServers = append(proxiedServers, createProxyServerHandler(server, stop))
		}
		client := clientProxiedHandler{
			conn,
			proxiedServers,
			stop,
		}
		client.startSession()
	}()
}

func (client *clientProxiedHandler) handleReply(c architecture.Command) error {
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

func (client *clientProxiedHandler) startSession() {
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

		case <-client.stop:
			return
		}
	}
}

func createProxyServerHandler(url string, stop chan bool) chan string {
	com := make(chan string)
	go func() {
		rAddr, err := net.ResolveTCPAddr("tcp4", url)
		if err != nil {
			panic(err)
		}

		rConn, err := net.DialTCP("tcp", nil, rAddr)
		if err != nil {
			panic(err)
		}
		defer rConn.Close()

		scan := make(chan string)
		go func() {
			scanner := bufio.NewScanner(rConn)
			for scanner.Scan() {
				scan <- scanner.Text()
			}
		}()

		for {
			select {
			case command := <- com:
				_, err := rConn.Write([]byte(command + "\r\n"))
				if err != nil {
					log.Print(err)
					return err
				}
			case reply := <- scan:
				com <- reply
			case <- stop:
				return
			}
		}
	}()
	return com
}

