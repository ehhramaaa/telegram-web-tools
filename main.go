package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"telegram-web/helper"
	"telegram-web/launcher"
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

	tools := []func(){
		func() {
			helper.PrettyLog("0", "Exit Program")
		},
		func() {
			helper.PrettyLog("1", "Get Local Storage")
		},
		func() {
			helper.PrettyLog("2", "Get Query Data Tools")
		},
	}

	isRepeat := true

	for isRepeat {
		helper.ClearTerminal()

		fmt.Println(`✩░▒▓▆▅▃▂▁𝐭𝐞𝐥𝐞𝐠𝐫𝐚𝐦 𝐰𝐞𝐛 𝐭𝐨𝐨𝐥𝐬▁▂▃▅▆▓▒░✩`)
		fmt.Println("ρσɯҽɾҽԃ Ⴆყ : ԋσʅყƈαɳ")

		for _, tool := range tools {
			tool()
		}

		fmt.Print("\n")

		choice, err := strconv.Atoi(strings.TrimSpace(helper.InputTerminal("Masukan Pilihan: ")))
		if err != nil {
			helper.PrettyLog("error", err.Error())
			return
		}

		if choice > (len(tools) - 1) {
			helper.PrettyLog("error", "Pilihan tidak valid")
			time.Sleep(2 * time.Second)
			continue
		}

		switch choice {
		case 0:
			helper.PrettyLog("success", "Exiting Program...")
			os.Exit(0)
		case 1:
			launcher.GetLocalStorage()
		case 2:
			launcher.GetQueryData()
		case 3:
			return
		}

		isInvalid := true
		for isInvalid {
			choice := strings.ToLower(strings.TrimSpace(helper.InputTerminal("Repeat Program ? (y/n): ")))

			switch choice {
			case "n":
				isRepeat = false
				isInvalid = false
			case "y":
				isInvalid = false
			default:
				helper.PrettyLog("error", "Pilihan tidak valid")
			}
		}
	}
}
