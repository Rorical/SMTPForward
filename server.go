package main

import (
	"flag"
	"github.com/Rorical/SMTPForward/server"
	"github.com/Rorical/SMTPForward/util"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.json", "config file path")
	flag.Parse()

	cfg, err := util.ReadServerConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	err = server.Serve(cfg)
	if err != nil {
		panic(err)
	}
}
