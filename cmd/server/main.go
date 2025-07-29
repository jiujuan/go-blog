package main

import (
	"log"

	"go-blog/internal/database"
	"go-blog/internal/handlers"
	"go-blog/internal/middleware"
	"go-blog/internal/repositories"
	"go-blog/internal/services"
	"go-blog/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.ConnectWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	articleRepo := repositories.NewArticleRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	commentRepo := repositories.NewCommentRepository(db)
	likeRepo := repositories.NewLikeRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret)
	userService := services.NewUserService(userRepo)
	userService.SetArticleRepository(articleRepo) // Inject article repository for user articles
	articleService := services.NewArticleService(articleRepo, userRepo, categoryRepo, tagRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	tagService := services.NewTagService(tagRepo)
	commentService := services.NewCommentService(commentRepo, articleRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	articleHandler := handlers.NewArticleHandler(articleService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	tagHandler := handlers.NewTagHandler(tagService)
	commentHandler := handlers.NewCommentHandler(commentService)

	// Setup router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Setup routes
	setupRoutes(router, authHandler, userHandler, articleHandler, categoryHandler, tagHandler, commentHandler, authService)

	// Start server
	log.Printf("Server starting on %s", cfg.GetServerAddress())
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	articleHandler *handlers.ArticleHandler,
	categoryHandler *handlers.CategoryHandler,
	tagHandler *handlers.TagHandler,
	commentHandler *handlers.CommentHandler,
	authService *services.AuthService,
) {
	api := router.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/me", middleware.Auth(authService), authHandler.Me)
	}

	// User routes
	users := api.Group("/users")
	{
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", middleware.Auth(authService), userHandler.Update)
		users.GET("/:id/articles", userHandler.GetUserArticles)
	}

	// Article routes
	articles := api.Group("/articles")
	{
		articles.GET("", articleHandler.List)
		articles.POST("", middleware.Auth(authService), articleHandler.Create)
		articles.GET("/search", articleHandler.Search)
		articles.GET("/:slug", articleHandler.GetBySlug)
		articles.PUT("/:id", middleware.Auth(authService), articleHandler.Update)
		articles.DELETE("/:id", middleware.Auth(authService), articleHandler.Delete)
		articles.POST("/:id/like", middleware.Auth(authService), articleHandler.ToggleLike)
	}

	// Category routes
	categories := api.Group("/categories")
	{
		categories.GET("", categoryHandler.List)
		categories.POST("", middleware.Auth(authService), categoryHandler.Create)
		categories.GET("/:slug", categoryHandler.GetBySlug)
		categories.PUT("/:id", middleware.Auth(authService), categoryHandler.Update)
		categories.DELETE("/:id", middleware.Auth(authService), categoryHandler.Delete)
		categories.GET("/:slug/articles", categoryHandler.GetCategoryArticles)
	}

	// Tag routes
	tags := api.Group("/tags")
	{
		tags.GET("", tagHandler.List)
		tags.POST("", middleware.Auth(authService), tagHandler.Create)
		tags.GET("/:slug", tagHandler.GetBySlug)
		tags.GET("/:slug/articles", tagHandler.GetTagArticles)
	}

	// Comment routes
	api.GET("/articles/:id/comments", commentHandler.GetByArticle)
	api.POST("/articles/:id/comments", middleware.Auth(authService), commentHandler.Create)
	api.PUT("/comments/:id", middleware.Auth(authService), commentHandler.Update)
	api.DELETE("/comments/:id", middleware.Auth(authService), commentHandler.Delete)

	// Archive routes
	archive := api.Group("/archive")
	{
		archive.GET("", articleHandler.GetArchive)
		archive.GET("/:year/:month", articleHandler.GetArchiveByMonth)
	}
}