package router

import (
	ctrl "go-blog/controller"
	"go-blog/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

func showPosts(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.HTML(200, "index.html", ctrl.GetMainData())
}

func showPostable(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	c.HTML(200, "postable.html", ctrl.GetPostableData())
}

func showPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	post := ctrl.GetPost(id)
	if post == nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.HTML(200, "post.html", ctrl.GetPostData(post))
}

func sendPost(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	description := c.PostForm("description")
	password := c.PostForm("password")

	if !ctrl.CheckPassword(password) {
		c.JSON(401, gin.H{"error": "Incorrect password"})
		return
	}

	bodyMarkdown := string(blackfriday.Run([]byte(body)))

	post := models.MakeBlogPost(title, description, bodyMarkdown)
	ctrl.AddPost(post)
	c.JSON(200, gin.H{"status": "posted"})
}
