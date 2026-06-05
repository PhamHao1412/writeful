package route

import (
	apiv1 "auth-service/internal/handler/rest/v1"

	"github.com/gin-gonic/gin"
)

func V1Router(r *gin.Engine, healthHandler *apiv1.HealthHandler, authHandler *apiv1.AuthHandler, followerHandler *apiv1.FollowerHandler) {

	// Root health check
	r.GET("/health", healthHandler.Health)

	v1 := r.Group("/api/v1")
	{

		// Auth endpoints (public)
		userRoute := v1.Group("/user")
		{
			userRoute.GET("/list", authHandler.List)
			userRoute.GET("/:username", authHandler.Detail)
			userRoute.GET("/profile", authHandler.Profile)
			userRoute.POST("/signup", authHandler.SignUp)
			userRoute.POST("/login", authHandler.Login)
			userRoute.POST("/refresh", authHandler.RefreshToken)
			userRoute.POST("/logout", authHandler.Logout)
			userRoute.POST("/update-info", authHandler.UpdateInfo)
			userRoute.POST("/active", authHandler.UpdateActiveStatus)
		}

		followerRoute := v1.Group("/user")
		{
			followerRoute.POST("/follow/:id", followerHandler.Follow)
			followerRoute.DELETE("/unfollow/:id", followerHandler.Unfollow)
			followerRoute.GET("/followers/:id", followerHandler.GetFollowers)
			followerRoute.GET("/following/:id", followerHandler.GetFollowing)
		}
	}
}
