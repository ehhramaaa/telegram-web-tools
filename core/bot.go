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

var sessionStorage []string

func processGetLocalStorage(browser *rod.Browser, passwordAccount string, loginUrl string, sessionsPath string, country string) {
	var phone, otpCode string

	page := browser.MustPage()

	navigate(page, loginUrl)

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	isStop := false
	for !isStop {
		// Click Login By Phone
		clickElement(page, "#auth-qr-form > div > button")

		// Input Country
		inputText(page, country, "#sign-in-phone-code")

		time.Sleep(1 * time.Second)

		// Select Country
		clickElement(page, "#auth-phone-number-form > div > form > div.DropdownMenu.CountryCodeInput > div.Menu.compact.CountryCodeInput > div.bubble.menu-container.custom-scroll.opacity-transition.fast.left.top.shown.open > div")

		// Input Number In Terminal
		phone = strings.TrimSpace(helper.InputTerminal("Input Phone Number (Without +): "))
		if strings.Contains(phone, "+") {
			phone = strings.TrimPrefix(phone, "+")
		}

		// Input Phone Number
		inputText(page, phone, "#sign-in-phone-number")

		time.Sleep(1 * time.Second)

		// Click Next
		clickElement(page, "#auth-phone-number-form > div > form > button:nth-child(4)")

		time.Sleep(3 * time.Second)

		isPhoneValid := getText(page, "#auth-phone-number-form > div > form > div.input-group.touched.with-label > label")

		if isPhoneValid == "Invalid phone number." {
			helper.PrettyLog("error", "Phone Number Invalid, Please Try Again...")
			page.MustReload()
			page.MustWaitLoad()
			continue
		}

		if checkElement(page, "#sign-in-code") {
			// Input Otp In Terminal
			otpCode = strings.TrimSpace(helper.InputTerminal("Input Otp Code: "))

			if len(otpCode) < 5 || len(otpCode) > 5 {
				helper.PrettyLog("error", "Otp Code Must 5 Digit Number, Please Try Again...")
				continue
			}

			time.Sleep(1 * time.Second)

			// Input Otp Code
			inputText(page, otpCode, "#sign-in-code")

			time.Sleep(2 * time.Second)

			helper.PrettyLog("info", "Check Otp Code...")

			// Get Validation
			isOtpValid := getText(page, "#auth-code-form > div > div.input-group.with-label > label")

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
	isHavePassword := checkElement(page, "#sign-in-password")

	if isHavePassword {
		// Input Password
		inputText(page, passwordAccount, "#sign-in-password")

		// Click Next
		clickElement(page, "form > button")
	} else {
		helper.PrettyLog("warning", fmt.Sprintf("Account %v Not Have Password...", phone))
	}

	helper.PrettyLog("success", fmt.Sprintf("Login Account %v Successfully, Sleep 5s Before Get Local Storage...", phone))

	time.Sleep(5 * time.Second)

	navigate(page, "https://web.telegram.org/k/")

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

func processGetQueryData(browser *rod.Browser, file fs.DirEntry, localStoragePath string, botUsername string) {
	account, err := helper.ReadFileJson(filepath.Join(localStoragePath, file.Name()))
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to read file %s: %v", phoneNumber, file.Name(), err))
		return
	}

	// Membuka halaman kosong terlebih dahulu
	page := browser.MustPage()
	navigate(page, "https://web.telegram.org/k/")

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
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to Evaluate Local Storage: Unknown Data Type", phoneNumber))
		return
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Local storage successfully set. Navigating to Telegram Web...", phoneNumber))

	time.Sleep(2 * time.Second)

	// Reload Page
	page.MustReload()
	page.MustWaitLoad()

	// Search Bot
	searchBot(page, botUsername)

	// Click Launch App
	clickElement(page, "div.new-message-bot-commands")

	popupLaunchBot(page)

	time.Sleep(2 * time.Second)

	isIframe := checkElement(page, ".payment-verification")

	if !isIframe {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed To Launch Bot: Iframe Not Detected", phoneNumber))
		return
	}

	iframe := page.MustElement(".payment-verification")

	iframePage := iframe.MustFrame()

	iframePage.MustWaitDOMStable()

	helper.PrettyLog("info", fmt.Sprintf("| %s | Process Get Session Local Storage...", phoneNumber))

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
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to evaluate script: %v", phoneNumber, err))
		return
	}

	var queryData string

	if strings.Contains(res.Value.String(), "tgWebAppData=") {
		queryParamsString, err := helper.GetTextAfterKey(res.Value.String(), "tgWebAppData=")
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get text after key tgWebAppData=: %v", phoneNumber, err))
			return
		}

		queryData = queryParamsString
	} else {
		if res.Type == proto.RuntimeRemoteObjectTypeString {
			queryData = res.Value.String()
			helper.PrettyLog("success", fmt.Sprintf("| %s | Get Session Storage Successfully...", phoneNumber))
		} else {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Get Session Storage Failed...", phoneNumber))
			return
		}
	}

	if len(queryData) > 0 {
		sessionStorage = append(sessionStorage, queryData)
	}
}

func processStartBotWithAutoRef(browser *rod.Browser, file fs.DirEntry, localStoragePath string, botUsername string, refUrl string) {
	account, err := helper.ReadFileJson(filepath.Join(localStoragePath, file.Name()))
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to read file %s: %v", phoneNumber, file.Name(), err))
		return
	}

	// Membuka halaman kosong terlebih dahulu
	page := browser.MustPage()
	navigate(page, "https://web.telegram.org/k/")

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
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to Evaluate Local Storage: Unknown Data Type", phoneNumber))
		return
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Local storage successfully set. Navigating to Telegram Web...", phoneNumber))

	time.Sleep(2 * time.Second)

	// Reload Page
	page.MustReload()
	page.MustWaitLoad()

	// Search Bot
	searchBot(page, "+42777")

	isBot := false
	// Send Message Ref Url
	sendMessage(page, refUrl, isBot)

	time.Sleep(3 * time.Second)

	// Click Launch App
	clickElement(page, fmt.Sprintf(`a.anchor-url[href="%v"]`, refUrl))

	popupLaunchBot(page)

	time.Sleep(3 * time.Second)

	isIframe := checkElement(page, ".payment-verification")

	if isIframe {
		helper.PrettyLog("success", "Launch Bot")

		iframe := page.MustElement(".payment-verification")

		iframePage := iframe.MustFrame()

		iframePage.MustWaitDOMStable()

		selectors := config.Strings("FIRST_LAUNCH_BOT_SELECTOR")

		helper.PrettyLog("info", fmt.Sprintf("| %s | Process Clicking Selector Bot...", phoneNumber))

		for _, selector := range selectors {
			clickElement(iframePage, selector)
			time.Sleep(2 * time.Second)
			iframePage.MustWaitDOMStable()
		}

		helper.PrettyLog("success", fmt.Sprintf("| %s | Clicking Selector Bot Completed...", phoneNumber))

		time.Sleep(5 * time.Second)
	}
}
