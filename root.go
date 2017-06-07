package main

import (
	"net"
	"net/http"

	"github.com/beanstalkg/beanstalkg/architecture"
	"github.com/beanstalkg/beanstalkg/operation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "beanstalkg",
	Short: "Beanstalkg is a go implementation of beanstalkd - A fast, general-purpose work queue",
	Run:   startServer,
}

func startServer(cmd *cobra.Command, args []string) {
	port := viper.GetString("port")
	debugMode := viper.GetBool("debug_mode")

	initLogging(debugMode)
	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	stop := make(chan bool)

	tubeRegister := make(chan architecture.Command)
	// use this tube to send the channels for each individual tube to the clients when the do 'use' command
	useTubeConnectionReceiver := make(chan chan architecture.Command)
	watchedTubeConnectionsReceiver := make(chan chan architecture.Command)
	operation.NewTubeRegister(tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop)
	log.Info("BEANSTALKG listening on: ", port)

	for {
		// log.Println("BEANSTALKG Waiting..")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		operation.NewClientHandler(conn, tubeRegister, useTubeConnectionReceiver, watchedTubeConnectionsReceiver, stop)
	}
}

func init() {
	// set persistent flags
	serverCmd.PersistentFlags().String("port", "11300", "Port for beanstalkg server")
	serverCmd.PersistentFlags().Bool("debug_mode", false, "Start server in debug mode. Logs shall contain more information")

	// setup environmental config
	viper.SetEnvPrefix("beanstalkg")
	viper.AutomaticEnv()

	// include support for cli flags
	viper.BindPFlags(serverCmd.PersistentFlags())

	// read config from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}
