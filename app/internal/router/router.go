package router

import (
	ctrl "go-blog/internal/controller"

	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	ctrl.Init()
	setRoutes(e)

	//======= shuts down the database connection when the server is stopped
	defer ctrl.Shutdown()
}

func setRoutes(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) { showPosts(c) })            // home page (with posts)
	e.GET("/posts", func(c *gin.Context) { showPosts(c) })       // show all posts
	e.GET("/postable", func(c *gin.Context) { showPostable(c) }) // postable site
	e.GET("/posts/:id", func(c *gin.Context) { showPost(c) })    // show a single post
	e.POST("/posts", func(c *gin.Context) { sendPost(c) })       // send a post
	e.GET("/public/*filepath", func(c *gin.Context) {
		c.File("public/" + c.Param("filepath"))
	}) // serve static files
}

func Shutdown() {
	ctrl.Shutdown()
}
