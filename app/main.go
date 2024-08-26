package main

import (
	"go-blog/internal/models"
	"go-blog/internal/router"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	// "log"
	//
	// "github.com/joho/godotenv"
)

type BlogPost = models.BlogPost

var e *gin.Engine

// func loadEnv() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}
// }

func main() {
	// Initialize the database
	// loadEnv()

	initializeServer()

	router.Init(e)

	//======== Run the server
	err := e.Run(":8080")
	if err != nil {
		return
	}
}

func initializeServer() {
	e = gin.Default()
	e.Use(gzip.Gzip(gzip.DefaultCompression)) // use gzip for text compression
	e.LoadHTMLGlob("templates/*")
}
