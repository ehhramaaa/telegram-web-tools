package core

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"github.com/gookit/config/v2"
)

func initializeBrowser() *rod.Browser {
	extensionPath, _ := filepath.Abs("./extension/mini-app-android-spoof")

	launchOptions := launcher.New().
		Set("load-extension", extensionPath).
		Headless(config.Bool("HEADLESS_MODE")).
		MustLaunch()

	browser := rod.New().ControlURL(launchOptions).MustConnect()

	return browser
}

func (c *Client) checkElement(page *rod.Page, selector string) bool {
	// Recovery from panic, in case of unexpected errors
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	// Custom sleep function (3 seconds sleep)
	sleep := func() utils.Sleeper {
		return func(context.Context) error {
			time.Sleep(5 * time.Second)
			return nil
		}
	}

	// Try to find the element up to 3 times
	for attempt := 1; attempt <= 3; attempt++ {
		// Try to find the element with a timeout and custom sleeper
		_, err := page.Timeout(10 * time.Second).Sleeper(sleep).Element(selector)

		if err == nil {
			// Element found, return true
			return true
		} else if errors.Is(err, &rod.ElementNotFoundError{}) {
			// If the element is not found and we reached the max attempt, log and return false
			if attempt == 3 {
				helper.PrettyLog("warning", fmt.Sprintf("| %s | Element %v not found after %d attempts", c.phoneNumber, selector, attempt))
				return false
			}
			// Sleep between attempts
			time.Sleep(2 * time.Second)
		} else {
			// If another error occurs, panic
			panic(err)
		}
	}

	return false
}

func (c *Client) navigate(page *rod.Page, url string) {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	page.Timeout(3 * time.Second).Navigate(url)
	page.MustWaitLoad()
	page.MustWaitRequestIdle()
}

func (c *Client) clickElement(page *rod.Page, selector string) {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	c.checkElement(page, selector)

	page.Timeout(3 * time.Second).MustElement(selector).MustWaitVisible()

	page.Timeout(3 * time.Second).MustElement(selector).MustClick()

	page.MustWaitRequestIdle()
}

func (c *Client) inputText(page *rod.Page, value string, selector string) {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	c.checkElement(page, selector)

	page.Timeout(3 * time.Second).MustElement(selector).MustClick().MustInput(value)
}

func (c *Client) getText(page *rod.Page, selector string) string {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	c.checkElement(page, selector)

	text := page.Timeout(10 * time.Second).MustElement(selector).MustText()

	return text
}

func (c *Client) removeTextFormInput(page *rod.Page, selector string) {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	c.checkElement(page, selector)

	// Dapatkan elemen berdasarkan selector, periksa apakah elemen ada
	element, err := page.Element(selector)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Element not found: %v", c.phoneNumber, err))
		return
	}
	if element == nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Element is nil for selector: %s", c.phoneNumber, selector))
		return
	}

	// Periksa apakah elemen adalah input atau div
	tagName, err := element.Eval(`() => this.tagName.toLowerCase()`)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get tag name: %v", c.phoneNumber, err))
		return
	}

	switch tagName.Value.String() {
	case "input":
		// Jika elemen adalah input, hapus teks
		page.MustElement(selector).MustSelectAllText().MustInput("")
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to clear input text: %v", c.phoneNumber, err))
		}
	case "div":
		// Jika elemen adalah div, hapus teks menggunakan JavaScript
		_, err = element.Eval(`() => { this.textContent = ""; }`)
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Failed To Remove Text From Input Field: %v", c.phoneNumber, err))
		} else {
			helper.PrettyLog("info", fmt.Sprintf("| %s | Remove Text From Input Field", c.phoneNumber))
		}
	default:
		helper.PrettyLog("info", fmt.Sprintf("| %s | The element is not an input or div, skipping", c.phoneNumber))
	}
}
