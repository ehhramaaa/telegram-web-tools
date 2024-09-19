package core

import (
	"fmt"
	"os"
	"telegram-web/helper"
	"time"

	"github.com/gookit/config/v2"
)

var selectedTools int
var isHeadless bool

func LaunchBot() {
	if !helper.CheckFileOrFolder("./output") {
		os.Mkdir("./output", 0755)
	}

	headlessMode := config.Bool("HEADLESS")

	if headlessMode {
		isHeadless = true
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

		helper.PrettyLog("0", "Exit Program")
		helper.PrettyLog("1", "Get Local Storage")
		helper.PrettyLog("2", "Start Bot With Auto Ref")
		helper.PrettyLog("3", "Get Query Data Tools")
		helper.PrettyLog("4", "Set Username (Upcoming)")
		helper.PrettyLog("5", "Set First Name (Upcoming)")
		helper.PrettyLog("6", "Set Last Name (Upcoming)")

		helper.PrettyLog("input", "Select Your Choice: ")

		_, err := fmt.Scan(&selectedTools)
		if err != nil || selectedTools < 0 || selectedTools > 5 {
			helper.PrettyLog("error", "Invalid selection")
			time.Sleep(2 * time.Second)
			continue
		}

		processChoice(selectedTools)

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

func processChoice(selectedTools int) {
	if selectedTools == 0 {
		helper.PrettyLog("success", "Exiting Program...")
		os.Exit(0)
	} else {
		switch selectedTools {
		case 1:
			getLocalStorage()
			return
		case 2:
			startBotWithAutoRef()
			return
		case 3:
			getQueryData()
			return
		case 4, 5, 6:
			helper.PrettyLog("info", "Feature Is Upcoming...")
		}
	}
}
