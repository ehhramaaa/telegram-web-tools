package core

import (
	"fmt"
	"os"
	"path"
	"strings"
	"telegram-web/helper"
	"time"

	"github.com/gookit/config/v2"
)

func queryDataTools() {
	fmt.Println("<=====================[Query Data Tools]=====================>")

	subTools := []string{
		"Get Query Data",
		"Merge All Query Data",
	}

	for index, tool := range subTools {
		helper.PrettyLog(fmt.Sprintf("%v", index+1), tool)
	}

	selectedSubTools = helper.InputChoice(len(subTools) + 1)

	processSelectedSubTools()
}

func getQueryData() {
	fmt.Println("<=====================[Get Query Data]=====================>")

	// Membaca semua file dari folder localStorage
	files := helper.ReadFileDir(localStoragePath)

	botList := config.Get("BOT_LIST")

	selectedBotList, indexBotList = selectBot(botList)

	botUsername := config.String(fmt.Sprintf("BOT_LIST.%v.%s.BOT_USERNAME", indexBotList, selectedBotList))

	if !helper.CheckFileOrFolder(path.Join(queryDataPath, botUsername)) {
		os.MkdirAll(path.Join(queryDataPath, botUsername), 0755)
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = selectOptionsAccount()

	processOptionsAccount(files, true)
}

func mergeQueryData() {
	fmt.Println("<=====================[Merge Query Data]=====================>")

	accounts, err := os.ReadDir(localStoragePath)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read accounts directory %s: %v", queryDataPath, err))
		return
	}

	folders, err := os.ReadDir(queryDataPath)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read folders directory %s: %v", queryDataPath, err))
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Query Data Folder Detected", len(folders)))

	for index, folder := range folders {
		helper.PrettyLog(fmt.Sprintf("%v", index+1), folder.Name())
	}

	choice := helper.InputChoice(len(folders) + 1)

	files := helper.ReadFileDir(fmt.Sprintf("%s/%s", queryDataPath, folders[choice-1].Name()))
	if files == nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read directory %s: %v", fmt.Sprintf("%s/%s", queryDataPath, folders[choice-1]), err))
		return
	}

	if len(files) <= 1 {
		helper.PrettyLog("error", fmt.Sprintf("No query data or only one file in %s", folders[choice-1].Name()))
		return
	}

	if len(files) < len(accounts) {
		helper.PrettyLog("warning", fmt.Sprintf("%v Account(s) do not have query data in %s", len(accounts)-len(files), folders[choice-1].Name()))
		helper.PrettyLog("info", fmt.Sprintf("List of phone numbers that do not have query data %s:", folders[choice-1].Name()))

		for index, account := range accounts {
			accountName := strings.TrimSuffix(account.Name(), ".json")
			hasQueryData := false

			// Check if any file contains the account name (assumed that file names include account name)
			for _, file := range files {
				if strings.Contains(file.Name(), accountName) {
					hasQueryData = true
					break
				}
			}

			// If the account does not have any corresponding query data file, log it
			if !hasQueryData {
				helper.PrettyLog(fmt.Sprintf("%v", index+1), accountName)
			}
		}
	}

	var mergedData []string

	for _, file := range files {
		// Baca data dari file
		if !file.IsDir() {
			account, err := helper.ReadFileTxt(fmt.Sprintf("%s/%s/%s", queryDataPath, folders[choice-1].Name(), file.Name()))
			if err != nil {
				helper.PrettyLog("error", fmt.Sprintf("Failed to read file %s: %v", file.Name(), err))
				continue
			}

			for _, value := range account {
				mergedData = append(mergedData, value)
			}
		}
	}

	// Check Folder
	mergePath := fmt.Sprintf("%s/%s/merge/", queryDataPath, folders[choice-1].Name())
	fileName := fmt.Sprintf("merge_query_data_%s.txt", time.Now().Format("20060102150405"))

	if !helper.CheckFileOrFolder(mergePath) {
		os.Mkdir(fmt.Sprintf(mergePath), 0755)
	}

	if len(mergedData) > 0 {
		// Save Query Data To Txt
		for _, value := range mergedData {
			err := helper.SaveFileTxt(mergePath+fileName, value)
			if err != nil {
				helper.PrettyLog("error", fmt.Sprintf("Error saving file: %v", err))
			}
		}

		helper.PrettyLog("success", "Merge All Query Data Successfully")
	}
}
