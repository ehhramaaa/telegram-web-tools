package main

import (
	"bufio"
	"fmt"
	"os"
	"telegram-web/core"
	"telegram-web/helper"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func handleExit() {
	if r := recover(); r != nil {
		helper.PrettyLog("error", fmt.Sprintf("%v", r))
		helper.PrettyLog("info", "Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func main() {
	defer handleExit()

	// Load Config
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config/config.yml")
	if err != nil {
		panic(err)
	}

	err = config.LoadFiles("config/bot_config.yml")
	if err != nil {
		panic(err)
	}

	core.LaunchBot()
}
