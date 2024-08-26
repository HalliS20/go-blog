package controller

import (
	"fmt"
	"go-blog/models"
	"go-blog/service"
	"html/template"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	posts       []models.BlogPost
	faviconData string
	cssMain     template.CSS
	cssPost     template.CSS
	cssPostable template.CSS
	jsMain      template.JS
)

type BlogPost = models.BlogPost

func Init() {
	initializeDatabase()
	getStaticFiles()
}

func Shutdown() {
	service.CloseDatabase()
}

func initializeDatabase() {
	service.InitDatabase()
}

func getStaticFiles() {
	posts = service.GetBlogPosts()
	cssMain = getCSS("main.css")
	cssPost = getCSS("post.css")
	cssPostable = getCSS("postable.css")
	jsMain = getJS("main.js")
	faviconData = getFaviconData()
}

func GetMainData() gin.H {
	return gin.H{
		"cssContent":  cssMain,
		"jsFile":      jsMain,
		"faviconData": faviconData,
		"posts":       posts,
	}
}

func GetPostableData() gin.H {
	return gin.H{
		"cssContent":  cssPostable,
		"jsFile":      jsMain,
		"faviconData": faviconData,
		"posts":       posts,
	}
}

func GetPostData(post *BlogPost) gin.H {
	safeBody := template.HTML(post.Body)
	return gin.H{
		"cssContent":   cssPost,
		"jsFile":       jsMain,
		"faviconData":  faviconData,
		"canonicalURL": fmt.Sprintf("https://localhost:8080/posts/%d", post.ID),
		"post":         post,
		"title":        post.Title,
		"description":  post.Description,
		"body":         safeBody,
	}
}

func AddPost(post models.BlogPost) {
	service.CreateBlogPost(post)
	posts = service.GetBlogPosts()
}

func GetPost(id int) *BlogPost {
	//==== Find the correct post
	var post *BlogPost
	for _, p := range posts {
		if p.ID == id {
			post = &p
			break
		}
	}
	return post
}
