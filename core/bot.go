package core

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"sync"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gookit/config/v2"
)

func selectAccount(files []fs.DirEntry) int {
	var selectedAccount int
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
		fmt.Scan(&input)

		// Check if input is "N" (case-insensitive) to continue showing more files
		if input == "n" || input == "N" {
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

		selectedAccount = selectAccount

		break
	}

	return selectedAccount
}

func selectChoice(files []fs.DirEntry, isMultiThread bool) {
	choice := helper.InputChoice()

	var wg sync.WaitGroup
	var semaphore chan struct{}

	switch choice {
	case 1:
		if isMultiThread {
			if maxThread > len(files) {
				maxThread = len(files)
			}

			semaphore = make(chan struct{}, maxThread)
			for _, file := range files {
				wg.Add(1)
				go processAccountMultiThread(semaphore, &wg, file)
			}
			wg.Wait()
		} else {
			for _, file := range files {
				processAccountSingleThread(file)
			}
		}
	case 2:
		selectedAccount := selectAccount(files)

		if selectedAccount == 0 {
			return
		}

		processAccountSingleThread(files[selectedAccount-1])
		break
	}
}

func getLocalStorage() {
	for {
		helper.ClearTerminal()
		fmt.Println("<=====================[Get Local Storage Session]=====================>")
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

func getDetailAccount() {
	fmt.Println("<=====================[Get Detail Account]=====================>")

	files := helper.ReadFileDir(localStoragePath)

	helper.ClearInputTerminal()

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Get Detail All Account")
	helper.PrettyLog("2", "Get Detail One Account")

	selectChoice(files, true)
}

func setUsername() {
	fmt.Println("<=====================[Set Account Username]=====================>")

	files := helper.ReadFileDir(localStoragePath)

	helper.ClearInputTerminal()

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Set Username All Account")
	helper.PrettyLog("2", "Set Username One Account")

	selectChoice(files, false)
}

func startBotWithAutoRef() {
	fmt.Println("<=====================[Start Bot With Auto Ref]=====================>")

	files := helper.ReadFileDir(localStoragePath)

	helper.ClearInputTerminal()

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Start Bot With Auto Ref All Account")
	helper.PrettyLog("2", "Start Bot With Auto Ref One Account")

	selectChoice(files, true)
}

func getQueryData() {
	fmt.Println("<=====================[Get Query Data Tools]=====================>")

	// Membaca semua file dari folder localStorage
	files := helper.ReadFileDir(localStoragePath)

	helper.ClearInputTerminal()

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))
	helper.PrettyLog("1", "Get Query All Account")
	helper.PrettyLog("2", "Get Query One Account")

	selectChoice(files, true)
}

func mergeQueryData() {
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
