package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tipbk/sneakfeed-service/config"
	"github.com/tipbk/sneakfeed-service/dto"
	"github.com/tipbk/sneakfeed-service/model"
	"github.com/tipbk/sneakfeed-service/repository"
)

type ContentService interface {
	CreatePost(userID string, content string, imageUrl *string) (string, error)
	AddComment(userID string, postID string, content string) (string, error)
	GetPosts(userID string, limit int, timeFrom *time.Time, isFollowingPost bool) (*model.PostDetailPagination, error)
	GetPostByID(userID, postID string) (*model.PostDetail, error)
	GetCommentFromPostID(postID string) ([]model.Comment, error)
	FindPost(postID string) (*model.Post, error)
	ToggleLikeOnPost(userID string, postID string) (bool, error)
	CountLikeAndCommentOnPost(postID string) (int64, int64, error)
	GetMetadata(targetUrl string) (*dto.MetadataExternal, error)
}

type contentService struct {
	envConfig         *config.EnvConfig
	contentRepository repository.ContentRepository
}

func NewContentService(envConfig *config.EnvConfig, contentRepository repository.ContentRepository) ContentService {
	return &contentService{
		envConfig:         envConfig,
		contentRepository: contentRepository,
	}
}

func (s *contentService) CreatePost(userID string, content string, imageUrl *string) (string, error) {
	postID, err := s.contentRepository.CreatePost(userID, content, imageUrl)
	if err != nil {
		return "", err
	}
	return postID, nil
}

func (s *contentService) AddComment(userID string, postID string, content string) (string, error) {
	post, err := s.contentRepository.FindPost(postID)
	if err != nil {
		return "", errors.New("couldn't find post")
	}
	commentID, err := s.contentRepository.AddComment(userID, post.ID.Hex(), content)
	if err != nil {
		return "", err
	}
	return commentID, nil
}

func (s *contentService) GetCommentFromPostID(postID string) ([]model.Comment, error) {
	comments, err := s.contentRepository.GetCommentFromPostID(postID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *contentService) FindPost(postID string) (*model.Post, error) {
	post, err := s.contentRepository.FindPost(postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *contentService) ToggleLikeOnPost(userID string, postID string) (bool, error) {
	isLike, err := s.contentRepository.IsPostLikeByUserID(userID, postID)
	if err != nil {
		return false, err
	}
	if isLike { //do unlike
		err := s.contentRepository.UnlikePost(userID, postID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else { // do like
		_, err := s.contentRepository.LikePost(userID, postID)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func (s *contentService) CountLikeAndCommentOnPost(postID string) (int64, int64, error) {
	likeCount, commentCount, err := s.contentRepository.CountLikeAndCommentOnPost(postID)
	if err != nil {
		return 0, 0, err
	}
	return likeCount, commentCount, nil
}

func (s *contentService) GetPosts(userID string, limit int, timeFrom *time.Time, isFollowingPost bool) (*model.PostDetailPagination, error) {
	posts, err := s.contentRepository.GetPosts(userID, limit, timeFrom, isFollowingPost)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *contentService) GetPostByID(userID, postID string) (*model.PostDetail, error) {
	post, err := s.contentRepository.GetPostByID(userID, postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *contentService) GetMetadata(targetUrl string) (*dto.MetadataExternal, error) {
	client := &http.Client{}
	requestedUrl := fmt.Sprintf("%s/api/metadata/%s", s.envConfig.MetadataEndpoint, url.QueryEscape(targetUrl))
	req, err := http.NewRequest("GET", requestedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("referer", requestedUrl)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return nil, errors.New("response status error from metadata service")
	}
	var response dto.MetadataExternal
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
