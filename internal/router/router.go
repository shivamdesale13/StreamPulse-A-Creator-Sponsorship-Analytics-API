package router

import (
	"github.com/gin-gonic/gin"
	"github.com/streampulse/api/internal/auth"
	"github.com/streampulse/api/internal/handlers"
)

func New(
	jwt *auth.JWTManager,
	authH *handlers.AuthHandler,
	creatorH *handlers.CreatorHandler,
	sponsorH *handlers.SponsorHandler,
	dealH *handlers.DealHandler,
	analyticsH *handlers.AnalyticsHandler,
) *gin.Engine {
	r := gin.Default()
	authMW := auth.Middleware(jwt)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "streampulse"})
	})

	public := r.Group("/auth")
	{
		public.POST("/register", authH.Register)
		public.POST("/login", authH.Login)
	}

	api := r.Group("/api/v1", authMW)
	{
		creators := api.Group("/creators")
		{
			creators.GET("", creatorH.List)
			creators.POST("", creatorH.Create)
			creators.GET("/:id", creatorH.GetByID)
			creators.PUT("/:id", creatorH.Update)
			creators.DELETE("/:id", creatorH.Delete)
			creators.GET("/:id/analytics", creatorH.GetAnalytics)
		}

		sponsors := api.Group("/sponsors")
		{
			sponsors.GET("", sponsorH.List)
			sponsors.POST("", sponsorH.Create)
			sponsors.GET("/:id", sponsorH.GetByID)
			sponsors.PUT("/:id", sponsorH.Update)
			sponsors.DELETE("/:id", sponsorH.Delete)
		}

		deals := api.Group("/deals")
		{
			deals.GET("", dealH.List)
			deals.POST("", dealH.Create)
			deals.GET("/:id", dealH.GetByID)
			deals.PATCH("/:id/status", dealH.UpdateStatus)
			deals.DELETE("/:id", dealH.Delete)
			deals.GET("/:id/analytics", dealH.GetAnalytics)
			deals.POST("/:id/analytics", dealH.RecordAnalytics)
		}

		analytics := api.Group("/analytics")
		{
			analytics.GET("/summary", analyticsH.GetSummary)
		}
	}

	return r
}
