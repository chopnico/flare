package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const (
	defaultTimeout = 10
)

type App struct {
	Token string `yaml:"token"`
	Timeout int `yaml:"timeout"`
	Location string `yaml:"-"`
}

func (app *App) Init(token string) error {
	app.Token = token
	app.Timeout = defaultTimeout

	data, err := yaml.Marshal(&app)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(app.Location, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Read() error {
	file, err := ioutil.ReadFile(app.Location)
	if err != nil {
		return err
	}

	yaml.Unmarshal(file, &app)

	return nil
}
