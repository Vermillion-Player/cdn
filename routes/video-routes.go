package routes

import (
	"cdn/controllers"
	"cdn/middlewares"

	"github.com/gin-gonic/gin"
)

func VideoRoutes(r *gin.Engine) {
	protectedVideoRoutes := r.Group("/api/video")
	protectedVideoRoutes.Use(middlewares.JwtMiddleware())
	{
		protectedVideoRoutes.GET("/list", func(c *gin.Context) {
			controllers.ListVideoController(c)
		})
		protectedVideoRoutes.GET("/:id", func(c *gin.Context) {
			controllers.GetVideoController(c)
		})
		protectedVideoRoutes.POST("/new", func(c *gin.Context) {
			controllers.NewVideoController(c)
		})
		protectedVideoRoutes.PUT("/edit/:id", func(c *gin.Context) {
			controllers.EditVideoController(c)
		})
		protectedVideoRoutes.DELETE("/delete/:id", func(c *gin.Context) {
			controllers.DeleteVideoController(c)
		})
		protectedVideoRoutes.GET("/search/:q", func(c *gin.Context) {
			controllers.SearchVideoController(c)
		})
	}
}

