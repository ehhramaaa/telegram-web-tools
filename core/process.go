package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gookit/config/v2"
)

func (c *Client) setLocalStorage(page *rod.Page, file fs.DirEntry, localStoragePath string) {
	account, err := helper.ReadFileJson(filepath.Join(localStoragePath, file.Name()))
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to read file %s: %v", c.phoneNumber, file.Name(), err))
		return
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
		return
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Local storage successfully set. Navigating to Telegram Web...", c.phoneNumber))

	time.Sleep(2 * time.Second)
}

func (c *Client) processGetLocalStorage(passwordAccount string, sessionsPath string, country string) {
	var phone, otpCode string

	page := c.Browser.MustPage()

	c.navigate(page, "https://web.telegram.org/a/")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	isStop := false
	for !isStop {
		// Click Login By Phone
		c.clickElement(page, "#auth-qr-form > div > button")

		// Input Country
		c.inputText(page, country, "#sign-in-phone-code")

		time.Sleep(1 * time.Second)

		// Select Country
		c.clickElement(page, "#auth-phone-number-form > div > form > div.DropdownMenu.CountryCodeInput > div.Menu.compact.CountryCodeInput > div.bubble.menu-container.custom-scroll.opacity-transition.fast.left.top.shown.open > div")

		// TODO
		// Input Number In Terminal
		phone = strings.TrimSpace(helper.InputTerminal("Input Phone Number (Without +): "))
		if strings.Contains(phone, "+") {
			phone = strings.TrimPrefix(phone, "+")
		}

		// Input Phone Number
		c.inputText(page, phone, "#sign-in-phone-number")

		time.Sleep(1 * time.Second)

		// Click Next
		c.clickElement(page, "#auth-phone-number-form > div > form > button:nth-child(4)")

		time.Sleep(3 * time.Second)

		isPhoneValid := c.getText(page, "#auth-phone-number-form > div > form > div.input-group.touched.with-label > label")

		if isPhoneValid == "Invalid phone number." {
			helper.PrettyLog("error", "Phone Number Invalid, Please Try Again...")
			page.MustReload()
			page.MustWaitLoad()
			continue
		}

		if c.checkElement(page, "#sign-in-code") {
			// Input Otp In Terminal
			otpCode = strings.TrimSpace(helper.InputTerminal("Input Otp Code: "))

			if len(otpCode) < 5 || len(otpCode) > 5 {
				helper.PrettyLog("error", "Otp Code Must 5 Digit Number, Please Try Again...")
				continue
			}

			time.Sleep(1 * time.Second)

			// Input Otp Code
			c.inputText(page, otpCode, "#sign-in-code")

			time.Sleep(2 * time.Second)

			helper.PrettyLog("info", "Check Otp Code...")

			// Get Validation
			isOtpValid := c.getText(page, "#auth-code-form > div > div.input-group.with-label > label")

			if isOtpValid == "Invalid code." {
				helper.PrettyLog("error", "Otp Code Invalid, Please Try Again...")
				page.MustReload()
				page.MustWaitLoad()
				continue
			} else {
				isStop = true
			}
		} else {
			helper.PrettyLog("warning", "Selector Input Otp Not Found")
		}
	}

	helper.PrettyLog("info", "Check Account Password...")

	// Check Account Have Password Or Not
	isHavePassword := c.checkElement(page, "#sign-in-password")

	if isHavePassword {
		if passwordAccount == "" {
			passwordAccount = strings.TrimSpace(helper.InputTerminal("Input Password: "))
		}
		// Input Password
		c.inputText(page, passwordAccount, "#sign-in-password")

		// Click Next
		c.clickElement(page, "form > button")
	} else {
		helper.PrettyLog("warning", fmt.Sprintf("Account %v Not Have Password...", phone))
	}

	helper.PrettyLog("success", fmt.Sprintf("Login Account %v Successfully, Sleep 5s Before Get Local Storage...", phone))

	time.Sleep(5 * time.Second)

	c.navigate(page, "https://web.telegram.org/k/")

	// Extract local storage data
	localStorageData := page.MustEval(`() => JSON.stringify(localStorage);`).String()

	var telegramData map[string]interface{}

	if err := json.Unmarshal([]byte(localStorageData), &telegramData); err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to unmarshal localStorage data: %v", err))
		return
	}

	filePath := fmt.Sprintf("%s/%s.json", sessionsPath, phone)

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

	helper.PrettyLog("success", fmt.Sprintf("Data berhasil disimpan ke %s", filePath))
}

func (c *Client) processGetQueryData(file fs.DirEntry, localStoragePath string, queryDataPath string, botUsername string) {
	page := c.Browser.MustPage()

	// Set Local Storage
	c.setLocalStorage(page, file, localStoragePath)

	// Reload Page
	page.MustReload()
	page.MustWaitLoad()

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

func (c *Client) processStartBotWithAutoRef(file fs.DirEntry, localStoragePath string, botUsername string, refUrl string) {
	page := c.Browser.MustPage()

	// Set Local Storage
	c.setLocalStorage(page, file, localStoragePath)

	// Reload Page
	page.MustReload()
	page.MustWaitLoad()

	// Search Bot
	c.searchBot(page, "+42777")

	isBot := false
	// Send Message Ref Url
	c.sendMessage(page, refUrl, isBot)

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

		selectors := config.Strings("FIRST_LAUNCH_BOT_SELECTOR")

		helper.PrettyLog("info", fmt.Sprintf("| %s | Process Clicking Selector Bot...", c.phoneNumber))

		for _, selector := range selectors {
			c.clickElement(iframePage, selector)
			time.Sleep(2 * time.Second)
			iframePage.MustWaitDOMStable()
		}

		helper.PrettyLog("success", fmt.Sprintf("| %s | Clicking Selector Bot Completed...", c.phoneNumber))

		time.Sleep(5 * time.Second)
	}
}
