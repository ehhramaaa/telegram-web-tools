package core

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gookit/config/v2"
)

var phoneNumber string

func getLocalStorage() {
	isRepeat := true

	for isRepeat {
		helper.ClearTerminal()
		fmt.Println("<=====================[Get Local Storage Session]=====================>")
		loginUrl := config.String("GET_LOCAL_STORAGE.LOGIN_URL")
		localStoragePath := "./output/local-storage"
		countryAccount := config.String("GET_LOCAL_STORAGE.COUNTRY")
		passwordAccount := config.String("GET_LOCAL_STORAGE.PASSWORD")

		launchOptions := launcher.New().
			Headless(false).
			MustLaunch()

		browser := rod.New().ControlURL(launchOptions).MustConnect()

		defer browser.MustClose()

		processGetLocalStorage(browser, passwordAccount, loginUrl, localStoragePath, countryAccount)

		browser.MustClose()

		var choice string

		helper.PrettyLog("input", "Repeat Program ? (y/n): ")

		_, err := fmt.Scan(&choice)
		if err != nil || choice != "y" || choice != "n" {
			helper.PrettyLog("error", "Invalid selection")
			continue
		}

		if choice == "n" {
			isRepeat = false
		}
	}
}

func startBotWithAutoRef() {
	fmt.Println("<=====================[Start Bot With Auto Ref]=====================>")
	maxThread := config.Int("MAX_THREAD")
	botUsername := config.String("BOT_USERNAME")
	refUrl := config.String("START_BOT_WITH_AUTO_REF.REF_URL")
	localStoragePath := "./output/local-storage"

	files := helper.ReadFileDir(localStoragePath)

	var choice int
	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Start Bot With Auto Ref All Account")
	helper.PrettyLog("2", "Start Bot With Auto Ref One Account")

	helper.PrettyLog("input", "Select Choice: ")

	_, err := fmt.Scan(&choice)
	if err != nil || choice < 0 || choice > 2 {
		helper.PrettyLog("error", "Invalid selection")
		return
	}

	if maxThread > len(files) {
		maxThread = len(files)
	}

	var wg sync.WaitGroup
	var semaphore chan struct{}

	processAccount := func(file fs.DirEntry) {
		defer wg.Done()
		semaphore <- struct{}{}

		phoneNumber = strings.TrimSuffix(file.Name(), ".json")

		helper.PrettyLog("info", fmt.Sprintf("| %s | Start Processing Account...", phoneNumber))

		extensionPath, _ := filepath.Abs("./extension/mini-app-android-spoof")

		launchOptions := launcher.New().
			Set("load-extension", extensionPath).
			Headless(true).
			MustLaunch()

		browser := rod.New().ControlURL(launchOptions).MustConnect()

		defer browser.MustClose()

		processStartBotWithAutoRef(browser, file, localStoragePath, botUsername, refUrl)

		helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished...", phoneNumber))

		<-semaphore
	}

	switch choice {
	case 1:
		semaphore = make(chan struct{}, maxThread)
		for _, file := range files {
			wg.Add(1)
			go processAccount(file)
		}
		wg.Wait()
	case 2:
		maxThread = 1
		semaphore = make(chan struct{}, maxThread)
		filesPerBatch := 10
		totalFiles := len(files)

		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		for index := 0; index < totalFiles; index += filesPerBatch {
			var input string

			// Display files in batches of 10
			for i := index; i < index+filesPerBatch && i < totalFiles; i++ {
				helper.PrettyLog(fmt.Sprintf("%v", i+1), fmt.Sprintf("%v", files[i].Name()))
			}

			helper.PrettyLog("input", "Select Account (or press 'N' to see more): ")

			// Read user input
			_, err := fmt.Scan(&input)

			// Check if input is "N" (case-insensitive) to continue showing more files
			if err != nil || input == "n" || input == "N" {
				// If we're at the last batch, notify the user and stop asking
				if index+filesPerBatch >= totalFiles {
					helper.PrettyLog("info", "No more files to display.")
					break
				}
				// Continue to next batch
				continue
			}

			// Convert input to integer for file selection
			selectAccount, err := strconv.Atoi(input)
			if err != nil || selectAccount <= 0 || selectAccount >= totalFiles+1 {
				helper.PrettyLog("error", "Invalid selection. Please try again.")
			} else {
				wg.Add(1)
				go processAccount(files[selectAccount-1])
				wg.Wait()
				break
			}
		}
	}
}

func getQueryData() {
	fmt.Println("<=====================[Get Query Data Tools]=====================>")

	var choice int

	maxThread := config.Int("MAX_THREAD")
	botUsername := config.String("BOT_USERNAME")
	localStoragePath := "./output/local-storage"
	queryDataPath := "./output/query-data"

	if !helper.CheckFileOrFolder(queryDataPath) {
		os.Mkdir(queryDataPath, 0755)
	}

	// Membaca semua file dari folder localStorage
	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Get Query All Account")
	helper.PrettyLog("2", "Get Query One Account")

	helper.PrettyLog("input", "Select Choice: ")

	_, err := fmt.Scan(&choice)
	if err != nil || choice < 0 || choice > 2 {
		helper.PrettyLog("error", "Invalid selection")
		return
	}

	if maxThread > len(files) {
		maxThread = len(files)
	}

	var wg sync.WaitGroup
	var semaphore chan struct{}

	processAccount := func(file fs.DirEntry) {
		defer wg.Done()
		semaphore <- struct{}{}

		phoneNumber = strings.TrimSuffix(file.Name(), ".json")

		helper.PrettyLog("info", fmt.Sprintf("| %s | Start Processing Account...", phoneNumber))

		extensionPath, _ := filepath.Abs("./extension/mini-app-android-spoof")

		launchOptions := launcher.New().
			Set("load-extension", extensionPath).
			Headless(true).
			MustLaunch()

		browser := rod.New().ControlURL(launchOptions).MustConnect()

		defer browser.MustClose()

		processGetQueryData(browser, file, localStoragePath, botUsername)

		if !helper.CheckFileOrFolder(fmt.Sprintf("%s/%s", queryDataPath, botUsername)) {
			os.Mkdir(fmt.Sprintf("%s/%s", queryDataPath, botUsername), 0755)
		}

		filePath := fmt.Sprintf("%s/%s/%s", queryDataPath, botUsername, fmt.Sprintf("query_data_all_account_%s.txt", time.Now().Format("20060102150405")))
		filePathOneAccount := fmt.Sprintf("%s/%s/%s", queryDataPath, botUsername, fmt.Sprintf("query_data_%s.txt", phoneNumber))

		// Save Query Data
		for _, queryData := range sessionStorage {
			if maxThread == 1 {
				helper.SaveFileTxt(filePathOneAccount, queryData)
				break
			}
			helper.SaveFileTxt(filePath, queryData)
		}

		if helper.CheckFileOrFolder(filePath) || helper.CheckFileOrFolder(filePathOneAccount) {
			helper.PrettyLog("success", fmt.Sprintf("| %s | Query Data %s Successfully Saved", phoneNumber, botUsername))
		} else {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Query Data %s Failed Saved", phoneNumber, botUsername))
		}

		helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished...", phoneNumber))

		<-semaphore
	}

	switch choice {
	case 1:
		semaphore = make(chan struct{}, maxThread)
		for _, file := range files {
			wg.Add(1)
			go processAccount(file)
		}
		wg.Wait()
	case 2:
		maxThread = 1
		semaphore = make(chan struct{}, maxThread)
		filesPerBatch := 10
		totalFiles := len(files)

		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		for index := 0; index < totalFiles; index += filesPerBatch {
			var input string

			// Display files in batches of 10
			for i := index; i < index+filesPerBatch && i < totalFiles; i++ {
				helper.PrettyLog(fmt.Sprintf("%v", i+1), fmt.Sprintf("%v", files[i].Name()))
			}

			helper.PrettyLog("input", "Select Account (or press 'N' to see more): ")

			// Read user input
			_, err := fmt.Scan(&input)

			// Check if input is "N" (case-insensitive) to continue showing more files
			if err != nil || input == "n" || input == "N" {
				// If we're at the last batch, notify the user and stop asking
				if index+filesPerBatch >= totalFiles {
					helper.PrettyLog("info", "No more files to display.")
					break
				}
				// Continue to next batch
				continue
			}

			// Convert input to integer for file selection
			selectAccount, err := strconv.Atoi(input)
			if err != nil || selectAccount <= 0 || selectAccount >= totalFiles+1 {
				helper.PrettyLog("error", "Invalid selection. Please try again.")
			} else {
				wg.Add(1)
				go processAccount(files[selectAccount-1])
				wg.Wait()
				break
			}
		}
	}
}

func ProcessChoice(choice int) {
	if choice == 0 {
		helper.PrettyLog("success", "Exiting Program...")
		os.Exit(0)
	} else {
		switch choice {
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
