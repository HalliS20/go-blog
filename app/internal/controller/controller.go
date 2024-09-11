package controller

import (
	"fmt"
	"go-blog/internal/domain/models"
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

type (
	BlogPost    = models.BlogPost
	BlogService = service.BlogService
)

type Controller struct {
	blogService *BlogService
	posts       []models.BlogPost
	postsMutex  sync.RWMutex
	cssName     string
	cssData     template.CSS
	jsName      string
	faviconName string
}

func NewController(blogService *BlogService) *Controller {
	c := &Controller{blogService: blogService}
	c.Init()
	return c
}

func (c *Controller) Init() {
	c.getStaticFiles()
}

func (c *Controller) updatePostsFromDB() {
	newPosts := c.blogService.GetBlogPosts()
	postsMutex.Lock()
	posts = newPosts
	postsMutex.Unlock()
}

func (c *Controller) getStaticFiles() {
	cssName = "/public/miniStyles/total.min.css"
	cssData = getCSS("main.min.css")
	jsName = "/public/scripts/main.min.js"
	faviconName = "/public/static/favicon.ico"
	posts = c.blogService.GetBlogPosts()
}

func (c *Controller) GetMainData() gin.H {
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

func (c *Controller) GetPostData(post *BlogPost) gin.H {
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

func (c *Controller) AddPost(post models.BlogPost) {
	c.blogService.CreateBlogPost(post)
	c.updatePostsFromDB()
}

func (c *Controller) GetPost(id int) *BlogPost {
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
