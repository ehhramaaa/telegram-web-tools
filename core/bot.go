package core

import (
	"fmt"
	"strings"
	"telegram-web/helper"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gookit/config/v2"
)

func getLocalStorageSession() {
	for {
		helper.ClearTerminal()
		fmt.Println("<=====================[Get Local Storage Session]=====================>")
		passwordAccount := config.String("ACCOUNT_PASSWORD")

		launchOptions := launcher.New().
			Headless(config.Bool("HEADLESS_MODE")).
			MustLaunch()

		browser := rod.New().ControlURL(launchOptions).MustConnect()

		defer browser.MustClose()

		client := &Client{
			Browser: browser,
		}

		selectedCountry := selectCountry()

		client.processGetLocalStorageSession(passwordAccount, localStoragePath, selectedCountry)

		browser.MustClose()

		choice := strings.ToLower(helper.InputTerminal("Repeat Program ? (y/n): "))

		if choice != "y" || choice != "n" {
			helper.PrettyLog("error", "Invalid selection")
			return
		}

		if choice == "n" {
			return
		}
	}
}

func freeRoam() {
	fmt.Println("<=====================[Free Roam]=====================>")

	config.Set("HEADLESS_MODE", false)

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = 2

	processOptionsAccount(files, false)
}

func startBotWithAutoRef() {
	fmt.Println("<=====================[Start Bot With Auto Ref]=====================>")
	files := helper.ReadFileDir(localStoragePath)

	botList := config.Get("BOT_LIST")

	selectedBotList, indexBotList = selectBot(botList)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = selectOptionsAccount(files)

	processOptionsAccount(files, true)
}

func joinSkibidiSigmaCode() {
	fmt.Println("<=====================[Join Skibidi Sigma Code Community]=====================>")

	channelUsername = "skibidi_sigma_code"

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = 1

	processOptionsAccount(files, true)
}

func autoSubscribeChannel() {
	fmt.Println("<=====================[Auto Subscribe Telegram Channel]=====================>")

	channelUsername = strings.TrimSpace(helper.InputTerminal("Input Channel Name / Username: "))

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = 1

	processOptionsAccount(files, true)
}
