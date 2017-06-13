package main

import (
	"flag"

	"github.com/jinzhu/configor"
)

type ServerConfig struct {
	Port    string `default:"11300"`
	Debug   bool
	Backend string `default:"minheap"`
}

// getConfig sets values based on the following order of precedence:
// flags, environment variables, configuration files, and finally
// defaults.
func getConfig() *ServerConfig {
	cfg := &ServerConfig{}
	configor.New(&configor.Config{ENVPrefix: "BEANSTALKG"}).Load(cfg, "config.yml")

	flag.StringVar(&cfg.Port, "port", cfg.Port, "Port for beanstalkg server")
	flag.StringVar(&cfg.Backend, "backend", cfg.Backend, "Use this backend for job storage.")
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Start server in debug mode. Logs shall contain more information")
	flag.Parse()

	return cfg
}
