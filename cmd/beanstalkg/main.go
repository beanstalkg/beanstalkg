package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"github.com/beanstalkg/beanstalkg/pkg/architecture"
	"github.com/beanstalkg/beanstalkg/pkg/operation"
	"github.com/beanstalkg/beanstalkg/pkg/backend"
	// planning to do this, but seems not working
	//"github.com/beanstalkg/beanstalkg/config/"
	"github.com/op/go-logging"

	// temporarily here to merge "new folder struct"
	"flag"
	"github.com/jinzhu/configor"
)

// temporarily here to merge "new folder struct"
type ServerConfig struct {
	Port    string `default:"11300"`
	Debug   bool
	Backend string `default:"minheap"`
}

var log = logging.MustGetLogger("BEANSTALKG")

func main() {
	cfg := getConfig();
	initLogging(cfg.Debug)
	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	service := ":" + cfg.Port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	stop := make(chan bool)

	tubeRegister := make(chan architecture.Command)
	// use this tube to send the channels for each individual tube to the clients when the do 'use' command
	useTubeConnectionReceiver := make(chan chan architecture.Command)
	watchedTubeConnectionsReceiver := make(chan chan architecture.Command)
	operation.NewTubeRegister(tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop, backend.QueueCreator(cfg.Backend))
	log.Info("BEANSTALKG listening on: ", cfg.Port)

	for {
		// log.Println("BEANSTALKG Waiting..")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		operation.NewClientHandler(conn, tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error:", err.Error())
	}
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

// temporarily put this here so we can start with the new folder structure
func getConfig() *ServerConfig {
	cfg := &ServerConfig{}
	configor.New(&configor.Config{ENVPrefix: "BEANSTALKG"}).Load(cfg, "config.yml")

	flag.StringVar(&cfg.Port, "port", cfg.Port, "Port for beanstalkg server")
	flag.StringVar(&cfg.Backend, "backend", cfg.Backend, "Use this backend for job storage.")
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Start server in debug mode. Logs shall contain more information")
	flag.Parse()

	return cfg
}