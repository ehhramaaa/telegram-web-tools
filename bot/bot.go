package bot

import (
	"telegram-web/core"
	"telegram-web/helper"
	"time"

	"github.com/go-rod/rod"
)

// Search And Select
func SearchBot(page *rod.Page, text string) {
	core.InputText(page, text, "div.input-search > input")

	time.Sleep(3 * time.Second)

	isAlreadyChat := core.CheckElement(page, "div.search-super > div > div.search-super-tab-container > div > div:nth-child(1) > ul > a")

	if isAlreadyChat {
		core.ClickElement(page, "div.search-super > div > div.search-super-tab-container > div > div:nth-child(1) > ul > a")
	} else {
		core.ClickElement(page, "#search-container > div.scrollable.scrollable-y > div.search-super > div > div.search-super-tab-container.search-super-container-chats.tabs-tab.active > div > div.search-group.search-group-contacts.is-short > ul > a:nth-child(1)")
	}

	helper.PrettyLog("success", "Search Bot")
}

func SendMessage(page *rod.Page, text string) {
	// Click Start Button If Not Started Yet
	core.ClickElement(page, "div.chat-input-control.chat-input-wrapper > button:nth-child(1)")

	// Input Message
	core.InputText(page, text, "div.input-message-container > div:nth-child(1)")

	// Click Send Message
	core.ClickElement(page, "div.btn-send-container > button")

	helper.PrettyLog("success", "Send Message")
}

func GetLastChat(page *rod.Page) string {
	return core.GetText(page, "div.bubbles-group.bubbles-group-last > div > div > div > div > span.translatable-message")
}
