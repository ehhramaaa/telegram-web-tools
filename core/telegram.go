package core

import (
	"fmt"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
)

// Search And Select In Telegram
func (c *Client) searchBot(page *rod.Page, text string) {
	c.inputText(page, text, "div.input-search > input")

	time.Sleep(3 * time.Second)

	isAlreadyChat := c.checkElement(page, "div.search-super > div > div.search-super-tab-container > div > div:nth-child(1) > ul > a")

	if isAlreadyChat {
		c.clickElement(page, "div.search-super > div > div.search-super-tab-container > div > div:nth-child(1) > ul > a")
	} else {
		c.clickElement(page, "#search-container > div.scrollable.scrollable-y > div.search-super > div > div.search-super-tab-container.search-super-container-chats.tabs-tab.active > div > div.search-group.search-group-contacts.is-short > ul > a:nth-child(1)")
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Search User %s Successfully", c.phoneNumber, text))
}

// Send Message Chat In Telegram
func (c *Client) sendMessage(page *rod.Page, text string, isBot bool) {
	if isBot {
		// Click Start Button If Not Started Yet
		c.clickElement(page, "div.chat-input-control.chat-input-wrapper > button:nth-child(1)")
	}

	// Input Message
	c.inputText(page, text, "div.input-message-container > div:nth-child(1)")

	// Click Send Message
	c.clickElement(page, "div.btn-send-container > button")

	helper.PrettyLog("success", fmt.Sprintf("| %s | Send Message %s Successfully", c.phoneNumber, text))
}

// Get Last Chat In Telegram
func (c *Client) getLastChat(page *rod.Page) string {
	return c.getText(page, "div.bubbles-group.bubbles-group-last > div > div > div > div > span.translatable-message")
}

// Click Popup Launch Now In Telegram Bot
func (c *Client) popupLaunchBot(page *rod.Page) {
	// Click Popup Launch If Found
	isPopupLaunch := c.checkElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")

	if isPopupLaunch {
		c.clickElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")
	}
}

// Goto Setting Page
func (c *Client) gotoSetting(page *rod.Page) bool {
	isClicked := false

	// Cari div yang berisi span dengan teks "Settings"
	for attempt := 1; attempt <= 3; attempt++ {
		divs, err := page.Timeout(10 * time.Second).Elements(`#column-left > div > div > div.sidebar-header.can-have-forum > div.sidebar-header__btn-container > button > div.btn-menu.bottom-right.has-footer.active.was-open > div`)
		if err != nil {
			helper.PrettyLog("waring", fmt.Sprintf("| %s | Element Settings Not Found, Try To Find Again After 3s...", c.phoneNumber))
			time.Sleep(3 * time.Second)
			continue
		}

		for _, div := range divs {
			span, err := div.ElementR(`span`, `Settings`)
			if err == nil && span != nil {
				div.MustClick()
				isClicked = true
				break
			}
		}
	}

	if !isClicked {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Element Settings Not Found After 3 Attempts", c.phoneNumber))
		return false
	}

	return true
}
