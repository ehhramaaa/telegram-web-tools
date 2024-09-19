package core

import (
	"context"
	"errors"
	"fmt"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
)

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

	page.MustElement(selector).MustSelectAllText().MustInput("")
}
