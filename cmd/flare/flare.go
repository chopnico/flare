package main

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chopnico/flare"
	fc "github.com/chopnico/flare/internal/cli"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	AppName  string = "flare"
	AppUsage string = "a cloudflare cli/tui tool"
	// ldflags will be used to set this. check Makefile
	AppVersion string

	DefaultLoggingLevel = "info"
	DefaultPrintFormat  = "table"
	DefaultTimeOut      = 60
)

// sets http client options such as ignoring ssl, timeouts, and proxy
func httpOptions(c *cli.Context) (*http.Client, error) {
	var client http.Client

	tr := &http.Transport{}

	if c.Bool("ignore-ssl") {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if c.String("proxy") != "" {
		if u, err := url.Parse(c.String("proxy")); err != nil {
			return nil, err
		} else {
			tr.Proxy = http.ProxyURL(u)
		}
	}

	client.Transport = tr

	client.Timeout = time.Duration(c.Int("timeout")) * time.Second

	return &client, nil
}

func main() {
	// create a pretty console logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	api := &flare.Api{
		Logger: &logger,
	}

	app := cli.NewApp()
	app.Name = AppName
	app.Usage = AppUsage
	app.Version = AppVersion
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Usage:       "api `TOKEN`",
			EnvVars:     []string{"CLOUDFLARE_TOKEN"},
			Required:    false,
			DefaultText: "none",
		},
		&cli.StringFlag{
			Name:        "email",
			Usage:       "api `EMAIL`",
			EnvVars:     []string{"CLOUDFLARE_EMAIL"},
			Required:    false,
			DefaultText: "none",
		},
		&cli.StringFlag{
			Name:        "key",
			Usage:       "api `KEY`",
			EnvVars:     []string{"CLOUDFLARE_KEY"},
			Required:    false,
			DefaultText: "none",
		},
		&cli.BoolFlag{
			Name:  "ignore-ssl",
			Usage: "ignore ssl errors",
			Value: false,
		},
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "http `TIMEOUT`",
			Value: DefaultTimeOut,
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "printing `FORMAT` (json, list, table)",
			Value: DefaultPrintFormat,
		},
		&cli.StringFlag{
			Name:  "logging",
			Usage: "set `LOGGING` level",
			Value: DefaultLoggingLevel,
		},
		&cli.StringFlag{
			Name:  "proxy",
			Usage: "set http `PROXY`",
		},
	}

	app.Before = func(c *cli.Context) error {
		http, err := httpOptions(c)
		if err != nil {
			return err
		}

		if c.String("token") != "" {
			api.Client, err = cloudflare.NewWithAPIToken(c.String("token"), cloudflare.HTTPClient(http))
			if err != nil {
				return err
			}
		} else if c.String("email") != "" && c.String("key") != "" {
			api.Client, err = cloudflare.New(c.String("key"), c.String("email"), cloudflare.HTTPClient(http))
			if err != nil {
				return err
			}
		} else {
			return errors.New("you must either supply an api token or an email and an api key")
		}

		return nil
	}

	fc.NewCommands(app, api)

	err := app.Run(os.Args)
	if err != nil {
		api.Logger.Error().Err(err).Msg("failed to run flare")
	}

	os.Exit(0)
}
