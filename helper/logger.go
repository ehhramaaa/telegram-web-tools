package helper

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

func RecoverPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic : %v\n", r)
		}
	}()
}

func PrettyLog(level, message string) {
	level = strings.ToUpper(level)

	var levelColor *color.Color
	switch level {
	case "INFO":
		levelColor = color.New(color.FgWhite) // Blue for INFO
	case "ERROR":
		levelColor = color.New(color.FgRed) // Red for ERROR
	case "WARNING":
		levelColor = color.New(color.FgYellow) // Yellow for WARNING
	case "INPUT":
		levelColor = color.New(color.FgCyan) // Cyan for INPUT
	case "SUCCESS":
		levelColor = color.New(color.FgGreen) // Cyan for INPUT
	default:
		levelColor = color.New(color.FgWhite) // White for default
	}

	// Print the log message with color
	if level == "INPUT" {
		levelColor.Printf("[%s] ", level)
		fmt.Printf("%s", message)
	} else {
		levelColor.Printf("[%s] ", level)
		fmt.Printf("%s\n", message)

		if level == "ERROR" {
			fileName := fmt.Sprintf("./error/log_%s.txt", time.Now().Format("01-02-2006"))

			SaveFileTxt(fileName, fmt.Sprintf("| %s | %s", time.Now().Format("15:04:05"), message))
		}
	}
}

func PrintLogo() {
	levelColor := color.New(color.FgCyan)

	levelColor.Println(`
 /$$$$$$$$        /$$                 /$$      /$$           /$$             /$$$$$$$$                  /$$          
|__  $$__/       | $$                | $$  /$ | $$          | $$            |__  $$__/                 | $$          
   | $$  /$$$$$$ | $$  /$$$$$$       | $$ /$$$| $$  /$$$$$$ | $$$$$$$          | $$  /$$$$$$   /$$$$$$ | $$  /$$$$$$$
   | $$ /$$__  $$| $$ /$$__  $$      | $$/$$ $$ $$ /$$__  $$| $$__  $$         | $$ /$$__  $$ /$$__  $$| $$ /$$_____/
   | $$| $$$$$$$$| $$| $$$$$$$$      | $$$$_  $$$$| $$$$$$$$| $$  \ $$         | $$| $$  \ $$| $$  \ $$| $$|  $$$$$$ 
   | $$| $$_____/| $$| $$_____/      | $$$/ \  $$$| $$_____/| $$  | $$         | $$| $$  | $$| $$  | $$| $$ \____  $$
   | $$|  $$$$$$$| $$|  $$$$$$$      | $$/   \  $$|  $$$$$$$| $$$$$$$/         | $$|  $$$$$$/|  $$$$$$/| $$ /$$$$$$$/
   |__/ \_______/|__/ \_______/      |__/     \__/ \_______/|_______/          |__/ \______/  \______/ |__/|_______/ 
`)

	levelColor.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")
}

func ClearTerminal() {
	var clearCmd *exec.Cmd

	// Mengecek sistem operasi yang digunakan
	switch runtime.GOOS {
	case "linux", "darwin": // Untuk Linux dan macOS
		clearCmd = exec.Command("clear")
	case "windows": // Untuk Windows
		clearCmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	// Mengatur output ke terminal (Stdout)
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}
