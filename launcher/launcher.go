package launcher

import (
	"fmt"
	"strconv"
	"strings"
	"telegram-web/bot"
	"telegram-web/helper"
	"time"

	"github.com/gookit/config/v2"
)

func GetQueryData() {
	tools := []func(){
		func() {
			helper.PrettyLog("0", "Exit Tools")
		},
		func() {
			helper.PrettyLog("1", "Get Query All Account")
		},
		func() {
			helper.PrettyLog("2", "Get Query One Account")
		},
		func() {
			helper.PrettyLog("3", "Merge All Data")
		},
	}

	isRepeat := true
	for isRepeat {
		helper.ClearTerminal()

		fmt.Println("<=====================[Get Query Data Tools]=====================>")

		var isAutoRef bool

		botUsername := config.String("bot.username")
		refUrl := config.String("bot.ref-url")
		localStoragePath := config.String("folder-name.local-storage")
		queryDataPath := config.String("folder-name.query-data")
		mergePath := config.String("folder-name.merge-query-data")
		mergeFileName := fmt.Sprintf("merge_%v", time.Now().Format("2006-01-02-15-04-05"))

		// Membaca semua file dari folder localStorage
		files := helper.ReadFileDir(localStoragePath)

		helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

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
		} else if choice == 0 {
			break
		}

		isStop := false

		for !isStop {
			choice := strings.ToLower(strings.TrimSpace(helper.InputTerminal(fmt.Sprintf("Is First Time Launch Bot %v ? (y/n): ", botUsername))))

			switch choice {
			case "n":
				isAutoRef = false
				isStop = true
			case "y":
				if len(config.Strings("bot.selector")) >= 1 {
					isAutoRef = true
				} else {
					helper.PrettyLog("error", "Value Of Bot Selector Is Null, Please Check Your Config.yml")
					return
				}

				isStop = true
			default:
				helper.PrettyLog("error", "Pilihan tidak valid")
			}
		}

		switch choice {
		case 1:
			bot.GetAllAccount(botUsername, refUrl, isAutoRef, localStoragePath, queryDataPath, files)
		case 2:
			session := strings.TrimSpace(helper.InputTerminal("Masukan Nama File Session : "))

			if !strings.Contains(session, ".json") {
				session = session + ".json"
			}

			if !strings.Contains(session, "+62") {
				session = "+62" + session
			}

			if helper.CheckFileOrFolder(fmt.Sprintf("%v/%v", localStoragePath, session)) {
				bot.GetOneAccount(session, botUsername, refUrl, isAutoRef, localStoragePath, queryDataPath)
			} else {
				helper.PrettyLog("error", fmt.Sprintf("File Sessions %v Not Found", session))
			}
		case 3:
			bot.MergeData(botUsername, queryDataPath, mergePath, mergeFileName, files)
		}

		isInvalid := true
		for isInvalid {
			choice := strings.ToLower(strings.TrimSpace(helper.InputTerminal("Repeat Query Data Tools ? (y/n): ")))

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

func GetLocalStorage() {
	isRepeat := true
	for isRepeat {
		helper.ClearTerminal()

		fmt.Println("<=====================[Get Local Storage Session]=====================>")
		url := config.String("telegram.url")
		sessionsPath := config.String("folder-name.local-storage")
		passwordAccount := config.String("telegram.login.password")

		bot.SessionLocalStorage(passwordAccount, url, sessionsPath)

		isInvalid := true
		for isInvalid {
			choice := strings.ToLower(strings.TrimSpace(helper.InputTerminal("Repeat Tools ? (y/n): ")))

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
