package main

import (
	"github.com/gin-contrib/gzip"
	"go-blog/internal/router"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var e *gin.Engine

func main() {
	// Initialize the database

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
