package bot

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"telegram-web/core"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gookit/config/v2"
)

type SessionStorage struct {
	Username  string
	Id        string
	FirstName string
	LastName  string
	QueryData string
}

func parseText(text string) map[string]string {
	// Buat map untuk menyimpan hasil parsing
	result := make(map[string]string)

	// Pisahkan text menjadi baris-baris
	lines := strings.Split(text, "\n")

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

	return result
}

func GetAllAccount(botUsername string, refUrl string, isAutoRef bool, localStoragePath string, queryDataPath string, files []fs.DirEntry) {
	for _, file := range files {
		fmt.Printf("<=====================[%v]=====================>\n", file.Name())

		var sessionStorage []SessionStorage

		launchOptions := launcher.New().
			Headless(true).
			MustLaunch()

		browser := rod.New().ControlURL(launchOptions).MustConnect()

		defer browser.MustClose()

		// Membaca setiap file JSON
		account, err := helper.ReadFileJson(filepath.Join(localStoragePath, file.Name()))
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("Failed to read file %s: %v", file.Name(), err))
			continue
		}

		// Membuka halaman kosong terlebih dahulu
		page := browser.MustPage()
		core.Navigate(page, "https://web.telegram.org/k/")

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
			helper.PrettyLog("error", "Unknown data type")
		}

		helper.PrettyLog("success", "Local storage successfully set. Navigating to Telegram Web...")

		time.Sleep(3 * time.Second)

		// Reload Page
		page.MustReload()
		page.MustWaitLoad()

		// Search UserInfoBot
		SearchBot(page, "userinfobot")

		// Send Message Get Detail Account
		SendMessage(page, "/start")

		time.Sleep(3 * time.Second)

		// Get Text From Chat
		text := GetLastChat(page)

		// Save Detail Account
		detailAcc := parseText(text)

		// Send Message Ref Url
		SendMessage(page, refUrl)

		// Click Launch App
		core.ClickElement(page, fmt.Sprintf(`a.anchor-url[href="%v"]`, refUrl))

		// Click Popup Launch If Found
		isPopupLaunch := core.CheckElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")

		if isPopupLaunch {
			core.ClickElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")
		}

		time.Sleep(3 * time.Second)

		isIframe := core.CheckElement(page, ".payment-verification")

		if isIframe {
			helper.PrettyLog("success", "Launch Bot")

			iframe := page.MustElement(".payment-verification")

			iframePage := iframe.MustFrame()

			iframePage.MustWaitDOMStable()

			if isAutoRef {
				selectors := config.Strings("bot.selector")

				helper.PrettyLog("info", "Process Clicking Selector Bot...")

				for _, selector := range selectors {
					core.ClickElement(iframePage, selector)
					time.Sleep(2 * time.Second)
					iframePage.MustWaitDOMStable()
				}

				helper.PrettyLog("success", "Clicking Selector Bot Completed...")

				time.Sleep(3 * time.Second)
			}

			helper.PrettyLog("info", "Process Get Session Local Storage...")

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
				helper.PrettyLog("error", fmt.Sprintf("Failed to evaluate script: %v", err))
			}

			var queryParams string

			if strings.Contains(res.Value.String(), "tgWebAppData=") {
				queryParamsString, err := helper.GetTextAfterKey(res.Value.String(), "tgWebAppData=")
				if err != nil {
					helper.PrettyLog("error", err.Error())
					return
				}

				queryParams = queryParamsString
			} else {
				if res.Type == proto.RuntimeRemoteObjectTypeString {
					queryParams = res.Value.String()
					helper.PrettyLog("success", "Get Session Storage Successfully")
				} else {
					helper.PrettyLog("warning", "Get Session Storage Failed, Continue Next Account")
					continue
				}
			}

			if len(queryParams) > 0 {
				// Check Folder Query Data
				if !helper.CheckFileOrFolder(fmt.Sprintf("%v/%v", queryDataPath, botUsername)) {
					os.Mkdir(fmt.Sprintf("%v/%v", queryDataPath, botUsername), 0755)
				}

				// Check File Query Data
				if helper.CheckFileOrFolder(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, file.Name())) {
					os.Remove(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, file.Name()))
				}

				sessionStorage = append(sessionStorage, SessionStorage{
					Username:  detailAcc["username"],
					Id:        detailAcc["Id"],
					FirstName: detailAcc["First"],
					LastName:  detailAcc["Last"],
					QueryData: queryParams,
				})

				// Save Query Data
				helper.SaveFileJson(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, file.Name()), sessionStorage)

				helper.PrettyLog("success", fmt.Sprintf("Query data %v berhasil disimpan", file.Name()))

				browser.MustClose()

				randomNumber := helper.RandomNumber(5, 15)

				helper.PrettyLog("info", fmt.Sprintf("Rest For %v Second....\n", randomNumber))

				time.Sleep(time.Duration(randomNumber) * time.Second)
			}
		}
	}
}

func GetOneAccount(session string, botUsername string, refUrl string, isAutoRef bool, localStoragePath string, queryDataPath string) {
	fmt.Printf("<=====================[%v]=====================>\n", session)

	var sessionStorage []SessionStorage

	launchOptions := launcher.New().
		Headless(true).
		MustLaunch()

	browser := rod.New().ControlURL(launchOptions).MustConnect()

	defer browser.MustClose()

	// Membaca setiap file JSON
	account, err := helper.ReadFileJson(filepath.Join(localStoragePath, session))
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read file %s: %v", session, err))
	}

	// Membuka halaman kosong terlebih dahulu
	page := browser.MustPage()
	core.Navigate(page, "https://web.telegram.org/k/")

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
		helper.PrettyLog("error", "Unknown data type")
	}

	helper.PrettyLog("success", "Local storage successfully set. Navigating to Telegram Web...")

	time.Sleep(3 * time.Second)

	// Reload Page
	page.MustReload()
	page.MustWaitLoad()

	// Search UserInfoBot
	SearchBot(page, "userinfobot")

	// Send Message Get Detail Account
	SendMessage(page, "/start")

	time.Sleep(3 * time.Second)

	// Get Text From Chat
	text := GetLastChat(page)

	// Save Detail Account
	detailAcc := parseText(text)

	// Send Message Ref Url
	SendMessage(page, refUrl)

	// Click Launch App
	core.ClickElement(page, fmt.Sprintf(`a.anchor-url[href="%v"]`, refUrl))

	// Click Popup Launch If Found
	isPopupLaunch := core.CheckElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")

	if isPopupLaunch {
		core.ClickElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")
	}

	time.Sleep(3 * time.Second)

	isIframe := core.CheckElement(page, ".payment-verification")

	if isIframe {
		helper.PrettyLog("success", "Launch Bot")

		iframe := page.MustElement(".payment-verification")

		iframePage := iframe.MustFrame()

		iframePage.MustWaitDOMStable()

		if isAutoRef {
			selectors := config.Strings("bot.selector")

			helper.PrettyLog("info", "Process Clicking Selector Bot...")

			for _, selector := range selectors {
				core.ClickElement(iframePage, selector)
				time.Sleep(2 * time.Second)
				iframePage.MustWaitDOMStable()
			}

			helper.PrettyLog("success", "Clicking Selector Bot Completed...")

			time.Sleep(3 * time.Second)
		}

		helper.PrettyLog("info", "Process Get Session Local Storage...")

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
			helper.PrettyLog("error", fmt.Sprintf("Failed to evaluate script: %v", err))
		}

		var queryParams string

		if strings.Contains(res.Value.String(), "tgWebAppData=") {
			queryParamsString, err := helper.GetTextAfterKey(res.Value.String(), "tgWebAppData=")
			if err != nil {
				helper.PrettyLog("error", err.Error())
				return
			}

			queryParams = queryParamsString
		} else {
			if res.Type == proto.RuntimeRemoteObjectTypeString {
				queryParams = res.Value.String()
				helper.PrettyLog("success", "Get Session Storage Successfully")
			} else {
				helper.PrettyLog("warning", "Get Session Storage Failed, Continue Next Account")
			}
		}

		if len(queryParams) > 0 {
			// Check Folder Query Data
			if !helper.CheckFileOrFolder(fmt.Sprintf("%v/%v", queryDataPath, botUsername)) {
				os.Mkdir(fmt.Sprintf("%v/%v", queryDataPath, botUsername), 0755)
			}

			// Check File Query Data
			if helper.CheckFileOrFolder(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, session)) {
				os.Remove(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, session))
			}

			sessionStorage = append(sessionStorage, SessionStorage{
				Username:  detailAcc["username"],
				Id:        detailAcc["Id"],
				FirstName: detailAcc["First"],
				LastName:  detailAcc["Last"],
				QueryData: queryParams,
			})

			// Save Query Data
			helper.SaveFileJson(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, session), sessionStorage)

			helper.PrettyLog("success", fmt.Sprintf("Query data %v berhasil disimpan", session))

			browser.MustClose()

			randomNumber := helper.RandomNumber(5, 15)

			helper.PrettyLog("info", fmt.Sprintf("Rest For %v Second....\n", randomNumber))

			time.Sleep(time.Duration(randomNumber) * time.Second)
		}
	}
}

func MergeData(botUsername string, queryDataPath string, mergePath string, mergeFileName string, files []fs.DirEntry) {
	// Inisialisasi map untuk menampung hasil penggabungan data
	var mergedData []SessionStorage

	for _, file := range files {
		// Baca data JSON dari file
		account, err := helper.ReadFileJson(filepath.Join(queryDataPath+"/"+botUsername, file.Name()))
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("Failed to read file %s: %v", file.Name(), err))
			continue
		}

		switch v := account.(type) {
		case []map[string]string:
			// Jika data adalah array of maps
			for _, acc := range v {
				// Konversi map[string]string ke queryData
				for key, value := range acc {
					if key == "QueryData" {
						mergedData = append(mergedData, SessionStorage{
							QueryData: value,
						})
					} else if key == "Username" {
						mergedData = append(mergedData, SessionStorage{
							Username: value,
						})
					}
				}
			}
		case map[string]string:
			// Jika data adalah single map
			for key, value := range v {
				if key == "QueryData" {
					mergedData = append(mergedData, SessionStorage{
						QueryData: value,
					})
				} else if key == "Username" {
					mergedData = append(mergedData, SessionStorage{
						Username: value,
					})
				}
			}
		default:
			helper.PrettyLog("error", "Unknown data type")
		}
	}

	// Check Folder Query Data
	if !helper.CheckFileOrFolder(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, mergePath)) {
		os.Mkdir(fmt.Sprintf("%v/%v/%v", queryDataPath, botUsername, mergePath), 0755)
	}

	// Save Query Data To Json
	helper.SaveFileJson(fmt.Sprintf("%v/%v/%v/%v.json", queryDataPath, botUsername, mergePath, mergeFileName), mergedData)

	// Save Query Data To Txt
	for _, value := range mergedData {
		err := helper.SaveFileTxt(fmt.Sprintf("%v/%v/%v/%v.txt", queryDataPath, botUsername, mergePath, mergeFileName), value.QueryData)
		if err != nil {
			fmt.Printf("Error saving file: %v\n", err)
		}
	}

	helper.PrettyLog("success", fmt.Sprintf("Query Data Berhasil Di Simpan Di %v/%v/%v/%v.json", queryDataPath, botUsername, mergePath, mergeFileName))
}
