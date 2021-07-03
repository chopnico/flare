package cli

import (
	"context"
	"errors"

	"github.com/chopnico/flare"
	"github.com/cloudflare/cloudflare-go"

	"github.com/urfave/cli/v2"
)

// get zone details
func getZone(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "get",
		Usage:   "get details about a zone",
		Aliases: []string{"g"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "`ID` of zone",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "`NAME` of zone",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "properties",
				Aliases: []string{"p"},
				Usage:   "choose what `PROPERTIES` to display (list format only)",
			},
		},
		Action: func(c *cli.Context) error {
			var zone cloudflare.Zone
			var err error

			// get zone by id
			if c.String("id") != "" {
				zone, err = api.Client.ZoneDetails(context.Background(), c.String("id"))
				if err != nil {
					if c.String("logging") == "debug" {
						api.Logger.Debug().Err(err).Msg("unable to get zone")
					}
					return err
				}
				// get zone by name
			} else if c.String("name") != "" {
				// get zone id by name first
				id, err := api.Client.ZoneIDByName(c.String("name"))
				if err != nil {
					if c.String("logging") == "debug" {
						api.Logger.Debug().Err(err).Msg("unable to get zone id by name")
					}
					return err
				}

				// get zone details
				zone, err = api.Client.ZoneDetails(context.Background(), id)
				if err != nil {
					if c.String("logging") == "debug" {
						api.Logger.Debug().Err(err).Msg("unable to get zone")
					}
					return err
				}
				// error if name nor id was supplied
			} else if c.String("name") == "" && c.String("id") == "" {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to get zone")
				}
				return errors.New("you must supply either the zone id or the zone name")
			}

			// printing time
			switch c.String("format") {
			case "json":
				flare.PrintJson(zone)
			default:
				var l []interface{}
				l = append(l, zone)
				flare.PrintList(&l, c.String("properties"))
			}
			return nil
		},
	}
}

// list all zones
func listZones(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "list all zones",
		Aliases: []string{"l"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "properties",
				Aliases: []string{"p"},
				Usage:   "choose what `PROPERTIES` to display (list format only)",
			},
		},
		Action: func(c *cli.Context) error {
			// get list of zones
			zones, err := api.Client.ListZones(context.Background())
			if err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to list zones")
				}
				return err
			}

			// printing time
			switch c.String("format") {
			case "json":
				flare.PrintJson(zones)
			case "list":
				var l []interface{}
				for _, i := range zones {
					l = append(l, i)
				}
				flare.PrintList(&l, c.String("properties"))
			default:
				data := [][]string{}
				for _, i := range zones {
					data = append(data,
						[]string{i.ID, i.Name, i.Type, i.Status},
					)
				}

				headers := []string{"ID", "Name", "Type", "Status"}
				flare.PrintTable(data, headers)
			}
			return nil
		},
	}
}

// commbined all zone commands
func zoneCommands(app *cli.App, api *flare.Api) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		listZones(app, api),
		getZone(app, api),
	)

	return commands
}
