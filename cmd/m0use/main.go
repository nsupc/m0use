package main

import (
	"log"
	"log/slog"
	"m0use/pkg/config"
	"os"
)

func main() {
	var path string

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "config.yml"
	}

	conf, err := config.ReadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	config.InitLogger(conf)

	slog.Info("starting")

}
