package main

import (
	"library/task/controllers"
	"library/task/middleware"
	"library/task/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	models.ConnectDatabase()
}

func main() {

	r := gin.Default()

	// cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "token", "Content-Type", "Accept", "user"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/qrcodes", "./qrcodes")

	// Authentication routes
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	// r.GET("/allbooks/:libid", controllers.FetchAllBooks)

	// Authorized routes
	authRoutes := r.Group("/")
	authRoutes.Use(middleware.Authenticate())

	ownerRoutes := r.Group("/owner")
	ownerRoutes.Use(middleware.Authenticate())
	ownerRoutes.Use(middleware.Authorize("owner"))
	ownerRoutes.POST("/onboardAdmin", controllers.OnboardAdmin)

	// Admin routes
	adminRoutes := authRoutes.Group("/admin")
	adminRoutes.Use(middleware.Authenticate())
	adminRoutes.Use(middleware.Authorize("admin"))
	{
		adminRoutes.POST("/onboardReader", controllers.OnboardReader)
		adminRoutes.POST("/createlibrary", controllers.CreateLibrary)
		adminRoutes.POST("/createbookinventory", controllers.CreateBookInventory)
		adminRoutes.DELETE("/deletebook/:id", controllers.DeleteBook)
		adminRoutes.PATCH("/updatebook", controllers.UpdateBook)
		adminRoutes.POST("/approvedisapprove", controllers.ApproveDisapprove)
		adminRoutes.POST("/disapprove", controllers.Disapprove)
		adminRoutes.GET("/requests/:libid", controllers.Requests)
		adminRoutes.GET("/allbooks/:libid", controllers.FetchAllBooks)
	}

	readerRoutes := authRoutes.Group("/reader")
	readerRoutes.Use(middleware.Authenticate())
	readerRoutes.Use(middleware.Authorize("reader"))
	{
		readerRoutes.POST("/createissue", controllers.CreateIssueRequests)
		readerRoutes.POST("/searchbook", controllers.SearchBookBy)
		readerRoutes.POST("/returnbook", controllers.ReturnRequests)
		readerRoutes.GET("/allbooks/:libid", controllers.FetchAllBooks)
		readerRoutes.GET("/issued/:user", controllers.Issued)
	}

	r.Run()
}
