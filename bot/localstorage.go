package bot

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"telegram-web/core"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func SessionLocalStorage(passwordAccount string, url string, sessionsPath string) {
	var phone, otpCode string

	launchOptions := launcher.New().
		Headless(true).
		MustLaunch()

	browser := rod.New().ControlURL(launchOptions).MustConnect()

	defer browser.MustClose()

	page := browser.MustPage()

	core.Navigate(page, url)

	// Click Login By Phone
	core.ClickElement(page, "div.input-wrapper > button:nth-child(1)")

	isStop := false
	var countLimit int

	for !isStop {
		// Limit Check
		if countLimit == 3 {
			page.MustReload()
			page.MustWaitLoad()
		}

		if countLimit == 5 {
			helper.PrettyLog("error", "Limit Check Show Region For Input Number Reached, Please Check Your Connection...")
			return
		}

		// Check Region
		isShowRegion := core.GetText(page, "div.input-field.input-select > div.input-field-input > span")

		if isShowRegion == "Indonesia" {
			isStop = true
		} else {
			time.Sleep(3 * time.Second)
			countLimit++
		}
	}

	isStop = false
	for !isStop {
		// Input Number In Terminal
		phone = strings.TrimSpace(helper.InputTerminal("Input Phone Number :"))

		if !strings.Contains(phone, "+62") {
			phone = "+62" + phone
		}

		// Input Phone Number
		core.InputText(page, phone, "div.input-field.input-field-phone > div.input-field-input")

		// Click Next
		core.ClickElement(page, "button.btn-primary.btn-color-primary.rp")

		time.Sleep(3 * time.Second)

		isPhoneValid := core.GetText(page, "div.input-field.input-field-phone > label > span")

		if isPhoneValid == "Phone Number Invalid" {
			helper.PrettyLog("error", "Phone Number Invalid, Please Try Again...")
			page.MustReload()
			page.MustWaitLoad()
			continue
		} else {
			isStop = true
		}
	}

	isStop = false
	for !isStop {
		// Input Otp In Terminal
		otpCode = strings.TrimSpace(helper.InputTerminal("Input Otp Code :"))

		if len(otpCode) < 5 || len(otpCode) > 5 {
			helper.PrettyLog("error", "Otp Code Must 5 Digit Number, Please Try Again...")
			continue
		}

		// Input Otp Code
		core.InputText(page, otpCode, "div.tabs-tab.page-authCode.active > div > div.input-wrapper > div > input")

		time.Sleep(3 * time.Second)

		helper.PrettyLog("info", "Check Otp Code...")

		// Get Validation
		isOtpValid := core.GetText(page, "div.tabs-tab.page-authCode.active > div > div.input-wrapper > div > label > span")

		if isOtpValid == "Invalid code" {
			helper.PrettyLog("error", "Otp Code Invalid, Please Try Again...")
			page.MustReload()
			page.MustWaitLoad()
			continue
		} else {
			isStop = true
		}
	}

	helper.PrettyLog("info", "Check Account Password...")

	// Check Account Have Password Or Not
	isHavePassword := core.CheckElement(page, "div.tabs-tab.page-password.active > div > div.input-wrapper > div > input.input-field-input.is-empty")

	if isHavePassword {
		// Input Password
		core.InputText(page, passwordAccount, "#auth-pages > div > div.tabs-container.auth-pages__container > div.tabs-tab.page-password.active > div > div.input-wrapper > div > input.input-field-input.is-empty")

		// Click Next
		core.ClickElement(page, "#auth-pages > div > div.tabs-container.auth-pages__container > div.tabs-tab.page-password.active > div > div.input-wrapper > button")
	} else {
		helper.PrettyLog("warning", fmt.Sprintf("Account %v Not Have Password...", phone))
	}

	helper.PrettyLog("success", fmt.Sprintf("Login Account %v Successfully...", phone))

	page.MustReload()

	page.MustWaitLoad()

	time.Sleep(3 * time.Second)

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
