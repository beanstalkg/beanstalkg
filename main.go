package main

import (
	_ "net/http/pprof"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("BEANSTALKG")

func main() {
	if err := serverCmd.Execute(); err != nil {
		log.Fatal(err)
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
