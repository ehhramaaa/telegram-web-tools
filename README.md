[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/Screenshot_5.png)

# Telegram Web Tools

#### 🔥🔥 Os Tested: Windows 11

#### 🔥🔥 Using Rod Library, Similarly With Puppeteer

#### 🔥🔥 More Secure Than Use Session

## Recommendation before use

#### Go Version >= 1.23

## Features

|                  Feature                  | Supported |
| :---------------------------------------: | :-------: |
|          Auto Get Local Storage           |    ✅     |
|       Auto Start Bot With Auto Ref        |    ✅     |
|             Auto Set Username             |    ⏳     |
|            Auto Set Last Name             |    ⏳     |
|            Auto Set First Name            |    ⏳     |
| Multithreading (Except Get Local Storage) |    ✅     |

## [Settings](https://github.com/ehhramaaa/telegram-web-tools/blob/main/config.yml)

|                       Settings                        |                                            Description                                            |
| :---------------------------------------------------: | :-----------------------------------------------------------------------------------------------: |
|                   **BOT_USERNAME**                    |                           For Get Query Data / Start Bot With Auto Ref                            |
|                    **MAX_THREAD**                     |                                      Max Client Run Parallel                                      |
|             **GET_LOCAL_STORAGE.COUNTRY**             |                     country name of the number you want to take local storage                     |
|            **GET_LOCAL_STORAGE.PASSWORD**             | Password for number you use. if your number > 1 you must make all your account passwords the same |
|          **START_BOT_WITH_AUTO_REF.REF_URL**          |                                   Ref Url Of Your Main Account                                    |
| **START_BOT_WITH_AUTO_REF.FIRST_LAUNCH_BOT_SELECTOR** |                            All Selector Must Click In First Launch Bot                            |

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
![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/Screenshot_8.png)

# Start Bot With Auto Ref Demo
![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/Image.png)

# Get Query Data Tools Demo
![demo](https://raw.githubusercontent.com/ehhramaaa/telegram-web-tools/main/demo/query.png)