package controller

import (
	"html/template"
	"log"
	"os"
)

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getCSS(filename string) template.CSS {
	css, err := readFile("./public/miniStyles/" + filename)
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
