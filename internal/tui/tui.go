package tui

import (
	"fmt"

	"github.com/chopnico/flare/internal/config"

	"github.com/rs/zerolog"
)
type Tui struct {
	Config *config.App
	Logger *zerolog.Logger
}

func (t *Tui) Start() error {
	fmt.Println("Hello")
	return nil
}
