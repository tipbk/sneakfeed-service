package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tipbk/sneakfeed-service/dto"
	"github.com/tipbk/sneakfeed-service/model"
	"github.com/tipbk/sneakfeed-service/service"
	"github.com/tipbk/sneakfeed-service/util"
)

type ContentHandler interface {
	CreatePost(c *gin.Context)
	AddComment(c *gin.Context)
	GetPosts(c *gin.Context)
	GetPostByID(c *gin.Context)
	GetCommentByPostID(c *gin.Context)
	ToggleLikePostByID(c *gin.Context)
	GetMetadata(c *gin.Context)
}

type contentHandler struct {
	contentService     service.ContentService
	userService        service.UserService
	imageUploadService service.ImageUploaderService
}

func NewContentHandler(contentService service.ContentService, userService service.UserService, imageUploadService service.ImageUploaderService) ContentHandler {
	return &contentHandler{
		contentService:     contentService,
		userService:        userService,
		imageUploadService: imageUploadService,
	}
}

func (h *contentHandler) CreatePost(c *gin.Context) {
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	var createPostRequest dto.CreatePostRequest
	if err := c.ShouldBindJSON(&createPostRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	if createPostRequest.Content == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("content cannot be empty"))
		return
	}

	var imageUrl *string
	if createPostRequest.ImageBase64 != nil {
		uploadResponse, err := h.imageUploadService.UploadImage(*createPostRequest.ImageBase64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
			return
		}
		imageUrl = &uploadResponse.Data.Url
	}
	postID, err := h.contentService.CreatePost(user.ID.Hex(), createPostRequest.Content, imageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(postID))
}

// comment on post
func (h *contentHandler) AddComment(c *gin.Context) {
	postID := c.Param("postID")
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	var addCommentRequest dto.AddCommentRequest
	if err := c.ShouldBindJSON(&addCommentRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("body parse error: invalid json"))
		return
	}
	if addCommentRequest.Content == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("content cannot be empty"))
		return
	}

	if postID == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("postID cannot be empty"))
		return
	}

	commentID, err := h.contentService.AddComment(user.ID.Hex(), postID, addCommentRequest.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(commentID))
}

func (h *contentHandler) GetPosts(c *gin.Context) {
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	timeFromString := c.Query("from")
	var timeFrom *time.Time

	if timeFromString != "" {
		t, err := util.ParseStringToTime(timeFromString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
			return
		}
		addedTime := t.Add(time.Millisecond * -1)
		timeFrom = &addedTime
	}

	limitString := c.Query("limit")
	limit := 5
	if limitString != "" {
		limit, err = util.ConvertStringToInt(limitString)
		if err != nil {
			limit = 5
		}
	}

	isFindFollowingPost := false
	filter := c.Query("filter")
	if filter == "FOLLOWING_POST" {
		isFindFollowingPost = true
	}

	posts, err := h.contentService.GetPosts(user.ID.Hex(), limit, timeFrom, isFindFollowingPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(posts))
}

func (h *contentHandler) GetPostByID(c *gin.Context) {
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	postID := c.Param("postID")
	postDetail, err := h.contentService.GetPostByID(user.ID.Hex(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(postDetail))
}

func (h *contentHandler) GetCommentByPostID(c *gin.Context) {
	postID := c.Param("postID")
	comments, err := h.contentService.GetCommentFromPostID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	// get user profile
	userIDs := []string{}
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
	}
	users, err := h.userService.GetUsersByIDList(userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	userMap := make(map[string]model.User)
	for _, user := range users {
		userMap[user.ID.Hex()] = user
	}
	responses := []dto.GetCommentByPostIDResponse{}
	for _, comment := range comments {
		responses = append(responses, dto.GetCommentByPostIDResponse{
			Content:         comment.Content,
			CreatedDatetime: comment.CreatedDatetime,
			UserID:          comment.UserID,
			Username:        userMap[comment.UserID].Username,
			DisplayName:     userMap[comment.UserID].DisplayName,
			ProfileImage:    userMap[comment.UserID].ProfileImage,
			CommentID:       comment.ID.Hex(),
		})

	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(responses))
}

func (h *contentHandler) ToggleLikePostByID(c *gin.Context) {
	user, err := util.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	postID := c.Param("postID")
	_, err = h.contentService.FindPost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	isLike, err := h.contentService.ToggleLikeOnPost(user.ID.Hex(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.GenerateSuccessResponse(dto.ToggleLikeResponse{IsLike: isLike}))
}

func (h *contentHandler) GetMetadata(c *gin.Context) {
	var request dto.MetadataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse(err.Error()))
		return
	}
	if request.Url == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("param url is needed"))
		return
	}
	resp, err := h.contentService.GetMetadata(request.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(resp))
}
