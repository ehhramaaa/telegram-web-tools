package main

import (
	"fmt"
	"telegram-web/core"
	"telegram-web/helper"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	// Load Config
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")
	if err != nil {
		panic(err)
	}

	isRepeat := true
	for isRepeat {
		helper.ClearTerminal()

		fmt.Println(`
 /$$$$$$$$        /$$                 /$$      /$$           /$$             /$$$$$$$$                  /$$          
|__  $$__/       | $$                | $$  /$ | $$          | $$            |__  $$__/                 | $$          
   | $$  /$$$$$$ | $$  /$$$$$$       | $$ /$$$| $$  /$$$$$$ | $$$$$$$          | $$  /$$$$$$   /$$$$$$ | $$  /$$$$$$$
   | $$ /$$__  $$| $$ /$$__  $$      | $$/$$ $$ $$ /$$__  $$| $$__  $$         | $$ /$$__  $$ /$$__  $$| $$ /$$_____/
   | $$| $$$$$$$$| $$| $$$$$$$$      | $$$$_  $$$$| $$$$$$$$| $$  \ $$         | $$| $$  \ $$| $$  \ $$| $$|  $$$$$$ 
   | $$| $$_____/| $$| $$_____/      | $$$/ \  $$$| $$_____/| $$  | $$         | $$| $$  | $$| $$  | $$| $$ \____  $$
   | $$|  $$$$$$$| $$|  $$$$$$$      | $$/   \  $$|  $$$$$$$| $$$$$$$/         | $$|  $$$$$$/|  $$$$$$/| $$ /$$$$$$$/
   |__/ \_______/|__/ \_______/      |__/     \__/ \_______/|_______/          |__/ \______/  \______/ |__/|_______/ 
`)

		fmt.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")

		var choice int

		helper.PrettyLog("0", "Exit Program")
		helper.PrettyLog("1", "Get Local Storage")
		helper.PrettyLog("2", "Start Bot With Auto Ref")
		helper.PrettyLog("3", "Get Query Data Tools")
		helper.PrettyLog("4", "Set Username (Upcoming)")
		helper.PrettyLog("5", "Set First Name (Upcoming)")
		helper.PrettyLog("6", "Set Last Name (Upcoming)")

		helper.PrettyLog("input", "Select Your Choice: ")

		_, err := fmt.Scan(&choice)
		if err != nil || choice < 0 || choice > 5 {
			helper.PrettyLog("error", "Invalid selection")
			time.Sleep(2 * time.Second)
			continue
		}

		core.ProcessChoice(choice)

		isInvalid := true
		for isInvalid {
			var choice string
			helper.PrettyLog("input", "Repeat Program ? (y/n): ")

			_, err := fmt.Scan(&choice)
			if err != nil {
				helper.PrettyLog("error", "Invalid selection")
				time.Sleep(2 * time.Second)
				continue
			}

			switch choice {
			case "n":
				isRepeat = false
				isInvalid = false
			case "y":
				isInvalid = false
			default:
				helper.PrettyLog("error", "Invalid selection")
			}
		}
	}
}
