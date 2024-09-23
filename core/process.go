package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gookit/config/v2"
)

func processSelectedMainTools() {
	if selectedMainTools != 1 && !helper.CheckFileOrFolder(localStoragePath) {
		helper.PrettyLog("error", "You'r not have a local storage session, Please Get Local Storage Session First...")
		return
	}

	switch selectedMainTools {
	case 0:
		helper.PrettyLog("success", "Exiting Program...")
		os.Exit(0)
	case 1:
		getLocalStorageSession()
		return
	case 2:
		joinSkibidiSigmaCode()
		return
	case 3:
		settingAccountTools()
		return
	case 4:
		startBotWithAutoRef()
		return
	case 5:
		queryDataTools()
		return
	case 6:
		autoSubscribeChannel()
		return
	case 7:
		freeRoam()
		return
	}
}

func processSelectedSubTools() {
	switch selectedMainTools {
	case 3:
		switch selectedSubTools {
		case 1:
			getDetailAccount()
			return
		case 2:
			setUsername()
			return
		case 3:
			setFirstName()
			return
		case 4:
			setLastName()
			return
		}
	case 5:
		switch selectedSubTools {
		case 1:
			getQueryData()
			return
		case 2:
			mergeQueryData()
			return
		}
	}
}

func processOptionsAccount(files []fs.DirEntry, isMultiThread bool) {
	maxThread := config.Int("MAX_THREAD")

	var wg sync.WaitGroup
	var semaphore chan struct{}

	switch selectedOptionsAccount {
	case 1:
		if isMultiThread {
			if maxThread > len(files) {
				maxThread = len(files)
			}

			semaphore = make(chan struct{}, maxThread)
			for _, file := range files {
				wg.Add(1)
				go processAccountMultiThread(&semaphore, &wg, file)
			}
			wg.Wait()
		} else {
			if len(files) == 1 {
				processAccountSingleThread(files[0])
				return
			} else {
				for _, file := range files {
					processAccountSingleThread(file)
				}
			}
		}
	case 2:
		selectedAccount := selectAccount(files)

		if len(selectedAccount) == 0 {
			return
		}

		if len(selectedAccount) == 1 {
			processAccountSingleThread(files[selectedAccount[0]-1])
			return
		} else {
			if isMultiThread {
				if maxThread > len(selectedAccount) {
					maxThread = len(selectedAccount)
				}

				semaphore = make(chan struct{}, maxThread)

				for _, account := range selectedAccount {
					wg.Add(1)
					go processAccountMultiThread(&semaphore, &wg, files[account-1])
				}
				wg.Wait()
			} else {
				for _, account := range selectedAccount {
					processAccountSingleThread(files[account-1])
				}
			}
		}
	}
}

func processAccountMultiThread(semaphore *chan struct{}, wg *sync.WaitGroup, file fs.DirEntry) {
	defer wg.Done()
	*semaphore <- struct{}{}

	defer func() {
		<-*semaphore
	}()

	browser := initializeBrowser()

	defer browser.MustClose()

	client := &Client{
		phoneNumber: strings.TrimSuffix(file.Name(), ".json"),
		Browser:     browser,
	}

	helper.PrettyLog("info", fmt.Sprintf("| %s | Start Processing Account...", client.phoneNumber))

	client.selectProcess(file)

	helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished...", client.phoneNumber))
}

func processAccountSingleThread(file fs.DirEntry) {
	browser := initializeBrowser()

	defer browser.MustClose()

	client := &Client{
		phoneNumber: strings.TrimSuffix(file.Name(), ".json"),
		Browser:     browser,
	}

	helper.PrettyLog("info", fmt.Sprintf("| %s | Start Processing Account...", client.phoneNumber))

	client.selectProcess(file)

	helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished...", client.phoneNumber))
}

func (c *Client) processGetLocalStorageSession(passwordAccount string, sessionsPath string, country string) {
	var phone, otpCode string

	helper.PrettyLog("info", "Launch Browser...")

	page := c.Browser.MustPage()

	c.navigate(page, "https://web.telegram.org/a/")

	isStop := false
	for !isStop {
		// Click Login By Phone
		c.clickElement(page, "#auth-qr-form > div > button")

		// Input Country
		c.inputText(page, country, "#sign-in-phone-code")

		time.Sleep(1 * time.Second)

		// Select Country
		c.clickElement(page, "#auth-phone-number-form > div > form > div.DropdownMenu.CountryCodeInput > div.Menu.compact.CountryCodeInput > div.bubble.menu-container.custom-scroll.opacity-transition.fast.left.top.shown.open > div")

		helper.PrettyLog("info", fmt.Sprintf("Selected Country: %s", country))
		// Input Number In Terminal
		phone = helper.InputTerminal("Input Phone Number: ")
		phone = strings.ReplaceAll(phone, " ", "")

		c.phoneNumber = phone

		if strings.Contains(phone, "+") {
			phone = strings.TrimPrefix(phone, "+")
		}

		if strings.HasPrefix(phone, "62") {
			phone = strings.TrimPrefix(phone, "62")
		}

		time.Sleep(1 * time.Second)

		// Input Phone Number
		c.inputText(page, phone, "#sign-in-phone-number")

		helper.PrettyLog("info", "Checking Your Number...")

		time.Sleep(1 * time.Second)

		// Click Next
		c.clickElement(page, "#auth-phone-number-form > div > form > button:nth-child(4)")

		time.Sleep(2 * time.Second)

		// Get Validation
		if c.checkElement(page, "#sign-in-code") {
			isStop = true
		} else {
			helper.PrettyLog("error", "Phone Number Invalid, Please Try Again...")
			page.MustReload()
			page.MustWaitLoad()
			continue
		}
	}

	isStop = false
	for !isStop {
		// Input Otp In Terminal
		otpCode = helper.InputTerminal("Input Otp Code: ")
		otpCode = strings.ReplaceAll(otpCode, " ", "")

		if len(otpCode) < 5 || len(otpCode) > 5 {
			helper.PrettyLog("error", "Otp Code Must 5 Digit Number, Please Try Again...")
			continue
		}

		time.Sleep(1 * time.Second)

		// Input Otp Code
		c.inputText(page, otpCode, "#sign-in-code")

		helper.PrettyLog("info", "Checking Otp Code...")

		time.Sleep(2 * time.Second)

		// Get Validation
		if c.getText(page, "#auth-code-form > div > div.input-group.with-label > label") == "Invalid code." {
			helper.PrettyLog("error", "Otp Code Invalid, Please Input Correct Otp Code...")
			c.removeTextFormInput(page, "#sign-in-code")
			continue
		}

		isStop = true
	}

	// Check Account Have Password Or Not
	isHavePassword := c.checkElement(page, "#sign-in-password")

	if isHavePassword {
		isStop = false
		for !isStop {
			if passwordAccount == "" {
				passwordAccount = strings.TrimSpace(helper.InputTerminal("Input Password: "))
			}
			// Input Password
			c.inputText(page, passwordAccount, "#sign-in-password")

			helper.PrettyLog("info", "Checking Account Password...")

			// Click Next
			c.clickElement(page, "form > button")

			time.Sleep(2 * time.Second)

			if c.getText(page, "#auth-password-form > div > form > div > label") == "Incorrect password." {
				helper.PrettyLog("error", "Password Is Incorrect, Please Input Correct Password...")
				c.removeTextFormInput(page, "#sign-in-password")
				passwordAccount = ""
				continue
			}

			isStop = true
		}
	} else {
		helper.PrettyLog("warning", fmt.Sprintf("Phone Number %v Not Have Password...", phone))
	}

	helper.PrettyLog("success", fmt.Sprintf("Login Phone Number %v Successfully, Sleep 5s Before Get Local Storage...", phone))

	time.Sleep(5 * time.Second)

	c.navigate(page, "https://web.telegram.org/k/")

	// Extract local storage data
	localStorageData := page.MustEval(`() => JSON.stringify(localStorage);`).String()

	var telegramData map[string]interface{}

	if err := json.Unmarshal([]byte(localStorageData), &telegramData); err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to unmarshal localStorage data: %v", err))
		return
	}

	filePath := fmt.Sprintf("%s/%s.json", sessionsPath, c.phoneNumber)

	var existingData []map[string]interface{}

	// Check Folder Session
	if !helper.CheckFileOrFolder(fmt.Sprintf("%v", sessionsPath)) {
		os.Mkdir(fmt.Sprintf("%v", sessionsPath), 0755)
	}

	// Baca file JSON jika ada
	if helper.CheckFileOrFolder(filePath) {
		_, err := helper.ReadFileJson(filePath)
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("Failed to read file: %v", err))
			return
		}
	}

	// Tambahkan data baru ke data yang ada
	existingData = append(existingData, telegramData)

	// Simpan data ke file JSON
	if err := helper.SaveFileJson(filePath, existingData); err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to save file: %v", err))
		return
	}

	return
}

func (c *Client) processSetLocalStorage(page *rod.Page, file fs.DirEntry) bool {
	account, err := helper.ReadFileJson(filepath.Join(localStoragePath, file.Name()))
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to read file %s: %v", c.phoneNumber, file.Name(), err))
		return false
	}

	// Membuka halaman kosong terlebih dahulu
	c.navigate(page, "https://web.telegram.org/k/")

	page.MustWaitLoad()

	time.Sleep(2 * time.Second)

	// Evaluasi JavaScript untuk menyimpan data ke localStorage
	switch v := account.(type) {
	case []map[string]interface{}:
		// Jika data adalah array of maps
		for _, acc := range v {
			for key, value := range acc {
				page.Eval(fmt.Sprintf(`localStorage.setItem('%s', '%s');`, key, value))
			}
		}
	case map[string]interface{}:
		// Jika data adalah single map
		for key, value := range v {
			page.Eval(fmt.Sprintf(`localStorage.setItem('%s', '%s');`, key, value))
		}
	default:
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to Evaluate Local Storage: Unknown Data Type", c.phoneNumber))
		return false
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Local storage successfully set | Check Login Status...", c.phoneNumber))

	page.MustReload()
	page.MustWaitLoad()

	time.Sleep(3 * time.Second)

	isSessionExpired := c.checkElement(page, "#auth-pages > div > div.tabs-container.auth-pages__container > div.tabs-tab.page-signQR.active > div > div.input-wrapper > button")

	if isSessionExpired {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Session Expired Or Account Banned, Please Check Your Account...", c.phoneNumber))

		helper.PrettyLog("info", fmt.Sprintf("| %s | Move Session File To Expired Folder | You Can Try Get Local Storage Again After Check Account...", c.phoneNumber))

		if !helper.CheckFileOrFolder(localStorageExpiredPath) {
			os.Mkdir(localStorageExpiredPath, 0755)
		}

		os.Rename(filepath.Join(localStoragePath, file.Name()), filepath.Join(localStorageExpiredPath, file.Name()))

		return false
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Login successfully | Sleep 5s Before Navigate...", c.phoneNumber))

	time.Sleep(5 * time.Second)

	helper.PrettyLog("info", fmt.Sprintf("| %s | Navigating Telegram...", c.phoneNumber))

	return true
}

func (c *Client) processFreeRoam(file fs.DirEntry) {
	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	helper.InputTerminal("Just press enter to next account or completing free roam...")

	return
}

func (c *Client) processGetQueryData(file fs.DirEntry) {
	botUsername := config.String(fmt.Sprintf("BOT_LIST.%v.%s.BOT_USERNAME", indexBotList, selectedBotList))

	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Search Bot
	c.searchBot(page, botUsername)

	// Click Launch App
	c.clickElement(page, "div.new-message-bot-commands")

	c.popupLaunchBot(page)

	time.Sleep(2 * time.Second)

	isIframe := c.checkElement(page, ".payment-verification")

	if !isIframe {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed To Launch Bot: Iframe Not Detected", c.phoneNumber))
		return
	}

	iframe := page.MustElement(".payment-verification")

	iframePage := iframe.MustFrame()

	iframePage.MustWaitDOMStable()

	helper.PrettyLog("info", fmt.Sprintf("| %s | Process Get Session Local Storage...", c.phoneNumber))

	// Mengeksekusi JavaScript untuk mendapatkan nilai dari sessiontorage
	res, err := iframePage.Evaluate(rod.Eval(`() => {
			let initParams = sessionStorage.getItem("__telegram__initParams");
			if (initParams) {
				let parsedParams = JSON.parse(initParams);
				return parsedParams.tgWebAppData;
			}
		
			initParams = sessionStorage.getItem("telegram-apps/launch-params");
			if (initParams) {
				let parsedParams = JSON.parse(initParams);
				return parsedParams;
			}
		
			return null;
		}`))

	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to evaluate script: %v", c.phoneNumber, err))
		return
	}

	var queryData string

	if strings.Contains(res.Value.String(), "tgWebAppData=") {
		queryParamsString, err := helper.GetTextAfterKey(res.Value.String(), "tgWebAppData=")
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get text after key tgWebAppData=: %v", c.phoneNumber, err))
			return
		}

		queryData = queryParamsString
	} else {
		if res.Type == proto.RuntimeRemoteObjectTypeString {
			queryData = res.Value.String()
			helper.PrettyLog("success", fmt.Sprintf("| %s | Get Session Storage Successfully...", c.phoneNumber))
		} else {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Get Session Storage Failed...", c.phoneNumber))
			return
		}
	}

	if len(queryData) > 0 {
		if helper.CheckFileOrFolder(fmt.Sprintf("%s/%s", queryDataPath, botUsername)) {
			os.Mkdir(fmt.Sprintf("%s/%s", queryDataPath, botUsername), 0755)
		}

		filePath := fmt.Sprintf("%s/%s/%s", queryDataPath, botUsername, fmt.Sprintf("query_data_%s.txt", c.phoneNumber))

		if helper.CheckFileOrFolder(filePath) {
			os.Remove(filePath)
		}

		// Save Query Data
		helper.SaveFileTxt(filePath, queryData)

		if helper.CheckFileOrFolder(filePath) {
			helper.PrettyLog("success", fmt.Sprintf("| %s | Query Data Successfully Saved", c.phoneNumber))
		} else {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Query Data Failed Saved", c.phoneNumber))
		}
	} else {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed To Get Query Data", c.phoneNumber))
	}
}

func (c *Client) processGetDetailAccount(file fs.DirEntry) {
	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Search Bot
	c.searchBot(page, "userinfobot")

	// Send Message
	c.sendMessage(page, "/start", true)

	// Get Message
	message := c.getLastChat(page)

	result := make(map[string]string)

	// Pisahkan text menjadi baris-baris
	lines := strings.Split(message, "\n")

	// Iterasi setiap baris
	for _, line := range lines {
		// Jika baris mengandung ": ", kita pisah berdasarkan itu
		if strings.Contains(line, ": ") {
			parts := strings.SplitN(line, ": ", 2) // Split menjadi 2 bagian: kunci dan nilai
			key := parts[0]                        // Bagian pertama adalah kunci
			value := parts[1]                      // Bagian kedua adalah nilai
			result[key] = value                    // Masukkan ke dalam map
		} else if strings.HasPrefix(line, "@") {
			// Jika baris diawali dengan @, hapus @ dan simpan username
			result["username"] = strings.TrimPrefix(strings.TrimSpace(line), "@")
		}
	}

	filePath := fmt.Sprintf("%s/detail_account_%s.json", detailAccountPath, c.phoneNumber)

	helper.SaveFileJson(filePath, result)

	if helper.CheckFileOrFolder(filePath) {
		helper.PrettyLog("success", fmt.Sprintf(fmt.Sprintf("| %s | Detail Account Successfully Saved", c.phoneNumber)))
	} else {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Detail Account Failed Saved", c.phoneNumber))
	}
}

func (c *Client) processStartBotWithAutoRef(file fs.DirEntry) {

	refUrl := config.String(fmt.Sprintf("BOT_LIST.%v.%s.REF_URL", indexBotList, selectedBotList))
	selectors := config.Strings(fmt.Sprintf("BOT_LIST.%v.%s.CLICKABLE_SELECTOR", indexBotList, selectedBotList))

	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Search Bot
	c.searchBot(page, "+42777")

	// Send Message Ref Url
	c.sendMessage(page, refUrl, false)

	time.Sleep(3 * time.Second)

	// Click Launch App
	c.clickElement(page, fmt.Sprintf(`a.anchor-url[href="%v"]`, refUrl))

	c.popupLaunchBot(page)

	time.Sleep(3 * time.Second)

	isIframe := c.checkElement(page, ".payment-verification")

	if isIframe {
		helper.PrettyLog("success", "Launch Bot")

		iframe := page.MustElement(".payment-verification")

		iframePage := iframe.MustFrame()

		iframePage.MustWaitDOMStable()

		helper.PrettyLog("info", fmt.Sprintf("| %s | Process Clicking Selector Bot...", c.phoneNumber))

		for _, selector := range selectors {
			c.clickElement(iframePage, selector)
			time.Sleep(2 * time.Second)
			iframePage.MustWaitNavigation()
		}

		helper.PrettyLog("success", fmt.Sprintf("| %s | Clicking Selector Bot Completed...", c.phoneNumber))

		time.Sleep(5 * time.Second)
	}
}

func (c *Client) processSetAccountUsername(file fs.DirEntry) {
	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Click Ripple Button
	c.clickElement(page, "#column-left > div > div > div.sidebar-header.can-have-forum > div.sidebar-header__btn-container > button")

	time.Sleep(1 * time.Second)

	isSetting := c.gotoSetting(page)

	if !isSetting {
		return
	}

	time.Sleep(1 * time.Second)

	// Click Edit Profile
	c.clickElement(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.settings-container.profile-container.is-collapsed.active > div.sidebar-header > button:nth-child(3)")

	time.Sleep(1 * time.Second)

	firstName := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(2) > div.sidebar-left-section > div > div.input-wrapper > div:nth-child(1) > div.input-field-input")

	lastName := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(2) > div.sidebar-left-section > div > div.input-wrapper > div:nth-child(2) > div.input-field-input")

	currentUsername := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > input")

	isComplete := false
	for !isComplete {

		helper.PrettyLog("info", fmt.Sprintf("| %s | Full Name: %s", c.phoneNumber, fmt.Sprintf("%s %s", firstName, lastName)))

		helper.PrettyLog("info", fmt.Sprintf("| %s | Current Username: %s", c.phoneNumber, currentUsername))

		if currentUsername != "" {
			helper.PrettyLog("info", "Are You Sure Want To Change Current Username?")

			makeSure := strings.TrimSpace(helper.InputTerminal("Input Choice (y/n) (default = Next Account): "))

			if makeSure != "y" || makeSure != "Y" {
				return
			}
		}

		// Input Username
		username := strings.TrimSpace(helper.InputTerminal("Input Username: "))

		helper.PrettyLog("info", fmt.Sprintf("| %s | Checking Username...", c.phoneNumber))

		c.removeTextFormInput(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > input")

		c.inputText(page, username, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > input")

		time.Sleep(2 * time.Second)

		isUsernameAvailable := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > label > span")

		// Check Username
		if isUsernameAvailable == "Username is already taken" || isUsernameAvailable == "Username is invalid" {
			helper.PrettyLog("error", fmt.Sprintf("| %s | %s, Try Another Username", c.phoneNumber, isUsernameAvailable))
			c.removeTextFormInput(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > input")
			continue
		}

		helper.PrettyLog("info", fmt.Sprintf("| %s | Username: %s Available", c.phoneNumber, username))

		time.Sleep(1 * time.Second)

		c.clickElement(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > button")

		time.Sleep(1 * time.Second)

		helper.PrettyLog("success", fmt.Sprintf("| %s | Username Successfully Set", c.phoneNumber))

		isComplete = true
	}
}

func (c *Client) processSetFirstName(file fs.DirEntry) {
	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Click Ripple Button
	c.clickElement(page, "#column-left > div > div > div.sidebar-header.can-have-forum > div.sidebar-header__btn-container > button")

	time.Sleep(1 * time.Second)

	isSetting := c.gotoSetting(page)

	if !isSetting {
		return
	}

	time.Sleep(1 * time.Second)

	// Click Edit Profile
	c.clickElement(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.settings-container.profile-container.is-collapsed.active > div.sidebar-header > button:nth-child(3)")

	time.Sleep(1 * time.Second)

	firstName := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(2) > div.sidebar-left-section > div > div.input-wrapper > div:nth-child(1) > div.input-field-input")

	lastName := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(2) > div.sidebar-left-section > div > div.input-wrapper > div:nth-child(2) > div.input-field-input")

	isComplete := false
	for !isComplete {

		helper.PrettyLog("info", fmt.Sprintf("| %s | Current Full Name: %s", c.phoneNumber, fmt.Sprintf("%s %s", firstName, lastName)))

		// Input First Name
		firstName := strings.TrimSpace(helper.InputTerminal("Input New First Name: "))

		c.inputText(page, firstName, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > input")

		time.Sleep(2 * time.Second)

		c.clickElement(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > button")

		time.Sleep(1 * time.Second)

		helper.PrettyLog("success", fmt.Sprintf("| %s | First Name Successfully Set", c.phoneNumber))

		isComplete = true
	}
}

func (c *Client) processSetLastName(file fs.DirEntry) {
	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Click Ripple Button
	c.clickElement(page, "#column-left > div > div > div.sidebar-header.can-have-forum > div.sidebar-header__btn-container > button")

	time.Sleep(1 * time.Second)

	isSetting := c.gotoSetting(page)

	if !isSetting {
		return
	}

	time.Sleep(1 * time.Second)

	// Click Edit Profile
	c.clickElement(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.settings-container.profile-container.is-collapsed.active > div.sidebar-header > button:nth-child(3)")

	time.Sleep(1 * time.Second)

	firstName := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(2) > div.sidebar-left-section > div > div.input-wrapper > div:nth-child(1) > div.input-field-input")

	lastName := c.getText(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrolled-bottom.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(2) > div.sidebar-left-section > div > div.input-wrapper > div:nth-child(2) > div.input-field-input")

	isComplete := false
	for !isComplete {

		helper.PrettyLog("info", fmt.Sprintf("| %s | Current Full Name: %s", c.phoneNumber, fmt.Sprintf("%s %s", firstName, lastName)))

		if lastName != "" {
			// Input First Name
			lastName = strings.TrimSpace(helper.InputTerminal("Input New Last Name: "))
		}

		c.inputText(page, lastName, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrolled-top.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > div > div:nth-child(3) > div.sidebar-left-section > div > div.input-wrapper > div > input")

		time.Sleep(2 * time.Second)

		c.clickElement(page, "#column-left > div > div.tabs-tab.sidebar-slider-item.scrollable-y-bordered.edit-profile-container.active > div.sidebar-content > button")

		time.Sleep(1 * time.Second)

		helper.PrettyLog("success", fmt.Sprintf("| %s | Last Name Successfully Set", c.phoneNumber))

		isComplete = true
	}
}

func (c *Client) processAutoSubscribeChannel(file fs.DirEntry) {
	page := c.Browser.MustPage()

	// Set Local Storage
	isLogin := c.processSetLocalStorage(page, file)

	if !isLogin {
		return
	}

	// Search
	c.searchBot(page, channelUsername)

	helper.PrettyLog("info", fmt.Sprintf("| %s | Subscribing %s Channel...", c.phoneNumber, channelUsername))

	isSubscribeButton := c.checkElement(page, "#column-center > div > div > div.sidebar-header.topbar.has-avatar > div.chat-info-container > div.chat-utils > button.btn-primary.btn-color-primary.chat-join.rp")

	if isSubscribeButton {
		// Click Subscribe
		c.clickElement(page, "#column-center > div > div > div.sidebar-header.topbar.has-avatar > div.chat-info-container > div.chat-utils > button.btn-primary.btn-color-primary.chat-join.rp")

		helper.PrettyLog("success", fmt.Sprintf("| %s | Subscribing %s Channel Successfully...", c.phoneNumber, channelUsername))
	} else {
		helper.PrettyLog("success", fmt.Sprintf("| %s | Already Subscribing %s Channel...", c.phoneNumber, channelUsername))
	}
}
