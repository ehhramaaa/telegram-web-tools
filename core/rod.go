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

// Todo

func (c *Client) checkElement(page *rod.Page, selector string) bool {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.phoneNumber, r))
		}
	}()

	sleep := func() utils.Sleeper {
		return func(context.Context) error {
			time.Sleep(3 * time.Second)
			return nil
		}
	}

	for attempt := 1; attempt <= 3; attempt++ {
		var err error
		page.Timeout(5 * time.Second).Sleeper(sleep).Element(selector)
		_, err = page.Timeout(5 * time.Second).Sleeper(rod.NotFoundSleeper).Element(selector)

		if err == nil {
			return true
		} else if errors.Is(err, &rod.ElementNotFoundError{}) {
			if attempt == 3 {
				helper.PrettyLog("warning", fmt.Sprintf("| %s | Check element %v not found after %d attempts", c.phoneNumber, selector, attempt))
				return false
			}
			time.Sleep(3 * time.Second)
		} else {
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
