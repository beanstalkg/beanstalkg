package main

import (
	"encoding/json"
	"flag"
	"github.com/op/go-logging"
	"github.com/beanstalkg/beanstalkg/pkg/architecture"
	"github.com/beanstalkg/beanstalkg/pkg/operation"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var log = logging.MustGetLogger("BEANSTALKG")

func main() {
	port := flag.String("port", "11300", "Port for beanstalkg server")
	proxyMode := flag.Bool("proxy_mode", false, "Start server in proxy mode")
	env := flag.String("env", "local", "Which environment config to use")
	debugMode := flag.Bool("debug_mode", false, "Start server in debug mode. Logs shall contain more information")
	flag.Parse()
	initLogging(*debugMode)
	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	service := ":" + *port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	stop := make(chan bool)

	if !*proxyMode {
		tubeRegister := make(chan architecture.Command)
		// use this tube to send the channels for each individual tube to the clients when the do 'use' command
		useTubeConnectionReceiver := make(chan chan architecture.Command)
		watchedTubeConnectionsReceiver := make(chan chan architecture.Command)
		operation.NewTubeRegister(tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop)
		log.Info("BEANSTALKG listening on: ", *port)

		for {
			// log.Println("BEANSTALKG Waiting..")
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			operation.NewClientHandler(conn, tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop)
		}
	} else {
		config := getConfig(*env)
		log.Info("BEANSTALKG started in proxy mode, now listening on: ", *port)
		for {
			// log.Println("BEANSTALKG Waiting..")
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			operation.NewProxiedClientHandler(conn, config.Beanstalks, stop)
		}

	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error:", err.Error())
	}
}

type Configuration struct {
	Beanstalks []string `json:"beanstalks"`
}

func getConfig(env string) Configuration {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := make(map[string]Configuration)
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("error in parsing config:", err)
	}
	envConf, ok := configuration[env]
	if !ok {
		log.Fatal("No configuration found for the given environment name")
	}
	return envConf
}

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func initLogging(debugMode bool) {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stdout, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.ERROR, "")

	backend2Leveled := logging.AddModuleLevel(backend2Formatter)
	if debugMode {
		backend2Leveled.SetLevel(logging.DEBUG, "")
	} else {
		backend2Leveled.SetLevel(logging.INFO, "")
	}
	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Leveled)
}
