package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/chopnico/flare/internal/config"

	"github.com/chopnico/output"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func NewZoneCommand(app *cli.App, config *config.App, logger *zerolog.Logger) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name: "zone",
			Aliases: []string{"z"},
			Usage: "Interact with zones",
			Before: func(c *cli.Context) error {
				if err := config.Read(); err != nil {
					logger.Error().Err(err).Msg("")
					return err
				}
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name: "list",
					Usage: "List all zones",
					Aliases: []string{"l"},
					Action: func(c *cli.Context) error {
						api, err := cloudflare.NewWithAPIToken(config.Token)
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to generate API client. Please check logs.")
						}

						zones, err := api.ListZones(context.Background())
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to list zones. Please check logs.")
						}

						var list []interface{}
						for _, zone := range zones {
							list = append(list, zone)
						}

						fmt.Println(output.FormatList(list, []string{"ID", "Name"}))

						return nil
					},
				},
				{
					Name: "detail",
					Usage: "Print zone details",
					Aliases: []string{"d"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "id",
							Usage: "`ID` of zone",
							Required: false,
						},
						&cli.StringFlag{
							Name: "properties",
							Usage: "Filter `PROPERTIES` of zone (comma separated)",
							Required: false,
						},
						&cli.StringFlag{
							Name: "name",
							Usage: "`NAME` of zone",
							Required: false,
						},
					},
					Action: func(c *cli.Context) error {
						api, err := cloudflare.NewWithAPIToken(config.Token)
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to generate API client. Please check logs.")
						}

						var zone cloudflare.Zone

						if c.String("id") != "" {
							zone, err = api.ZoneDetails(context.Background(), c.String("id"))
						} else if c.String("name") != "" {
							id, err := api.ZoneIDByName(c.String("name"))
							if err != nil {
								logger.Error().Err(err).Msg("")
								return errors.New("Unable to lookup zone by name. Please check logs.")
							}

							zone, err = api.ZoneDetails(context.Background(), id)
						} else {
							return errors.New("You must supply either the ID of the zone or the Name of the zone.")
						}

						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to get zone details. Please check logs.")
						}

						var list []interface{}
						list = append(list, zone)

						if c.String("properties") == "" {
							fmt.Println(output.FormatList(list, nil))
						} else {
							properties := strings.Split(c.String("properties"), ",")
							fmt.Println(output.FormatList(list, properties))
						}

						return nil
					},
				},
			},
		},
	)
}
