package core

import (
	"fmt"
	"os"
	"telegram-web/helper"
	"time"
)

var selectedMainTools, selectedSubTools, selectedOptionsAccount int
var localStoragePath, localStorageExpiredPath, queryDataPath, detailAccountPath string
var selectedBotList string
var indexBotList int

func initConfig() {
	outputPath := "./output"
	errLogPath := "./error"
	localStoragePath = "./output/local-storage"
	localStorageExpiredPath = "./output/local-storage-expired"
	queryDataPath = "./output/query-data"
	detailAccountPath = "./output/detail-account"

	if !helper.CheckFileOrFolder(outputPath) {
		os.Mkdir(outputPath, 0755)
	}

	if !helper.CheckFileOrFolder(queryDataPath) {
		os.Mkdir(queryDataPath, 0755)
	}

	if !helper.CheckFileOrFolder(detailAccountPath) {
		os.Mkdir(detailAccountPath, 0755)
	}

	if !helper.CheckFileOrFolder(errLogPath) {
		os.Mkdir(errLogPath, 0755)
	}
}

func LaunchBot() {
	initConfig()

	isRepeat := true
	for isRepeat {
		helper.ClearTerminal()

		helper.PrintLogo()

		mainTools := []string{
			"Exits Program",
			"Get Local Storage Session",
			"Join Skibidi Sigma Code Community",
			"Setting Account Tools",
			"Start Bot With Auto Ref",
			"Query Data Tools",
			"Free Roam",
			"Auto Join Telegram Group (Upcoming)",
			"Auto Subscribe Telegram Channel (Upcoming)",
		}

		for index, tool := range mainTools {
			helper.PrettyLog(fmt.Sprintf("%v", index), tool)
		}

		selectedMainTools = helper.InputChoice(len(mainTools))

		processSelectedMainTools()

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
