package cli

import (
	"context"
	"errors"

	"github.com/chopnico/flare/internal/config"

	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func NewAccountCommand(app *cli.App, config *config.App, logger *zerolog.Logger) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name: "account",
			Aliases: []string{"a"},
			Usage: "Manage Cloudflare members",
			Before: func(c *cli.Context) error {
				if err := config.Read(); err != nil {
					logger.Error().Err(err)
					return err
				}
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name: "list",
					Usage: "List accounts",
					Aliases: []string{"l"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "properties",
							Aliases: []string{"p"},
							Usage: "Filter `PROPERTIES` of account (comma separated)",
							Required: false,
						},
					},
					Action: func(c *cli.Context) error {
						api, err := cloudflare.NewWithAPIToken(config.Token)
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to generate API client. Please check logs.")
						}

						accounts, _, err := api.Accounts(context.Background(), cloudflare.PaginationOptions{})
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to get list of accounts. Please check logs.")
						}

						var list []interface{}
						for _, account := range accounts {
							list = append(list, account)
						}

						writeOutput(&list, c.String("properties"), c.String("output-format"))

						return nil
					},
				},
				{
					Name: "member",
					Usage: "Manage members of an account.",
					Aliases: []string{"m"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "account-id",
							Usage: "`ACCOUNT-ID` of account",
							Aliases: []string{"ai"},
							Required: true,
						},
						&cli.StringFlag{
							Name: "id",
							Usage: "`ID` of account",
							Aliases: []string{"i"},
							Required: false,
						},
						&cli.StringFlag{
							Name: "properties",
							Usage: "`PROPERTIES` of account",
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

						members, _, err := api.AccountMembers(context.Background(), c.String("account-id"), cloudflare.PaginationOptions{})
						if err != nil {
							logger.Error().Err(err).Msg("")
							return errors.New("Unable to get list of accounts. Please check logs.")
						}

						var list []interface{}
						for _, member := range members {
							list = append(list, member)
						}

						writeOutput(&list, c.String("properties"), c.String("output-format"))

						return nil
					},
				},
			},
		},
	)
}
