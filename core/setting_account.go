package core

import (
	"fmt"
	"strings"
	"telegram-web/helper"
)

func settingAccountTools() {
	fmt.Println("<=====================[Setting Account Tools]=====================>")

	subTools := []string{
		"Get Detail Account",
		"Set Account Username",
		"Set First Name (Upcoming)",
		"Set Last Name (Upcoming)",
		"Set Account Password (Upcoming)",
	}

	for index, tool := range subTools {
		helper.PrettyLog(fmt.Sprintf("%v", index+1), tool)
	}

	selectedSubTools = helper.InputChoice(len(subTools) + 1)

	processSelectedSubTools()
}

func getDetailAccount() {
	fmt.Println("<=====================[Get Detail Account]=====================>")

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = selectOptionsAccount()

	processOptionsAccount(files, true)
}

func setUsername() {
	fmt.Println("<=====================[Set Account Username]=====================>")

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = selectOptionsAccount()

	processOptionsAccount(files, false)
}

func setFirstName() {
	fmt.Println("<=====================[Set Account First Name]=====================>")

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = selectOptionsAccount()

	processOptionsAccount(files, false)
}

func setLastName() {
	fmt.Println("<=====================[Set Account Last Name]=====================>")

	choice := helper.InputTerminal("Do You Want To Change Same Last Name For All Account ? (y/n) (default = n): ")

	if choice == "y" || choice == "Y" {
		lastName = strings.TrimSpace(helper.InputTerminal("Input Batch Last Name: "))
	}

	files := helper.ReadFileDir(localStoragePath)

	helper.PrettyLog("info", fmt.Sprintf("%v Session Local Storage Detected", len(files)))

	selectedOptionsAccount = selectOptionsAccount()

	processOptionsAccount(files, false)
}
