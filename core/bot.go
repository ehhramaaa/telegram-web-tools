package core

import (
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

func listAccount(files []fs.DirEntry) int {
	var selectAccount int
	filesPerBatch := 10
	totalFiles := len(files)

	helper.ClearInputTerminal()

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
			} else {
				// Continue to next batch
				continue
			}
		}

		// Convert input to integer for file selection
		selectAccount, err := strconv.Atoi(input)
		if err != nil || selectAccount <= 0 || selectAccount >= totalFiles+1 {
			helper.PrettyLog("error", "Invalid selection. Please try again.")
			return 0
		}
	}

	return selectAccount
}

func processAccount(semaphore chan struct{}, wg *sync.WaitGroup, file fs.DirEntry, localStoragePath string, queryDataPath string) {
	botUsername := config.String("BOT_USERNAME")
	refUrl := config.String("START_BOT_WITH_AUTO_REF.REF_URL")

	defer wg.Done()
	semaphore <- struct{}{}

	extensionPath, _ := filepath.Abs("./extension/mini-app-android-spoof")

	launchOptions := launcher.New().
		Set("load-extension", extensionPath).
		Headless(isHeadless).
		MustLaunch()

	browser := rod.New().ControlURL(launchOptions).MustConnect()

	defer browser.MustClose()

	client := &Client{
		phoneNumber: strings.TrimSuffix(file.Name(), ".json"),
		Browser:     browser,
	}

	helper.PrettyLog("info", fmt.Sprintf("| %s | Start Processing Account...", client.phoneNumber))

	switch selectedTools {
	case 2:
		client.processStartBotWithAutoRef(file, localStoragePath, botUsername, refUrl)
	case 3:
		client.processGetQueryData(file, localStoragePath, queryDataPath, botUsername)
	}

	helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished...", client.phoneNumber))

	<-semaphore
}

func getLocalStorage() {
	for {
		helper.ClearTerminal()
		fmt.Println("<=====================[Get Local Storage Session]=====================>")
		localStoragePath := "./output/local-storage"
		countryAccount := config.String("GET_LOCAL_STORAGE.COUNTRY")
		passwordAccount := config.String("GET_LOCAL_STORAGE.PASSWORD")

		launchOptions := launcher.New().
			Headless(isHeadless).
			MustLaunch()

		browser := rod.New().ControlURL(launchOptions).MustConnect()

		defer browser.MustClose()

		client := &Client{
			Browser: browser,
		}

		client.processGetLocalStorage(passwordAccount, localStoragePath, countryAccount)

		browser.MustClose()

		var choice string

		helper.ClearInputTerminal()

		helper.PrettyLog("input", "Repeat Program ? (y/n): ")

		_, err := fmt.Scan(&choice)
		if err != nil || choice != "y" || choice != "n" {
			helper.PrettyLog("error", "Invalid selection")
			continue
		}

		if choice == "n" {
			return
		}
	}
}

func startBotWithAutoRef() {
	fmt.Println("<=====================[Start Bot With Auto Ref]=====================>")
	maxThread := config.Int("MAX_THREAD")
	localStoragePath := "./output/local-storage"

	files := helper.ReadFileDir(localStoragePath)

	var choice int
	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Start Bot With Auto Ref All Account")
	helper.PrettyLog("2", "Start Bot With Auto Ref One Account")

	helper.ClearInputTerminal()

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

	switch choice {
	case 1:
		semaphore = make(chan struct{}, maxThread)
		for _, file := range files {
			wg.Add(1)
			go processAccount(semaphore, &wg, file, localStoragePath, "")
		}
		wg.Wait()
	case 2:
		maxThread = 1
		semaphore = make(chan struct{}, maxThread)

		selectAccount := listAccount(files)
		if selectAccount == 0 {
			return
		}

		wg.Add(1)
		go processAccount(semaphore, &wg, files[selectAccount-1], localStoragePath, "")
		wg.Wait()
		break
	}
}

func getQueryData() {
	fmt.Println("<=====================[Get Query Data Tools]=====================>")

	var choice int

	maxThread := config.Int("MAX_THREAD")
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
	helper.PrettyLog("3", "Merge All Query Data")

	helper.PrettyLog("input", "Select Choice: ")

	_, err := fmt.Scan(&choice)
	if err != nil || choice < 0 || choice > 3 {
		helper.PrettyLog("error", "Invalid selection")
		return
	}

	if maxThread > len(files) {
		maxThread = len(files)
	}

	var wg sync.WaitGroup
	var semaphore chan struct{}

	switch choice {
	case 1:
		semaphore = make(chan struct{}, maxThread)
		for _, file := range files {
			wg.Add(1)
			go processAccount(semaphore, &wg, file, localStoragePath, queryDataPath)
		}
		wg.Wait()
	case 2:
		maxThread = 1
		semaphore = make(chan struct{}, maxThread)

		selectedAccount := listAccount(files)
		if selectedAccount == 0 {
			return
		}

		wg.Add(1)
		go processAccount(semaphore, &wg, files[selectedAccount-1], localStoragePath, queryDataPath)
		wg.Wait()
		break
	case 3:
		mergeQueryData(queryDataPath)
	}
}

func mergeQueryData(queryDataPath string) {
	fmt.Println("<=====================[Merge Query Data]=====================>")

	folders, err := os.ReadDir(queryDataPath)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read directory %s: %v", queryDataPath, err))
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Query Data Folder Detected", len(folders)))

	for index, folder := range folders {
		helper.PrettyLog(fmt.Sprintf("%v", index+1), folder.Name())
	}

	helper.ClearInputTerminal()

	var choice int

	helper.PrettyLog("input", "Select Folder: ")

	fmt.Scan(&choice)
	if choice <= 0 || choice > len(folders)+1 {
		helper.PrettyLog("error", "Invalid selection")
		return
	}

	files := helper.ReadFileDir(fmt.Sprintf("%s/%s", queryDataPath, folders[choice-1].Name()))
	if files == nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read directory %s: %v", fmt.Sprintf("%s/%s", queryDataPath, folders[choice-1]), err))
		return
	}

	var mergedData []string

	for _, file := range files {
		// Baca data dari file
		account, err := helper.ReadFileTxt(fmt.Sprintf("%s/%s/%s", queryDataPath, folders[choice-1].Name(), file.Name()))
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("Failed to read file %s: %v", file.Name(), err))
			continue
		}

		for _, value := range account {
			mergedData = append(mergedData, value)
		}
	}

	// Check Folder
	mergePath := fmt.Sprintf("%s/%s/merge/", queryDataPath, folders[choice-1].Name())
	fileName := fmt.Sprintf("merge_query_data_%s.txt", time.Now().Format("20060102150405"))

	if !helper.CheckFileOrFolder(mergePath) {
		os.Mkdir(fmt.Sprintf(mergePath), 0755)
	}

	// Save Query Data To Txt
	for _, value := range mergedData {
		err := helper.SaveFileTxt(mergePath+fileName, value)
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("Error saving file: %v", err))
		}
	}

	helper.PrettyLog("success", fmt.Sprintf("Merge Query Data Successfully Saved In: %s", mergePath+fileName))
}
