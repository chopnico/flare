package cli

import (
	"github.com/chopnico/flare"

	"github.com/urfave/cli/v2"
)

func NewCommands(app *cli.App, api *flare.Api) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name:        "zone",
			Aliases:     []string{"z"},
			Usage:       "interact with zones",
			Subcommands: zoneCommands(app, api),
		},
		&cli.Command{
			Name:        "dns",
			Aliases:     []string{"d"},
			Usage:       "interact with dns",
			Subcommands: dnsCommands(app, api),
		},
	)
}
