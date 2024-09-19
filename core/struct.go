package core

import "github.com/go-rod/rod"

type Client struct {
	phoneNumber string
	Browser     *rod.Browser
}
