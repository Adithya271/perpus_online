package routes

import (
	"sistem_perpus/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Books routes
	r.GET("/books", controllers.GetBooks)
	r.GET("/books/:id", controllers.GetBookDetail)
	
	// Loans routes
	r.POST("/loans", controllers.BorrowBook)
	r.POST("/loans/return", controllers.ReturnBook)
	r.GET("/loans/user/:id", controllers.GetUserLoans)
	
	// Recommendations routes (
	r.GET("/recommendations/:id/by-search", controllers.GetRecommendationsBySearchHistory)
	
	// Search history routes
	r.POST("/search-history", controllers.SaveSearchHistory)
}