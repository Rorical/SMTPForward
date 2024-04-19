package main

import (
	"flag"
	"github.com/Rorical/SMTPForward/client"
	"github.com/Rorical/SMTPForward/util"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.json", "config file path")
	flag.Parse()

	cfg, err := util.ReadClientConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	err = client.Listen(cfg)
	if err != nil {
		panic(err)
	}
}
