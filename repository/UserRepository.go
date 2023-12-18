package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tipbk/sneakfeed-service/config"
	"github.com/tipbk/sneakfeed-service/model"
	"github.com/tipbk/sneakfeed-service/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(username string, password string, email string) (*model.User, error)
	LoginUser(username string, password string) (*model.User, error)
	FindUserWithUserID(userID string) (*model.User, error)
	FindUserWithUsername(username string) (*model.User, error)
	GetUsersByIDList(userIDs []string) ([]model.User, error)
	UpdateProfile(userID string, updatedUser *model.User) error
	FollowUser(userID string, followUserID string) (string, error)
	UnfollowUser(userID string, followUserID string) error
	IsUserFollowed(userID string, followUserID string) (bool, error)
}

type userRepository struct {
	envConfig   *config.EnvConfig
	mongoClient *mongo.Client
}

func NewUserRepository(envConfig *config.EnvConfig, mongoClient *mongo.Client) UserRepository {
	return &userRepository{
		envConfig:   envConfig,
		mongoClient: mongoClient,
	}
}

func (r *userRepository) CreateUser(username string, password string, email string) (*model.User, error) {
	_, err := r.FindUser(username)
	if err == nil {
		return nil, errors.New("username already taken")
	}
	_, err = r.FindUserByEmail(email)
	if err == nil {
		return nil, errors.New("email already taken")
	}
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}
	newUser := model.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: hashPassword,
		Email:    email,
	}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")
	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to create user")
	}
	return &newUser, nil
}

func (r *userRepository) FindUser(username string) (*model.User, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")

	var existingUser model.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	return &existingUser, err
}

func (r *userRepository) FindUserByEmail(email string) (*model.User, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")

	var existingUser model.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&existingUser)
	return &existingUser, err
}

func (r *userRepository) FindUserWithUserID(userID string) (*model.User, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")

	userHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("couldn't find a user")
	}
	var existingUser model.User
	err = collection.FindOne(context.Background(), bson.M{"_id": userHex}).Decode(&existingUser)
	return &existingUser, err
}

func (r *userRepository) FindUserWithUsername(username string) (*model.User, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")

	var existingUser model.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)

	return &existingUser, err
}

func (r *userRepository) LoginUser(username string, password string) (*model.User, error) {
	user, err := r.FindUser(username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	if !util.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

func (r *userRepository) GetUsersByIDList(userIDs []string) ([]model.User, error) {
	userPrimitiveObjects := []primitive.ObjectID{}
	for _, userID := range userIDs {
		object, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		userPrimitiveObjects = append(userPrimitiveObjects, object)
	}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")
	query := bson.M{"_id": bson.M{"$in": userPrimitiveObjects}}
	cursor, err := collection.Find(context.Background(), query)
	if err != nil {
		fmt.Println("Error finding users:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	var users []model.User
	for cursor.Next(context.Background()) {
		var user model.User
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		fmt.Println("Error iterating cursor:", err)
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateProfile(userID string, updatedUser *model.User) error {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("user")
	refinedUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": refinedUserID}, bson.M{"$set": updatedUser})
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	return nil
}

func (r *userRepository) FollowUser(userID string, followUserID string) (string, error) {
	now := time.Now()
	follow := model.Follow{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		FollowUserID:    followUserID,
		CreatedDatetime: &now,
	}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("follow")
	result, err := collection.InsertOne(context.Background(), follow)
	if err != nil {
		return "", err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("there are some errors when following a user")
	}

	return oid.Hex(), nil
}

func (r *userRepository) UnfollowUser(userID string, followUserID string) error {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("follow")
	filter := bson.M{"userID": userID, "followUserID": followUserID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Println("Error deleting document:", err)
		return err
	}
	if result.DeletedCount == 1 {
		fmt.Println("Successfully deleted one document")
		return nil
	}
	return errors.New("no documents were deleted")
}

func (r *userRepository) IsUserFollowed(userID string, followUserID string) (bool, error) {
	filter := bson.M{"userID": userID, "followUserID": followUserID}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("follow")
	var follow model.Follow
	err := collection.FindOne(context.Background(), filter).Decode(&follow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
