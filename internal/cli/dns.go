package cli

import (
	"context"
	"errors"
	"strings"

	"github.com/chopnico/flare"
	"github.com/cloudflare/cloudflare-go"

	"github.com/imdario/mergo"
	"github.com/urfave/cli/v2"
)

func deleteDnsRecord(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "`DELETE` a dns record on a particuar zone",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "`ID` of record",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "zone-id",
				Usage:    "`ZONE-ID` of record",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-name",
				Usage:    "`ZONE-NAME` of record",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			var zoneId *string

			zoneId, err = getZoneId(c, api)
			if err != nil {
				return err
			}

			err = api.Client.DeleteDNSRecord(context.Background(), *zoneId, c.String("id"))
			if err != nil {
				return err
			}

			api.Logger.Info().Msg("record with id " + c.String("id") + " was deleted")

			return nil
		},
	}
}

func createDnsRecord(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create a dns record on a particuar zone",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "`NAME` of record",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "value",
				Usage:    "record `VALUE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Usage:    "`TYPE` of record",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "proxy",
				Usage:    "should the record be proxied",
				Value:    false,
				Required: false,
			},
			&cli.IntFlag{
				Name:     "ttl",
				Usage:    "`TTL` of record in seconds",
				Value:    300,
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-id",
				Usage:    "`ZONE-ID` of record",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-name",
				Usage:    "`ZONE-NAME` of record",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			var zoneId *string

			zoneId, err = getZoneId(c, api)
			if err != nil {
				return err
			}

			proxy := c.Bool("proxy")
			record := cloudflare.DNSRecord{
				Name:    c.String("name"),
				Content: c.String("value"),
				Type:    strings.ToUpper(c.String("type")),
			}

			if proxy {
				record.TTL = 1
				record.Proxied = &proxy
			}

			responseRecord, err := api.Client.CreateDNSRecord(context.Background(), *zoneId, record)
			if err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to create record with name " + c.String("name"))
				}
				return err
			}

			// printing time
			switch c.String("format") {
			case "json":
				flare.PrintJson(responseRecord.Result)
			default:
				var l []interface{}
				l = append(l, responseRecord.Result)
				flare.PrintList(&l, c.String("properties"))
			}

			return nil
		},
	}
}

func updateDnsRecord(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "update a dns record on a particuar zone",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "`ID` of dns record",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "`NAME` of record",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "value",
				Usage:    "record `VALUE`",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "type",
				Usage:    "`TYPE` of record",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "proxy",
				Usage:    "should the record be proxied",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "ttl",
				Usage:    "`TTL` of record in seconds",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-id",
				Usage:    "`ZONE-ID` of record",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-name",
				Usage:    "`ZONE-NAME` of record",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			var sourceRecord cloudflare.DNSRecord
			var zoneId *string

			zoneId, err = getZoneId(c, api)
			if err != nil {
				return err
			}

			sourceRecord, err = api.Client.DNSRecord(context.Background(), *zoneId, c.String("id"))
			if err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to get record with id " + c.String("id"))
				}
				return err
			}

			destinationRecord := cloudflare.DNSRecord{
				Name:    c.String("name"),
				Content: c.String("value"),
				Type:    c.String("type"),
			}

			// creates the an updated record by merging with the original record
			if err = mergo.Merge(&destinationRecord, sourceRecord); err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to update record with id " + c.String("id"))
				}
				return err
			}

			var proxy bool = c.Bool("proxy")
			if proxy {
				destinationRecord.TTL = 1
				destinationRecord.Proxied = &proxy
			}

			err = api.Client.UpdateDNSRecord(context.Background(), *zoneId, c.String("id"), destinationRecord)
			if err != nil {
				if err != nil {
					if c.String("logging") == "debug" {
						api.Logger.Debug().Err(err).Msg("unable to update record with id " + c.String("id"))
					}
					return err
				}
			}

			record, err := api.Client.DNSRecord(context.Background(), *zoneId, destinationRecord.ID)
			if err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to retrieve updated record with id " + c.String("id"))
				}
				return err
			}

			// printing time
			switch c.String("format") {
			case "json":
				flare.PrintJson(record)
			default:
				var l []interface{}
				l = append(l, record)
				flare.PrintList(&l, c.String("properties"))
			}

			return nil
		},
	}
}

func getDnsRecord(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get a dns record from a particuar zone",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "`ID` of dns record",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "zone-id",
				Usage:    "`ID` of zone",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-name",
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
			var record cloudflare.DNSRecord
			var err error
			var zoneId *string

			zoneId, err = getZoneId(c, api)
			if err != nil {
				return err
			}

			record, err = api.Client.DNSRecord(context.Background(), *zoneId, c.String("id"))
			if err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to get record")
				}
				return err
			}

			// printing time
			switch c.String("format") {
			case "json":
				flare.PrintJson(record)
			default:
				var l []interface{}
				l = append(l, record)
				flare.PrintList(&l, c.String("properties"))
			}
			return nil
		},
	}
}

func listDnsRecords(app *cli.App, api *flare.Api) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list all dns records for a particular zone",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "zone-id",
				Usage:    "`ID` of zone",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "zone-name",
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
			var records []cloudflare.DNSRecord
			var err error
			var zoneId *string

			zoneId, err = getZoneId(c, api)
			if err != nil {
				return err
			}

			records, err = api.Client.DNSRecords(context.Background(), *zoneId, cloudflare.DNSRecord{})
			if err != nil {
				if c.String("logging") == "debug" {
					api.Logger.Debug().Err(err).Msg("unable to get records")
				}
				return err
			}

			// printing time
			switch c.String("format") {
			case "json":
				flare.PrintJson(records)
			case "list":
				var l []interface{}
				for _, i := range records {
					l = append(l, i)
				}
				flare.PrintList(&l, c.String("properties"))
			default:
				data := [][]string{}
				for _, i := range records {
					data = append(data,
						[]string{i.Name, i.Content, i.Type},
					)
				}

				headers := []string{"Name", "Value", "Type"}
				flare.PrintTable(data, headers)
			}
			return nil
		},
	}
}

func getZoneId(c *cli.Context, api *flare.Api) (*string, error) {
	var err error
	var zoneId string

	if c.String("zone-id") != "" {
		zoneId = c.String("zone-id")
		return &zoneId, nil
	} else if c.String("zone-name") != "" {
		zoneId, err = api.Client.ZoneIDByName(c.String("zone-name"))
		if err != nil {
			if c.String("logging") == "debug" {
				api.Logger.Debug().Err(err).Msg("unable to get get zone id with name " + c.String("zone-name"))
			}
			return nil, err
		}
		return &zoneId, nil
	} else if c.String("zone-id") == "" && c.String("zone-name") == "" {
		if c.String("logging") == "debug" {
			api.Logger.Debug().Err(err).Msg("unable to get zone")
		}
		return nil, errors.New("you must supply either the zone id or the zone name")
	}

	return nil, errors.New("unable to retrieve zone id")
}

// commbined all zone commands
func dnsCommands(app *cli.App, api *flare.Api) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		listDnsRecords(app, api),
		getDnsRecord(app, api),
		updateDnsRecord(app, api),
		createDnsRecord(app, api),
		deleteDnsRecord(app, api),
	)

	return commands
}
