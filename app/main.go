package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-blog/config"
	"go-blog/internal/router"
	"log"
)

var e *gin.Engine

func main() {
	// Initialize the database

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	initializeServer()

	router.Init(e)

	//======== Shutdown the all layers when the server is closed
	defer router.Shutdown()
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
