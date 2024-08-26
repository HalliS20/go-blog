package controller

import (
	"encoding/base64"
	"html/template"
	"log"
	"os"
)

func getFaviconData() string {
	faviconBytes, err := os.ReadFile("./public/static/favicon.ico")
	if err != nil {
		log.Fatal("Error reading favicon: ", err)
	}
	faviconData = base64.StdEncoding.EncodeToString(faviconBytes)
	return faviconData
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getJS(filename string) template.JS {
	js, err := readFile("./public/scripts/" + filename)
	if err != nil {
		println("Failed to read JS file %s: %v", filename, err)
		return ""
	}
	return template.JS(js)
}

func getCSS(filename string) template.CSS {
	css, err := readFile("./public/styling/" + filename)
	if err != nil {
		println("Failed to read CSS file %s: %v", filename, err)
		return ""
	}
	return template.CSS(css)
}

func CheckPassword(password string) bool {
	if password != os.Getenv("PASSWORD") {
		log.Println("Password mismatch")
		log.Println(os.Getenv("PASSWORD"))
		return false
	}
	return true
}
