package main

import (
	"net"
	"os"
	"fmt"
	"github.com/vimukthi-git/beanstalkg/operation"
	"encoding/json"
	"log"
)

func main() {
	service := ":11300"
	beanstalkg := operation.TubeRegister{}
	beanstalkg.Init()
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
		go beanstalkg.HandleClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
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