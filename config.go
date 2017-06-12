package main

import (
	"flag"

	"github.com/beanstalkg/beanstalkg/architecture"
	"github.com/beanstalkg/beanstalkg/backend"
	"github.com/jinzhu/configor"
)

type ServerConfig struct {
	Port    string `default:"11300"`
	Debug   bool
	Backend string `default:"minheap"`

	queueCreator architecture.PriorityQueueCreator
}

const defaultBackend = "minheap"

var validBackends map[string]architecture.PriorityQueueCreator = map[string]architecture.PriorityQueueCreator{
	"minheap": func() architecture.PriorityQueue { return &backend.MinHeap{} },
}

// getConfig sets values based on the following order of precedence:
// flags, environment variables, configuration files, and finally
// defaults.
func getConfig() *ServerConfig {
	cfg := &ServerConfig{}
	configor.New(&configor.Config{ENVPrefix: "BEANSTALKG"}).Load(cfg, "config.yml")

	flag.StringVar(&cfg.Port, "port", cfg.Port, "Port for beanstalkg server")
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Start server in debug mode. Logs shall contain more information")
	flag.Parse()

	// Fallback to default backend if the provided one is invalid.
	// This guarantees that the server will start.
	if _, ok := validBackends[cfg.Backend]; !ok {
		log.Debugf("%s backend not supported, falling back to %s.", cfg.Backend, defaultBackend)
		cfg.Backend = defaultBackend
	}
	cfg.queueCreator = validBackends[cfg.Backend]

	return cfg
}
