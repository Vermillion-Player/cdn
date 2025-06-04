package routes

import (
	"cdn/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/api")
	{
		userRoutes.POST("/register", func(c *gin.Context) {
			controllers.RegisterController(c)
		})
		userRoutes.POST("/login", func(c *gin.Context) {
			controllers.LoginController(c)
		})
		userRoutes.PATCH("/change_password", func(c *gin.Context) {
			controllers.ChangePasswordController(c)
		})
		userRoutes.DELETE("/delete_user", func(c *gin.Context) {
			controllers.DeleteUserController(c)
		})
	}
}
