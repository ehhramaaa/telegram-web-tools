package core

import (
	"fmt"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
)

// Search And Select
func searchBot(page *rod.Page, text string) {
	inputText(page, text, "div.input-search > input")

	time.Sleep(3 * time.Second)

	isAlreadyChat := checkElement(page, "div.search-super > div > div.search-super-tab-container > div > div:nth-child(1) > ul > a")

	if isAlreadyChat {
		clickElement(page, "div.search-super > div > div.search-super-tab-container > div > div:nth-child(1) > ul > a")
	} else {
		clickElement(page, "#search-container > div.scrollable.scrollable-y > div.search-super > div > div.search-super-tab-container.search-super-container-chats.tabs-tab.active > div > div.search-group.search-group-contacts.is-short > ul > a:nth-child(1)")
	}

	helper.PrettyLog("success", fmt.Sprintf("| %s | Search User %s Successfully", phoneNumber, text))
}

func sendMessage(page *rod.Page, text string, isBot bool) {
	if isBot {
		// Click Start Button If Not Started Yet
		clickElement(page, "div.chat-input-control.chat-input-wrapper > button:nth-child(1)")
	}

	// Input Message
	inputText(page, text, "div.input-message-container > div:nth-child(1)")

	// Click Send Message
	clickElement(page, "div.btn-send-container > button")

	helper.PrettyLog("success", fmt.Sprintf("| %s | Send Message %s Successfully", phoneNumber, text))
}

func getLastChat(page *rod.Page) string {
	return getText(page, "div.bubbles-group.bubbles-group-last > div > div > div > div > span.translatable-message")
}

func popupLaunchBot(page *rod.Page) {
	// Click Popup Launch If Found
	isPopupLaunch := checkElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")

	if isPopupLaunch {
		clickElement(page, "body > div.popup.popup-peer.popup-confirmation.active > div > div.popup-buttons > button:nth-child(1)")
	}
}
