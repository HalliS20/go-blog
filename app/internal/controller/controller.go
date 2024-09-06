package controller

import (
	"fmt"
	"go-blog/internal/models"
	"go-blog/internal/service"
	"html/template"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	posts       []models.BlogPost
	postsMutex  sync.RWMutex
	cssName     string
	cssData     template.CSS
	jsName      string
	faviconName string
)

type BlogPost = models.BlogPost

func Init() {
	service.InitDatabase()
	getStaticFiles()

}

func Shutdown() {
	service.CloseDatabase()
}

func updatePostsFromDB() {
	newPosts := service.GetBlogPosts()
	postsMutex.Lock()
	posts = newPosts
	postsMutex.Unlock()
}

func getStaticFiles() {
	cssName = "/public/miniStyles/total.min.css"
	cssData = getCSS("main.min.css")
	jsName = "/public/scripts/main.min.js"
	faviconName = "/public/static/favicon.ico"
	posts = service.GetBlogPosts()
}

func GetMainData() gin.H {
	postsMutex.RLock()
	defer postsMutex.RUnlock()
	return gin.H{
		"canonicalURL": fmt.Sprintf("https://localhost:8080/posts"),
		"cssName":      cssName,
		"cssData":      cssData,
		"jsName":       jsName,
		"faviconName":  faviconName,
		"posts":        posts,
	}
}

func GetPostableData() gin.H {
	postsMutex.RLock()
	defer postsMutex.RUnlock()
	return gin.H{
		"canonicalURL": fmt.Sprintf("https://localhost:8080/postable"),
		"cssName":      cssName,
		"jsName":       jsName,
		"posts":        posts,
		"faviconName":  faviconName,
	}
}

func GetPostData(post *BlogPost) gin.H {
	safeBody := template.HTML(post.Body)
	return gin.H{
		"canonicalURL": fmt.Sprintf("https://localhost:8080/posts/%d", post.ID),
		"cssName":      cssName,
		"jsName":       jsName,
		"post":         post,
		"title":        post.Title,
		"description":  post.Description,
		"body":         safeBody,
		"faviconName":  faviconName,
	}
}

func AddPost(post models.BlogPost) {
	service.CreateBlogPost(post)
	updatePostsFromDB()
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
