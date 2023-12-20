package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tipbk/sneakfeed-service/config"
	"github.com/tipbk/sneakfeed-service/dto"
	"github.com/tipbk/sneakfeed-service/service"
	"github.com/tipbk/sneakfeed-service/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	RefreshToken(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	GetUserByOthers(c *gin.Context)
	ToggleFollowUser(c *gin.Context)
}

type userHandler struct {
	envConfig            *config.EnvConfig
	userService          service.UserService
	imageUploaderService service.ImageUploaderService
}

func NewUserHandler(envConfig *config.EnvConfig, userService service.UserService, imageUploaderService service.ImageUploaderService) UserHandler {
	return &userHandler{
		envConfig:            envConfig,
		userService:          userService,
		imageUploaderService: imageUploaderService,
	}
}

func (h *userHandler) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	if registerRequest.Username == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("username cannot be empty"))
		return
	}
	if registerRequest.Password == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("password cannot be empty"))
		return
	}
	if registerRequest.Email == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("email cannot be empty"))
		return
	}
	_, err := h.userService.CreateUser(strings.ToLower(registerRequest.Username), registerRequest.Password, strings.ToLower(registerRequest.Email))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse("register done"))
}

func (h *userHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	if loginRequest.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be empty"})
		return
	}
	if loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be empty"})
		return
	}
	user, err := h.userService.LoginUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse("Username or password is incorrect or username does not exist."))
		return
	}
	accessToken, err := util.GenerateAccessToken(h.envConfig.AccessTokenSecret, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	refreshToken, err := util.GenerateRefreshToken(h.envConfig.RefreshTokenSecret, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}

func (h *userHandler) GetProfile(c *gin.Context) {
	currentUser, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(currentUser))
}

func (h *userHandler) GetUserByOthers(c *gin.Context) {
	currentUser, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("username cannot be empty"))
		return
	}
	userView, err := h.userService.FindUserViewByOthers(currentUser.ID.Hex(), username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, util.GenerateFailedResponse("user doesn't not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(dto.GetUserByUsernameResponse{
		IsYourUser:       userView.ID.Hex() == currentUser.ID.Hex(),
		UserViewByOthers: userView,
	}))
}

func (h *userHandler) RefreshToken(c *gin.Context) {
	var refreshTokenRequest dto.RefreshtokenRequest
	if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	if refreshTokenRequest.RefreshToken == "" {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse("refresh token cannot be empty"))
		return
	}

	jwt, err := util.ValidateRefreshToken(h.envConfig.RefreshTokenSecret, refreshTokenRequest.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}

	userID := jwt["userID"].(string)
	user, err := h.userService.FindUserWithUserID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}

	accessToken, err := util.GenerateAccessToken(h.envConfig.AccessTokenSecret, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	refreshToken, err := util.GenerateRefreshToken(h.envConfig.RefreshTokenSecret, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}

func (h *userHandler) UpdateUserProfile(c *gin.Context) {
	currentUser, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	var request dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}

	imageUrl := ""

	if request.ImageBase64 != "" {
		uploadResponse, err := h.imageUploaderService.UploadImage(request.ImageBase64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
			return
		}
		imageUrl = uploadResponse.Data.Url
	}

	if imageUrl != "" {
		currentUser.ProfileImage = imageUrl
	}

	if request.DisplayName != "" {
		currentUser.DisplayName = request.DisplayName
	}

	err = h.userService.UpdateProfile(currentUser.ID.Hex(), currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse("updated"))
}

func (h *userHandler) ToggleFollowUser(c *gin.Context) {
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}

	var request dto.ToggleFollowUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}

	if request.FollowUserID == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("followUserID cannot be empty"))
		return
	}

	if request.FollowUserID == user.ID.Hex() {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("you cannot follow yourself"))
		return
	}

	isFollowed, err := h.userService.ToggleFollowOnUser(user.ID.Hex(), request.FollowUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(dto.ToggleFollowUserResponse{IsFollowed: isFollowed}))
}
