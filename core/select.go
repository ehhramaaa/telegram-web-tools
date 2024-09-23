package core

import (
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"telegram-web/helper"
)

func selectOptionsAccount() int {
	options := []string{
		"Select All Account",
		"Select Custom Account",
	}

	for index, option := range options {
		helper.PrettyLog(fmt.Sprintf("%v", index+1), option)
	}

	choice := helper.InputChoice(len(options) + 1)

	return choice
}

func selectAccount(files []fs.DirEntry) []int {
	var selectedAccounts []int
	filesPerBatch := 10
	totalFiles := len(files)

	for index := 0; index < totalFiles; index += filesPerBatch {
		var input string

		// Display files in batches of 10
		for i := index; i < index+filesPerBatch && i < totalFiles; i++ {
			helper.PrettyLog(fmt.Sprintf("%v", i+1), fmt.Sprintf("%v", files[i].Name()))
		}

		helper.PrettyLog("input", "Select Account(s) (e.g., '1,3,5-7' or press 'N' to see more): ")

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

		// Split input by comma or space to support multiple choices or ranges
		choices := strings.FieldsFunc(input, func(r rune) bool {
			return r == ',' || r == ' '
		})

		// Iterate through the selected choices
		for _, choice := range choices {
			if strings.Contains(choice, "-") {
				// If input is a range (e.g., 1-5)
				rangeBounds := strings.Split(choice, "-")
				if len(rangeBounds) != 2 {
					helper.PrettyLog("error", fmt.Sprintf("Invalid range: %v. Please try again.", choice))
					return nil
				}

				// Convert range bounds to integers
				start, err1 := strconv.Atoi(rangeBounds[0])
				end, err2 := strconv.Atoi(rangeBounds[1])
				if err1 != nil || err2 != nil || start <= 0 || end >= totalFiles+1 || start > end {
					helper.PrettyLog("error", fmt.Sprintf("Invalid range: %v. Please try again.", choice))
					return nil
				}

				// Append all numbers within the range
				for i := start; i <= end; i++ {
					selectedAccounts = append(selectedAccounts, i)
				}
			} else {
				// Convert single selection to integer
				selectAccount, err := strconv.Atoi(choice)
				if err != nil || selectAccount <= 0 || selectAccount >= totalFiles+1 {
					helper.PrettyLog("error", fmt.Sprintf("Invalid selection: %v. Please try again.", choice))
					return nil
				}

				// Append valid selections to the slice
				selectedAccounts = append(selectedAccounts, selectAccount)
			}
		}

		// Exit the loop after valid selections
		break
	}

	return selectedAccounts
}

func (c *Client) selectProcess(file fs.DirEntry) {
	switch selectedMainTools {
	case 2:
		c.processJoinSkibidiSigmaCode(file)
	case 3:
		switch selectedSubTools {
		case 1:
			c.processGetDetailAccount(file)
		case 2:
			c.processSetAccountUsername(file)
		}
	case 4:
		c.processStartBotWithAutoRef(file)
	case 5:
		switch selectedSubTools {
		case 1:
			c.processGetQueryData(file)
		case 2:
			mergeQueryData()
		}
	case 6:
		c.processFreeRoam(file)
	}
}

func selectCountry() string {
	var selectedCountry string

	// Read the country list from a file
	listCountry, err := helper.ReadFileTxt("./config/countryList.txt")
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Error reading file: %v", err))
		return ""
	}

	filesPerBatch := 10
	totalCountries := len(listCountry)

	for index := 0; index < totalCountries; index += filesPerBatch {
		var input string

		// Display countries in batches of 10
		for i := index; i < index+filesPerBatch && i < totalCountries; i++ {
			helper.PrettyLog(fmt.Sprintf("%v", i+1), listCountry[i])
		}

		// Prompt the user for input
		helper.PrettyLog("input", "Select Country(s) (e.g., '1', input country name, or press 'N' to see more): ")

		// Read user input
		fmt.Scan(&input)

		// If input is 'N' or 'n', show next batch
		if input == "n" || input == "N" {
			if index+filesPerBatch >= totalCountries {
				helper.PrettyLog("info", "No more countries to display.")
				return ""
			} else {
				// Continue to next batch
				continue
			}
		}

		// Try to parse the input as an integer (country index)
		if countryIndex, err := strconv.Atoi(input); err == nil {
			if countryIndex > 0 && countryIndex <= totalCountries {
				selectedCountry = listCountry[countryIndex-1]
				break
			} else {
				helper.PrettyLog("error", "Invalid selection. Please try again.")
				continue
			}
		}

		// If input is not a number, treat it as a country name search
		input = strings.TrimSpace(strings.ToLower(input)) // Normalize input to lower case for matching
		for _, country := range listCountry {
			if strings.ToLower(country) == input || strings.Contains(strings.ToLower(country), input) {
				helper.PrettyLog("info", fmt.Sprintf("Country Found: %v", country))

				confirm := helper.InputTerminal("Confirm Country (y/n) (n = Next) : ")

				if confirm == "y" || confirm == "Y" {
					selectedCountry = country
					break
				} else {
					continue
				}
			}
		}

		// If country is found, break out of the loop
		if selectedCountry != "" {
			break
		} else {
			helper.PrettyLog("error", "Country not found. Please try again.")
		}
	}

	return selectedCountry
}

func selectBot(botList interface{}) (string, int) {
	helper.PrettyLog("info", fmt.Sprintf("%v Bot Detected", len(botList.([]interface{}))))

	for index, bot := range botList.([]interface{}) {
		// Convert each bot to a map to access keys and values
		botMap := bot.(map[string]interface{})

		// Loop through the map to get the key names (like "MIDAS")
		for botName := range botMap {
			helper.PrettyLog(fmt.Sprintf("%v", index+1), botName)
		}
	}

	choice := helper.InputChoice(len(botList.([]interface{})) + 1)

	selectedBot := botList.([]interface{})[choice-1].(map[string]interface{})

	// Print the selected bot name
	for botName := range selectedBot {
		return botName, choice - 1
	}

	return "", 0
}
