package main

import (
	"database/sql"
	"go-blog/service"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the database
	service.InitDatabase()
	defer func(Db *sql.DB) {
		err := Db.Close()
		if err != nil {
			panic(err)
		}
	}(service.Db)

	e := gin.Default()
	// use gzip for text compression
	e.Use(gzip.Gzip(gzip.DefaultCompression))

	e.LoadHTMLGlob("templates/*")

	posts := service.GetBlogPosts()

	// Define the routes
	//=========== GET / - Display the list of blog posts
	e.GET("/", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.HTML(200, "index.html", gin.H{"posts": posts})
	})

	//=========== GET /postable - can post a new blog post
	e.GET("/postable", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.HTML(200, "postable.html", gin.H{"posts": posts})
	})

	//=========== GET images and such
	e.GET("/public/*filepath", func(c *gin.Context) {
		c.File("public/" + c.Param("filepath"))
	})

	//=========== GET /post/:id - Display a single blog post
	e.GET("/posts/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid post ID"})
			return
		}
		post := posts[id-1]
		c.HTML(200, "post.html", gin.H{"post": post})
	})

	//========== POST /post - Create a new blog post
	e.POST("/post", func(c *gin.Context) {
		title := c.PostForm("title")
		body := c.PostForm("body")
		password := c.PostForm("password")

		if password != os.Getenv("PASSWORD") {
			log.Println("Password mismatch")
			log.Println(os.Getenv("PASSWORD"))
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		post := service.BlogPost{
			Title: title,
			Body:  body,
		}
		service.CreateBlogPost(post)
		c.JSON(200, gin.H{"status": "posted"})
	})

	//======== Run the server
	err := e.Run(":8080")
	if err != nil {
		return
	}
}
