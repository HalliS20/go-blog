package main

import (
	"go-blog/config"
	"go-blog/internal/router"
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	e.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/public/miniStyles/main.min.css" || c.Request.URL.Path == "/public/scripts/main.min.js" {
			c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		}
		c.Next()
	})
}
