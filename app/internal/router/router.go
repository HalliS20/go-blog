package router

import (
	ctrl "go-blog/internal/controller"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	ctrl.Init()
	setRoutes(router)

	//======= shuts down the database connection when the server is stopped
}

func setRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) { showPosts(c) })            // home page (with posts)
	router.GET("/posts", func(c *gin.Context) { showPosts(c) })       // show all posts
	router.GET("/postable", func(c *gin.Context) { showPostable(c) }) // postable site
	router.GET("/posts/:id", func(c *gin.Context) { showPost(c) })    // show a single post
	router.POST("/posts", func(c *gin.Context) { sendPost(c) })       // send a post
	router.GET("/sitemap.xml", func(c *gin.Context) { c.File("public/sitemap.xml") })
	router.GET("/public/*filepath", func(c *gin.Context) {
		c.File("public/" + c.Param("filepath"))
	}) // serve static files
}

func Shutdown() {
	ctrl.Shutdown()
}
