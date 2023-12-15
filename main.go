package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tipbk/sneakfeed-service/config"
	"github.com/tipbk/sneakfeed-service/handler"
	"github.com/tipbk/sneakfeed-service/middleware"
	"github.com/tipbk/sneakfeed-service/repository"
	"github.com/tipbk/sneakfeed-service/service"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	envConfig := config.GetEnvConfig()

	mongoClient, err := repository.CreateMongoConnection(envConfig)
	if err != nil {
		// take service down
		panic(err)
	}

	imageUploaderService := service.NewImageUploaderService()
	userRepository := repository.NewUserRepository(envConfig, mongoClient)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(envConfig, userService, imageUploaderService)
	contentRepository := repository.NewContentReepository(envConfig, mongoClient)
	contentService := service.NewContentService(contentRepository)
	contentHandler := handler.NewContentHandler(contentService, userService, imageUploaderService)
	authMiddleware := middleware.NewAuthMiddleware(envConfig, userService)

	r.GET("/ping")
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/refresh", userHandler.RefreshToken)

	authorized := r.Group("/")
	authorized.Use(authMiddleware.AuthAccessTokenMiddleware)
	{
		authorized.GET("/posts", contentHandler.GetPosts)
		authorized.GET("/posts/:postID", contentHandler.GetPostByID)
		authorized.GET("/posts/:postID/comments", contentHandler.GetCommentByPostID)
		authorized.POST("/posts", contentHandler.CreatePost)
		authorized.POST("/posts/:postID/comments", contentHandler.AddComment)
		authorized.GET("/profiles", userHandler.GetProfile)
		authorized.PATCH("/profiles/image", userHandler.UpdateProfilePicture)
		authorized.POST("/posts/:postID/like", contentHandler.ToggleLikePostByID)
	}

	r.Run()
}
