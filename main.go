package main

import (
	"encoding/json"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"github.com/vimukthi-git/beanstalkg/operation"
	"log"
	"net"
	"os"
	"flag"
)

func main() {
	port := flag.String("p", "11300", "Port for beanstalkg server")
	flag.Parse()
	service := ":" + *port
	tubeRegister := make(chan architecture.Command)
	// use this tube to send the channels for each individual tube to the clients when the do 'use' command
	useTubeConnectionReceiver := make(chan chan architecture.Command)
	watchedTubeConnectionsReceiver := make(chan chan architecture.Command)
	stop := make(chan bool)
	operation.NewTubeRegister(tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	log.Println("BEANSTALKG listening on: ", *port)

	for {
		log.Println("BEANSTALKG Waiting..")
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
