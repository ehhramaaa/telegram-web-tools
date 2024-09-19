package helper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"strings"
)

func PointerString(s string) *string {
	return &s
}

func ReadFileDir(path string) []fs.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		PrettyLog("error", "Failed to read directory: %v")
	}

	return files
}

func InputTerminal(prompt string) string {
	PrettyLog("input", prompt)

	reader := bufio.NewReader(os.Stdin)

	value, _ := reader.ReadString('\n')

	return strings.TrimSpace(value)
}

func SaveFileJson(filePath string, data interface{}) error {
	// Membuka file untuk menulis data JSON
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Membuat encoder dengan SetEscapeHTML(false)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")  // Mengatur indentasi
	encoder.SetEscapeHTML(false) // Nonaktifkan HTML escaping

	// Encode data ke file
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	PrettyLog("success", fmt.Sprintf("Data berhasil disimpan ke %s", filePath))

	return nil
}

func ReadFileJson(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Coba unmarshal sebagai array of generic maps (map[string]interface{})
	var dataArray []map[string]interface{}
	if err := json.Unmarshal(byteValue, &dataArray); err == nil {
		PrettyLog("success", fmt.Sprintf("Data array berhasil dibaca dari %s", filePath))
		return dataArray, nil
	}

	// Jika gagal, coba unmarshal sebagai generic map
	var dataObject map[string]interface{}
	if err := json.Unmarshal(byteValue, &dataObject); err == nil {
		PrettyLog("success", fmt.Sprintf("Data object berhasil dibaca dari %s", filePath))
		return dataObject, nil
	}

	return nil, fmt.Errorf("failed to unmarshal JSON from file %s", filePath)
}

func SaveFileTxt(filePath string, data string) error {
	// Cek apakah file sudah ada
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// Jika file tidak ada, tulis data baru
		err = os.WriteFile(filePath, []byte(data+"\n"), 0644)
	} else {
		// Jika file sudah ada, tambahkan data ke akhir file
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(data + "\n")
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func CheckFileOrFolder(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func RandomNumber(min int, max int) int {
	return (rand.Intn(max-min) + min)
}

func GetTextAfterKey(urlData, key string) (string, error) {
	// Temukan lokasi key
	keyIndex := strings.Index(urlData, key)
	if keyIndex == -1 {
		return "", fmt.Errorf("key %s tidak ditemukan", key)
	}

	// Ambil substring setelah key
	startIndex := keyIndex + len(key)
	endIndex := strings.Index(urlData[startIndex:], "&")
	if endIndex == -1 {
		return urlData[startIndex:], nil
	}

	return urlData[startIndex : startIndex+endIndex], nil
}
