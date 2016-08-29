package main

import (
	"flag"
	"youtube/server"
	"youtube/utils/loadConfig"
	"youtube/utils/logs"
)

var (
	flConfigFile = flag.String("config-file", "config-default.json", "Load config from file")
	l            = logs.New("API_server")
)

func main() {
	flag.Parse()

	var cfg server.Config
	err := loadConfig.FromFileAndEnv(&cfg, *flConfigFile)
	if err != nil {
		l.Fatalln("Error loading config:", err)
	}
	server.Start(cfg)

}
