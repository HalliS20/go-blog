package main

import (
	"database/sql"
	"fmt"
	"go-blog/service"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getCSS() (template.CSS, template.CSS, template.CSS, error) {
	cssMain, err := readFile("./public/styling/main.css")
	if err != nil {
		return "", "", "", err
	}
	cssPost, err := readFile("./public/styling/post.css")
	if err != nil {
		return "", "", "", err
	}
	cssPostable, err := readFile("./public/styling/postable.css")
	if err != nil {
		return "", "", "", err
	}

	safeCssMain := template.CSS(cssMain)
	safeCssPost := template.CSS(cssPost)
	safeCssPostable := template.CSS(cssPostable)
	return safeCssMain, safeCssPost, safeCssPostable, nil
}

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

	cssMain, cssPost, cssPostable, err := getCSS()
	if err != nil {
		log.Fatal("Error reading CSS files: ", err)
	}

	e.StaticFile("/favicon.ico", "./public/static/favicon.ico")
	e.Static("/styling", "./public/styling")

	newPost := 0

	// Define the routes
	//=========== GET / - Display the list of blog posts
	e.GET("/", func(c *gin.Context) {
		if newPost == 1 {
			posts = service.GetBlogPosts()
			newPost = 0
		} 
		c.Header("Cache-Control", "no-cache")
		c.HTML(200, "index.html", gin.H{
			"posts":      posts,
			"cssContent": cssMain,
		})
	})

	//=========== GET /postable - can post a new blog post
	e.GET("/postable", func(c *gin.Context) {
		if newPost == 1 {
			posts = service.GetBlogPosts()
			newPost = 0
		} 
		c.Header("Cache-Control", "no-cache")
		c.HTML(200, "postable.html", gin.H{
			"posts":      posts,
			"cssContent": cssPostable,
		})
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
		post := posts[len(posts)-id]
		c.HTML(200, "post.html", gin.H{
			"post":         post,
			"title":        post.Title,
			"description":  post.Description,
			"canonicalURL": fmt.Sprintf("https://localhost:8080/posts/%d", post.ID),
			"cssContent":   cssPost,
		})
	})

	//========== POST /post - Create a new blog post
	e.POST("/post", func(c *gin.Context) {
		title := c.PostForm("title")
		body := c.PostForm("body")
		description := c.PostForm("description")
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
			Description: description,
		}
		service.CreateBlogPost(post)
		c.JSON(200, gin.H{"status": "posted"})
		newPost = 1
	})

	//======== Run the server
	err = e.Run(":8080")
	if err != nil {
		return
	}
}
