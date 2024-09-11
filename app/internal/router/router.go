package router

import (
	ctrl "go-blog/internal/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ctrl *ctrl.Controller
}

// New Router takes in pointer to controller
// we could fix this by changing the controller to an interface
func NewRouter(ctrl *ctrl.Controller) *Router {
	r := &Router{ctrl: ctrl}
	return r
}

func (r *Router) Init(ginHandler *gin.Engine) {
	r.ctrl.Init()
	r.setRoutes(ginHandler)
}

func (r *Router) setRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) { r.showPosts(c) })            // home page (with posts)
	router.GET("/posts", func(c *gin.Context) { r.showPosts(c) })       // show all posts
	router.GET("/postable", func(c *gin.Context) { r.showPostable(c) }) // postable site
	router.GET("/posts/:id", func(c *gin.Context) { r.showPost(c) })    // show a single post
	router.POST("/posts", func(c *gin.Context) { r.sendPost(c) })       // send a post
	router.GET("/sitemap.xml", func(c *gin.Context) { c.File("public/sitemap.xml") })
	router.GET("/public/*filepath", func(c *gin.Context) {
		c.File("public/" + c.Param("filepath"))
	}) // serve static files
}
