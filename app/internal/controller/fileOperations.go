package controller

import (
	"log"
	"os"
)

func CheckPassword(password string) bool {
	if password != os.Getenv("PASSWORD") {
		log.Println("Password mismatch")
		log.Println(os.Getenv("PASSWORD"))
		return false
	}
	return true
}
