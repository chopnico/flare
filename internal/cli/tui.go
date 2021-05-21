package cli

import (
	"github.com/chopnico/flare/internal/tui"
	"github.com/chopnico/flare/internal/config"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func NewTuiCommand(app *cli.App, config *config.App, logger *zerolog.Logger) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name: "tui",
			Aliases: []string{"t"},
			Usage: "Run the terminal UI",
			Before: func(c *cli.Context) error {
				if err := config.Read(); err != nil {
					logger.Error().Err(err)
					return err
				}
				return nil
			},
			Action: func(c *cli.Context) error {
				t := tui.Tui{
					Config: config,
					Logger: logger,
				}

				t.Start()

				return nil
			},
		},
	)
}
