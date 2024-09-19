# Junior MiniApp Bypass

## Overview
Junior MiniApp Bypass is a browser extension designed to bypass Telegram mini apps. It modifies the behavior of Telegram Web to allow users to interact with mini apps that might otherwise be restricted.

Version: 1.2

## Features
- Bypasses restrictions on Telegram mini apps
- Modifies iframe sources to appear as iOS platform
- Temporarily changes the page title to indicate successful bypass

## Installation
1. Clone this repository or download the source code.
2. Open your browser's extension management page:
   - Chrome: chrome://extensions
   - Edge: edge://extensions
   - Firefox: about:addons
3. Enable "Developer mode"
4. Click "Load unpacked" and select the directory containing the extension files.

## Usage
Once installed, the extension will automatically work on the Telegram Web interface (https://web.telegram.org/*). No additional configuration is required.

## Files Description
- `telegram.js`: Contains the main logic for bypassing mini apps.
- `manifest.json`: Defines the extension's properties and permissions.
- `popup.html`: Provides a simple information popup for the extension.

## How It Works
The extension uses the following techniques:
1. Listens for new iframe elements using the Arrive library.
2. Modifies the `src` attribute of iframes, changing the platform identifier from web to iOS.
3. Temporarily changes the page title to indicate successful bypass.

## Permissions
This extension requires the following permissions:
- activeTab
- scripting
- webNavigation
- declarativeNetRequest
- declarativeNetRequestFeedback

## Community
For more information, guides, and updates, join our community:
- Telegram Channel: [Airdrop_DailyOfficial](https://t.me/Airdrop_DailyOfficial)
- YouTube Channel: [CryptoInsightNews](https://www.youtube.com/@CryptoInsightNews/)

## Disclaimer
This extension is for educational purposes only. Use it responsibly and in accordance with Telegram's terms of service.

## Creator
Created by [Junior](https://t.me/Airdrop_DailyOfficial)

## License
[Include license information here]
