package cli

import (
	"context"
	"errors"
	"strings"

	"github.com/chopnico/flare/internal/config"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func NewDnsCommand(app *cli.App, config *config.App, logger *zerolog.Logger) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name: "dns",
			Aliases: []string{"d"},
			Usage: "Interact with DNS records",
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
					Usage: "List DNS records for a zone",
					Aliases: []string{"l"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "id",
							Usage: "`ID` of zone",
							Aliases: []string{"i"},
							Required: false,
						},
						&cli.StringFlag{
							Name: "name",
							Usage: "`NAME` of zone",
							Aliases: []string{"n"},
							Required: false,
						},
						&cli.StringFlag{
							Name: "types",
							Usage: "Filter records by `TYPE` (comma separated)",
							Aliases: []string{"t"},
							Required: false,
						},
						&cli.StringFlag{
							Name: "properties",
							Usage: "Filter `PROPERTIES` of records (comma separated)",
							Aliases: []string{"p"},
							Required: false,
						},
					},
					Action: func(c *cli.Context) error {
						api, err := cloudflare.NewWithAPIToken(config.Token)
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to generate API client. Please check logs.")
						}

						var zoneId string

						if c.String("id") != "" {
							zoneId = c.String("id")
						} else if c.String("name") != "" {
							zoneId, err = api.ZoneIDByName(c.String("name"))
							if err != nil {
								logger.Error().Err(err).Msg("")
								return errors.New("Unable to get zone ID from supplied name. Please check logs.")
							}
						} else {
							return errors.New("You must either supply the name of the zone or the zone ID.")
						}

						records, err := api.DNSRecords(context.Background(), zoneId, cloudflare.DNSRecord{})
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to get DNS records. Please check logs.")
						}

						var list []interface{}
						for _, record := range records {
							if c.String("types") != "" {
								types := strings.Split(c.String("types"), ",")
								for _, t := range types {
									if record.Type == t {
										list = append(list, record)
									}
								}
							} else {
								list = append(list, record)
							}
						}

						writeOutput(&list, c.String("properties"), c.String("output-format"))

						return nil
					},
				},
			},
		},
	)
}
