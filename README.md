[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/demo.png)

# Telegram Web Tools

#### 🔥🔥 Os Tested: Windows 11

#### 🔥🔥 Using Rod Library, Similarly With Puppeteer

#### 🔥🔥 More Secure Than Use Session

## Recommendation before use

#### Go Version >= 1.23

## Features

|         Feature         | Supported |
| :---------------------: | :-------: |
|    Get Local Storage    |    ✅     |
|   Get Detail Account    |    ✅     |
|      Set Username       |    ✅     |
| Start Bot With Auto Ref |    ✅     |
|     Get Query Data      |    ✅     |
|    Merge Query Data     |    ✅     |
|     Set First Name      |    ⏳     |
|      Set Last Name      |    ⏳     |
|     Multithreading      |    ✅     |

## [Settings](https://github.com/ehhramaaa/telegram-web-tools/blob/main/config.yml)

|                       Settings                        |                                    Description                                    |
| :---------------------------------------------------: | :-------------------------------------------------------------------------------: |
|                   **BOT_USERNAME**                    |                   For get query data & start bot with auto ref                    |
|                    **MAX_THREAD**                     |                     Max client run parallel at the same time                      |
|                     **HEADLESS**                      |                     false = browser open, true = browser hide                     |
|             **GET_LOCAL_STORAGE.COUNTRY**             |             Country name of the number you want to get local storage              |
|            **GET_LOCAL_STORAGE.PASSWORD**             | Password for the number you use, or you can put "" for input password in terminal |
|          **START_BOT_WITH_AUTO_REF.REF_URL**          |                         Ref bot url of your main account                          |
| **START_BOT_WITH_AUTO_REF.FIRST_LAUNCH_BOT_SELECTOR** |                   All clickable selector when first launch bot                    |

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) **Go Version Tested 1.23.1**
- Remove .example From Files With Have That Name And Insert With Your Data
- If You Already Have Local Storage File With JSON Extension, You Can Put It At output/local-storage

## Installation

You can download the [**repository**](https://github.com/ehhramaaa/telegram-web-tools.git) by cloning it to your system and installing the necessary dependencies:

```shell
git clone https://github.com/ehhramaaa/telegram-web-tools.git
cd telegram-web-tools
go run .
```

## Usage

```shell
go run .
```

Or

```shell
go run main.go
```

## Or you can do build application by typing:

Windows:

```shell
go build -o telegramWebTools.exe
```

Linux:

```shell
go build -o telegramWebTools
chmod +x telegramWebTools
./telegramWebTools
```

# Get Local Storage Session Demo

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/local_storage.png)

# Get Detail Account Demo

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/get_detail.png)

# Set Username Demo

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/set_username.png)

# Start Bot With Auto Ref Demo

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/start_auto_ref.png)

# Get Query Data Tools Demo

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/query.png)
