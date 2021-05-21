package cli

import (
	"github.com/chopnico/flare/internal/config"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func NewInitCommand(app *cli.App, config *config.App, logger *zerolog.Logger) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name: "init",
			Aliases: []string{"i"},
			Usage: "Initialize configuration",
			Flags: []cli.Flag {
				&cli.StringFlag{
					Name: "token",
					Usage: "Cloudflare API token",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				err := config.Init(c.String("token"))
				if err != nil {
					logger.Error().Err(err)
					return err
				}
				return nil
			},
		},
	)
}
