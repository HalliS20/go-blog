package controller

import (
	"fmt"
	"go-blog/internal/models"
	"go-blog/internal/service"
	"html/template"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	posts       []models.BlogPost
	postsMutex  sync.RWMutex
	faviconData string
	cssMain     template.CSS
	cssPost     template.CSS
	cssPostable template.CSS
	jsMain      template.JS
)

type BlogPost = models.BlogPost

func Init() {
	service.InitDatabase()
	getStaticFiles()
	setupPostListener()
}

func Shutdown() {
	service.CloseDatabase()
}

func setupPostListener() {
	go func() {
		ch := service.SetupListener()
		for notification := range ch {
			log.Println("Received notification:", notification)
			updatePostsFromDB()
		}
	}()
}

func updatePostsFromDB() {
	newPosts := service.GetBlogPosts()
	postsMutex.Lock()
	posts = newPosts
	postsMutex.Unlock()
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
	postsMutex.RLock()
	defer postsMutex.RUnlock()
	return gin.H{
		"cssContent":  cssMain,
		"jsFile":      jsMain,
		"faviconData": faviconData,
		"posts":       posts,
	}
}

func GetPostableData() gin.H {
	postsMutex.RLock()
	defer postsMutex.RUnlock()
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
}

func GetPost(id int) *BlogPost {
	postsMutex.RLock()
	defer postsMutex.RUnlock()

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
