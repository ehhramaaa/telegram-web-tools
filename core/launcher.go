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
var localStoragePath, queryDataPath, detailAccountPath string
var maxThread int

func initConfig() {
	outputPath := "./output"
	localStoragePath = "./output/local-storage"
	queryDataPath = "./output/query-data"
	detailAccountPath = "./output/detail-account"
	maxThread = config.Int("MAX_THREAD")

	headlessMode := config.Bool("HEADLESS")

	if headlessMode {
		isHeadless = true
	}

	if !helper.CheckFileOrFolder(outputPath) {
		os.Mkdir(outputPath, 0755)
	}

	if !helper.CheckFileOrFolder(queryDataPath) {
		os.Mkdir(queryDataPath, 0755)
	}

	if !helper.CheckFileOrFolder(detailAccountPath) {
		os.Mkdir(detailAccountPath, 0755)
	}
}

func LaunchBot() {
	initConfig()

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
		helper.PrettyLog("1", "Get Local Storage Session")
		helper.PrettyLog("2", "Join Skibidi Sigma Code Community")
		helper.PrettyLog("3", "Free Roam")
		helper.PrettyLog("4", "Get Detail Account")
		helper.PrettyLog("5", "Set Account Username")
		helper.PrettyLog("6", "Start Bot With Auto Ref")
		helper.PrettyLog("7", "Get Query Data Tools")
		helper.PrettyLog("8", "Merge All Query Data")
		helper.PrettyLog("9", "Join Telegram Group")
		// helper.PrettyLog("8", "Set First Name (Upcoming)")
		// helper.PrettyLog("9", "Set Last Name (Upcoming)")
		// helper.PrettyLog("10", "Set Account Password (Upcoming)")

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
			getLocalStorageSession()
			return
		case 2:
			joinSkibidiSigmaCode()
			return
		case 3:
			freeRoam()
			return
		case 4:
			getDetailAccount()
			return
		case 5:
			setUsername()
			return
		case 6:
			startBotWithAutoRef()
			return
		case 7:
			getQueryData()
			return
		case 8:
			mergeQueryData()
			return
		case 9:
			helper.PrettyLog("info", "Feature Is Upcoming...")
		}
	}
}
