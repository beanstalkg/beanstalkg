package operation

import (
	"bufio"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"net"
	"reflect"
)

type clientProxiedHandler struct {
	conn           net.Conn
	proxiedServers []chan string
	serverStatus   []bool
	stop           chan bool
	error          chan int
}

func NewProxiedClientHandler(
	conn net.Conn,
	proxiedServerNames []string,
	stop chan bool,
) {
	go func() {
		defer conn.Close()
		proxiedServers := []chan string{}
		errorChannel := make(chan int)
		for index, server := range proxiedServerNames {
			proxiedServers = append(proxiedServers, createProxyServerHandler(index, server, stop, errorChannel))
		}
		client := clientProxiedHandler{
			conn:           conn,
			proxiedServers: proxiedServers,
			stop:           stop,
			error:          errorChannel,
		}
		client.startSession()
		log.Info("PROXIED_CLIENT_HANDLER exit")
		return
	}()
}

func (client *clientProxiedHandler) startSession() {
	// convert scan to a selectable
	scan := make(chan string)
	exit := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(client.conn)
		c := architecture.NewCommand()
		for scanner.Scan() {
			done, _ := c.Parse(scanner.Text())
			if done {
				scan <- c.RawCommand
				c = architecture.NewCommand()
			}
		}
		exit <- true
	}()

	for {
		select {
		case rawCommand := <-scan:
			log.Info("PROXIED_CLIENT_HANDLER", rawCommand)
			serverCount := 0
			for _, proxyHandler := range client.proxiedServers {
				proxyHandler <- rawCommand
				serverCount++
			}
			var chosenReply string
			for i := 0; i < serverCount; i++ {
				cases := make([]reflect.SelectCase, len(client.proxiedServers))
				for i, ch := range client.proxiedServers {
					cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
				}
				chosen, value, _ := reflect.Select(cases)
				choose := 0 // TODO configurable
				if chosen == choose {
					chosenReply = value.Interface().(string)
				}
			}
			_, err := client.conn.Write([]byte(chosenReply + "\r\n"))
			if err != nil {
				log.Error(err)
				return
			}
		case <-client.stop:
			return
		case id := <-client.error:
			log.Info("PROXIED_CLIENT_HANDLER error from: ", id, len(client.proxiedServers))
			client.proxiedServers[id] = client.proxiedServers[len(client.proxiedServers)-1] // Replace id with the last one.
			client.proxiedServers = client.proxiedServers[:len(client.proxiedServers)-1]    // Chop off the last one.
		case <-exit:
			return
		}
	}
}

func createProxyServerHandler(id int, url string, stop chan bool, error chan int) chan string {
	com := make(chan string)
	go func() {
		rAddr, err := net.ResolveTCPAddr("tcp4", url)
		if err != nil {
			log.Error(err)
			error <- id
			return
		}

		rConn, err := net.DialTCP("tcp", nil, rAddr)
		if err != nil {
			log.Error(err)
			error <- id
			return
		}
		defer rConn.Close()

		scan := make(chan string)
		go func() {
			scanner := bufio.NewScanner(rConn)
			for scanner.Scan() {
				scan <- scanner.Text()
			}
			error <- id
		}()

		for {
			select {
			case command := <-com:
				_, err := rConn.Write([]byte(command + "\r\n"))
				if err != nil {
					log.Error(err)
					error <- id
					return
				}
			case reply := <-scan:
				com <- reply
			case <-stop:
				return
			}
		}
	}()
	return com
}
