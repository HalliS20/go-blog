package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/russross/blackfriday/v2"
	"go-blog/service"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	posts []service.BlogPost
	e     *gin.Engine
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

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	//Initialize the database
	loadEnv()
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
	e.StaticFile("/favicon.ico", "./public/static/favicon.ico")
	e.Static("/styling", "./public/styling")

}

func setRoutes() {
	cssMain, cssPost, cssPostable, err := getCSS()
	if err != nil {
		log.Fatal("Error reading CSS files: ", err)
	}

	e.GET("/", func(c *gin.Context) { showPosts(c, cssMain) })                // show all posts
	e.GET("/postable", func(c *gin.Context) { showPostable(c, cssPostable) }) // postable site
	e.GET("/posts/:id", func(c *gin.Context) { showPost(c, cssPost) })        // show a single post
	e.POST("/post", func(c *gin.Context) { sendPost(c) })                     // send a post
	e.GET("/public/*filepath", func(c *gin.Context) {
		c.File("public/" + c.Param("filepath"))
	}) // serve static files
}

func showPosts(c *gin.Context, cssMain template.CSS) {
	c.Header("Cache-Control", "no-cache")
	c.HTML(200, "index.html", gin.H{
		"posts":      posts,
		"cssContent": cssMain,
	})
}

func showPostable(c *gin.Context, cssPostable template.CSS) {
	c.Header("Cache-Control", "no-cache")
	c.HTML(200, "postable.html", gin.H{
		"posts":      posts,
		"cssContent": cssPostable,
	})
}

func showPost(c *gin.Context, cssPost template.CSS) {
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
