package service

import (
	"errors"
	"regexp"

	"github.com/tipbk/sneakfeed-service/model"
	"github.com/tipbk/sneakfeed-service/repository"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	CreateUser(username string, password string, email string) (*model.User, error)
	LoginUser(username string, password string) (*model.User, error)
	FindUserWithUserID(userID string) (*model.User, error)
	FindUserWithUsername(username string) (*model.User, error)
	FindUserViewByOthers(currentUserID, targetUsername string) (*model.UserViewByOthers, error)
	GetUsersByIDList(userIDs []string) ([]model.User, error)
	UpdateProfile(userID string, updatedUser *model.User) error
	ToggleFollowOnUser(userID string, followUserID string) (bool, error)
	IsUserFollowed(userID, followUserID string) (bool, error)
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(username string, password string, email string) (*model.User, error) {
	err := s.validateRegisterInput(username, password, email)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.CreateUser(username, password, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) LoginUser(username string, password string) (*model.User, error) {
	user, err := s.userRepository.LoginUser(username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) FindUserWithUserID(userID string) (*model.User, error) {
	user, err := s.userRepository.FindUserWithUserID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) FindUserWithUsername(username string) (*model.User, error) {
	user, err := s.userRepository.FindUserWithUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) FindUserViewByOthers(currentUserID, targetUsername string) (*model.UserViewByOthers, error) {
	user, err := s.userRepository.FindUserViewByOthers(currentUserID, targetUsername)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUsersByIDList(userIDs []string) ([]model.User, error) {
	users, err := s.userRepository.GetUsersByIDList(userIDs)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) UpdateProfile(userID string, updatedUser *model.User) error {
	err := s.userRepository.UpdateProfile(userID, updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) validateRegisterInput(username, password, email string) error {
	usernameRegex := `^[0-9a-z]{5,15}$`
	passwordRegex := `^[a-zA-Z0-9!@#$%^&*]{6,16}$`
	emailRegex := `^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`

	matched, err := regexp.Match(usernameRegex, []byte(username))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("username is invalid")
	}

	matched, err = regexp.Match(passwordRegex, []byte(password))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("password is invalid")
	}

	matched, err = regexp.Match(emailRegex, []byte(email))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("email is invalid")
	}

	return nil
}

func (s *userService) IsUserFollowed(userID, followUserID string) (bool, error) {
	isFollowed, err := s.userRepository.IsUserFollowed(userID, followUserID)
	if err != nil {
		return false, err
	}
	return isFollowed, nil
}

func (s *userService) ToggleFollowOnUser(userID string, followUserID string) (bool, error) {
	isFollowed, err := s.userRepository.IsUserFollowed(userID, followUserID)
	if err != nil {
		return false, err
	}
	if isFollowed { //do unfollow
		err := s.userRepository.UnfollowUser(userID, followUserID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else { // do follow
		_, err := s.userRepository.FollowUser(userID, followUserID)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}
