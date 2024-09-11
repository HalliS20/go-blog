package main

import (
	"go-blog/internal/config"
	"go-blog/internal/controller"
	"go-blog/internal/database"
	"go-blog/internal/repositories"
	"go-blog/internal/router"
	"go-blog/internal/service"
	"log"
	"os"

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

	connStr := getConnStr()

	//====== Initialize the database
	postgresDB, err := database.NewPostgresConnection(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.CloseDatabase(postgresDB)

	//====== Initialize components
	postgresRepo := repositories.NewPostgresRepository(postgresDB)
	blogService := service.NewBlogService(postgresRepo)
	blogController := controller.NewController(blogService)
	blogRouter := router.NewRouter(blogController)

	initializeServer()

	blogRouter.Init(e)

	//======== Run the server
	err = e.Run(":8080")
	if err != nil {
		return
	}
}

func initializeServer() {
	e = gin.Default()
	e.Use(gzip.Gzip(gzip.DefaultCompression)) // use gzip for text compression
	e.LoadHTMLGlob("templates/*")
	e.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/public/miniStyles/total.min.css" || c.Request.URL.Path == "/public/scripts/main.min.js" {
			c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		}
		c.Next()
	})
}

func getConnStr() string {
	user := "Blog_owner"
	password := os.Getenv("DB_PASSWORD")
	hostName := "ep-late-sun-a5p8yfr7.us-east-2.aws.neon.tech"
	stringTail := " dbname=Blog port=5432 sslmode=require TimeZone=UTC"
	connStr := "host=" + hostName + " user=" + user + " password=" + password + stringTail
	return connStr
}
