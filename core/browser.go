package core

import (
	"context"
	"errors"
	"fmt"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/utils"
)

func CheckElement(page *rod.Page, element string) bool {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("error", fmt.Sprintf("Recovered from panic : %v", r))
		}
	}()

	sleep := func() utils.Sleeper {
		return func(context.Context) error {
			time.Sleep(5 * time.Second)
			return nil
		}
	}

	for attempt := 1; attempt <= 3; attempt++ {
		page.Timeout(5 * time.Second).Sleeper(sleep).Element(element)
		_, err := page.Timeout(5 * time.Second).Sleeper(rod.NotFoundSleeper).Element(element)

		if err == nil {
			return true
		} else if errors.Is(err, &rod.ElementNotFoundError{}) {
			if attempt == 3 {
				helper.PrettyLog("warning", fmt.Sprintf("Check element %v not found after %d attempts", element, attempt))
				return false
			}
			time.Sleep(3 * time.Second)
		} else {
			panic(err)
		}
	}

	return false
}

func Navigate(page *rod.Page, url string) {
	defer helper.RecoverPanic()

	page.Timeout(3 * time.Second).Navigate(url)
	page.MustWaitLoad()
	page.MustWaitRequestIdle()
}

func ClickElement(page *rod.Page, element string) {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("error", fmt.Sprintf("Recovered from panic : %v", r))
		}
	}()

	CheckElement(page, element)

	page.Timeout(3 * time.Second).MustElement(element).MustWaitVisible()

	page.Timeout(3 * time.Second).MustElement(element).MustClick()

	page.MustWaitRequestIdle()
}

func InputText(page *rod.Page, text string, element string) {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("error", fmt.Sprintf("Recovered from panic : %v", r))
		}
	}()

	CheckElement(page, element)

	page.Timeout(3 * time.Second).MustElement(element).MustClick().MustInput(text)
}

func GetText(page *rod.Page, element string) string {
	defer func() {
		if r := recover(); r != nil {
			helper.PrettyLog("error", fmt.Sprintf("Recovered from panic : %v", r))
		}
	}()

	CheckElement(page, element)

	text := page.Timeout(10 * time.Second).MustElement(element).MustText()

	return text
}

func DeleteText(page *rod.Page, element string) {
	ClickElement(page, element)

	page.Timeout(3 * time.Second).MustElement(element).MustKeyActions().Press(input.ControlLeft).Press(input.KeyA).Press(input.Backspace).Release(input.KeyA).Release(input.ControlLeft).Release(input.Backspace)
}
