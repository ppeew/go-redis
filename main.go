package main

import (
	"fmt"
	"go-redis/config"
	"go-redis/lib/logger"
	"go-redis/resp/handler"
	"go-redis/tcp"
	"os"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0.",
	Port: 6379,
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	_, err := os.Stat(configFile)
	if err != nil {
		config.Properties = defaultProperties
	} else {
		config.SetupConfig(configFile)
	}
	err = tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		},
		handler.MakeHandler())
	if err != nil {
		logger.Error(err)
	}
}
