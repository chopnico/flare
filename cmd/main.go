package main

import (
	"fmt"
	"os"

	"github.com/chopnico/flare/internal/config"
	"github.com/chopnico/flare/internal/cli"

	"github.com/rs/zerolog"
)

var (
	home, _ 		= os.UserHomeDir()
	root 			= home + "/.flare"
	logLocation		= root + "/app.log"
	configLocation 	= root + "/config.yml"
)

func init() {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		os.Mkdir(root, 0700)
	}

	if f, _ := os.Stat(logLocation); f != nil {
		os.Rename(logLocation, logLocation + ".1")
	}
}

func main() {
	var logger zerolog.Logger

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log, err := os.OpenFile(logLocation, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger = zerolog.New(log).With().Logger()
	} else {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}

	config := config.App { Location: configLocation }

	app := cli.App{
		Config: &config,
		Logger: &logger,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}
}
