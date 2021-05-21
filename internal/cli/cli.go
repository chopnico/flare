package cli

import (
	"github.com/chopnico/flare/internal/config"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type App struct {
	Name	string
	Version	string
	Config 	*config.App
	Logger 	*zerolog.Logger
}

func (a *App) Run(args []string) error {
	app := &cli.App{
		Name: a.Name,
		Version: a.Version,
	}

	NewInitCommand(app, a.Config, a.Logger)
	NewTuiCommand(app, a.Config, a.Logger)
	NewZoneCommand(app, a.Config, a.Logger)
	NewDnsCommand(app, a.Config, a.Logger)

	if err := app.Run(args); err != nil {
		a.Logger.Error().Err(err)
		return err
	}

	return nil
}