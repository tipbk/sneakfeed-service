package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tipbk/sneakfeed-service/config"
	"github.com/tipbk/sneakfeed-service/dto"
	"github.com/tipbk/sneakfeed-service/model"
	"github.com/tipbk/sneakfeed-service/service"
	"github.com/tipbk/sneakfeed-service/util"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	PartiallyUpdateProfile(c *gin.Context)
	RefreshToken(c *gin.Context)
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
	_, err := h.userService.CreateUser(registerRequest.Username, registerRequest.Password)
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
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
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

func (h *userHandler) UpdateProfile(c *gin.Context) {
	currentUser, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	var request model.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	if request.ProfileImage == "" {
		c.JSON(http.StatusOK, util.GenerateSuccessResponse("nothing to update"))
		return
	}
	err = h.userService.UpdateProfile(currentUser.ID.Hex(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse("updated"))
}

func (h *userHandler) GetProfile(c *gin.Context) {
	currentUser, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(currentUser))
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

func (h *userHandler) PartiallyUpdateProfile(c *gin.Context) {
	currentUser, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}
	var request dto.PartiallyUpdateProfile
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}

	if request.ImageBase64 == "" {
		c.JSON(http.StatusUnauthorized, util.GenerateFailedResponse("imageBase64 cannot be empty"))
		return
	}

	uploadResponse, err := h.imageUploaderService.UploadImage(request.ImageBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	imageUrl := uploadResponse.Data.Url
	currentUser.ProfileImage = imageUrl
	err = h.userService.UpdateProfile(currentUser.ID.Hex(), currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse("updated"))
}
