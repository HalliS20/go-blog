package router

import (
	"go-blog/internal/controller"
	"go-blog/internal/domain/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

func (r *Router) showPosts(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=172800")
	c.HTML(200, "index.html", r.ctrl.GetMainData())
}

func (r *Router) showPostable(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=172800")
	c.HTML(200, "postable.html", r.ctrl.GetMainData())
}

func (r *Router) showPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	post := r.ctrl.GetPost(id)
	if post == nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=172800")
	c.HTML(200, "post.html", r.ctrl.GetPostData(post))
}

func (r *Router) sendPost(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	description := c.PostForm("description")
	password := c.PostForm("password")

	if !controller.CheckPassword(password) {
		c.JSON(401, gin.H{"error": "Incorrect password"})
		return
	}

	bodyMarkdown := string(blackfriday.Run([]byte(body)))

	post := models.MakeBlogPost(title, description, bodyMarkdown)
	r.ctrl.AddPost(post)
	c.JSON(200, gin.H{"status": "posted"})
}
