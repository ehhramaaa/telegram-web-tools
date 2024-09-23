package main

import (
	"telegram-web/core"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	// Load Config
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config/config.yml")
	if err != nil {
		panic(err)
	}

	err = config.LoadFiles("config/start_bot_with_auto_ref.yml")
	if err != nil {
		panic(err)
	}

	core.LaunchBot()
}
