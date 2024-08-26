package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"go-blog/service"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/gzip"
	"github.com/russross/blackfriday/v2"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	posts       []service.BlogPost
	e           *gin.Engine
	faviconData string
)

func main() {
	// Initialize the database
	// loadEnv()
	service.InitDatabase()
	defer func(Db *sql.DB) {
		err := Db.Close()
		if err != nil {
			panic(err)
		}
	}(service.Db)

	initializeServer()

	posts = service.GetBlogPosts()

	setRoutes()

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
	// e.Static("/styling", "./public/styling")

	faviconBytes, err := os.ReadFile("./public/static/favicon.ico")
	if err != nil {
		log.Fatal("Error reading favicon: ", err)
	}
	faviconData = base64.StdEncoding.EncodeToString(faviconBytes)
}

func setRoutes() {
	cssMain, cssPost, cssPostable, err := getCSS()
	if err != nil {
		log.Fatal("Error reading CSS files: ", err)
	}

	jsMain, err := getJS()
	if err != nil {
		log.Fatal("Error reading JS file: ", err)
	}

	e.GET("/", func(c *gin.Context) { showPosts(c, cssMain, jsMain) })                // home page (with posts)
	e.GET("/posts", func(c *gin.Context) { showPosts(c, cssMain, jsMain) })           // show all posts
	e.GET("/postable", func(c *gin.Context) { showPostable(c, cssPostable, jsMain) }) // postable site
	e.GET("/posts/:id", func(c *gin.Context) { showPost(c, cssPost, jsMain) })        // show a single post
	e.POST("/posts", func(c *gin.Context) { sendPost(c) })                            // send a post
	e.GET("/public/*filepath", func(c *gin.Context) {
		c.File("public/" + c.Param("filepath"))
	}) // serve static files
}

func showPosts(c *gin.Context, cssMain template.CSS, jsMain template.JS) {
	c.Header("Cache-Control", "no-cache")
	c.HTML(200, "index.html", gin.H{
		"posts":       posts,
		"cssContent":  cssMain,
		"jsFile":      jsMain,
		"faviconData": faviconData,
	})
}

func showPostable(c *gin.Context, cssPostable template.CSS, jsMain template.JS) {
	c.Header("Cache-Control", "no-cache")
	c.HTML(200, "postable.html", gin.H{
		"posts":       posts,
		"cssContent":  cssPostable,
		"jsFile":      jsMain,
		"faviconData": faviconData,
	})
}

func showPost(c *gin.Context, cssPost template.CSS, jsMain template.JS) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	//==== Find the correct post
	var post *service.BlogPost
	for _, p := range posts {
		if p.ID == id {
			post = &p
			break
		}
	}

	if post == nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	safeBody := template.HTML(post.Body)
	c.HTML(200, "post.html", gin.H{
		"post":         post,
		"title":        post.Title,
		"description":  post.Description,
		"body":         safeBody,
		"canonicalURL": fmt.Sprintf("https://localhost:8080/posts/%d", post.ID),
		"cssContent":   cssPost,
		"jsFile":       jsMain,
		"faviconData":  faviconData,
	})
}

func sendPost(c *gin.Context) {
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

	bodyMarkdown := blackfriday.Run([]byte(body))

	post := service.BlogPost{
		Title:       title,
		Body:        string(bodyMarkdown),
		Description: description,
	}
	service.CreateBlogPost(post)
	c.JSON(200, gin.H{"status": "posted"})
	posts = service.GetBlogPosts()
}

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

func getJS() (template.JS, error) {
	js, err := readFile("./public/scripts/main.js")
	if err != nil {
		return "", err
	}
	return template.JS(js), nil
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
