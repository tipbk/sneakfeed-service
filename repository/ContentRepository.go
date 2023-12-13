package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tipbk/sneakfeed-service/config"
	"github.com/tipbk/sneakfeed-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContentRepository interface {
	CreatePost(userID string, content string, imageUrl *string) (string, error)
	AddComment(userID string, postID string, content string) (string, error)
	FindPost(postID string) (*model.Post, error)
	GetPosts(userID string) ([]model.PostDetail, error)
	GetPostByID(userID, postID string) (*model.PostDetail, error)
	GetCommentFromPostID(postID string) ([]model.Comment, error)
	IsPostLikeByUserID(userID string, postID string) (bool, error)
	LikePost(userID string, postID string) (string, error)
	UnlikePost(userID string, postID string) error
	CountLikeAndCommentOnPost(postID string) (int64, int64, error)
}

type contentRepository struct {
	envConfig   *config.EnvConfig
	mongoClient *mongo.Client
}

func NewContentReepository(envConfig *config.EnvConfig, mongoClient *mongo.Client) ContentRepository {
	return &contentRepository{
		envConfig:   envConfig,
		mongoClient: mongoClient,
	}
}

func (r *contentRepository) CreatePost(userID string, content string, imageUrl *string) (string, error) {
	now := time.Now()
	newPost := model.Post{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Content:         content,
		CreatedDatetime: &now,
		ImageUrl:        imageUrl,
	}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("post")
	result, err := collection.InsertOne(context.Background(), newPost)
	if err != nil {
		return "", errors.New("failed to create new post")
	}

	// converting primitive object to string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", errors.New("there are some errors when creating a new post")
}

func (r *contentRepository) FindPost(postID string) (*model.Post, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("post")

	postHex, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New("couldn't find a post")
	}
	var existingPost model.Post
	err = collection.FindOne(context.Background(), bson.M{"_id": postHex}).Decode(&existingPost)
	return &existingPost, err
}

func (r *contentRepository) GetCommentFromPostID(postID string) ([]model.Comment, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("comment")
	query := bson.M{"postID": postID}
	cursor, err := collection.Find(context.Background(), query)
	if err != nil {
		fmt.Println("Error finding comment with postID:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	var comments []model.Comment
	for cursor.Next(context.Background()) {
		var comment model.Comment
		err := cursor.Decode(&comment)
		if err != nil {
			fmt.Println("Error decoding comment:", err)
			continue
		}
		comments = append(comments, comment)
	}
	if err := cursor.Err(); err != nil {
		fmt.Println("Error iterating cursor:", err)
		return nil, err
	}
	return comments, nil
}

func (r *contentRepository) AddComment(userID string, postID string, content string) (string, error) {
	now := time.Now()
	newComment := model.Comment{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		PostID:          postID,
		Content:         content,
		CreatedDatetime: &now,
	}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("comment")
	result, err := collection.InsertOne(context.Background(), newComment)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("failed to add new comment")
	}

	// converting primitive object to string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", errors.New("there are some errors when adding a new comment")
}

func (r *contentRepository) IsPostLikeByUserID(userID string, postID string) (bool, error) {
	filter := bson.M{"userID": userID, "postID": postID}
	likeCollection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("like")
	var likePost model.LikePost
	err := likeCollection.FindOne(context.Background(), filter).Decode(&likePost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (r *contentRepository) LikePost(userID string, postID string) (string, error) {
	likePost := model.LikePost{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		PostID: postID,
	}
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("like")
	result, err := collection.InsertOne(context.Background(), likePost)
	if err != nil {
		return "", err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("there are some errors when liking a post")
	}

	return oid.Hex(), nil
}

func (r *contentRepository) UnlikePost(userID string, postID string) error {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("like")
	filter := bson.M{"userID": userID, "postID": postID}
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

func (r *contentRepository) CountLikeAndCommentOnPost(postID string) (int64, int64, error) {
	likeCollection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("like")
	filter := bson.M{"postID": postID}

	likeCount, err := likeCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, 0, err
	}

	commentCollection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("comment")
	commentCount, err := commentCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, 0, err
	}

	return likeCount, commentCount, nil
}

func (r *contentRepository) GetPosts(userID string) ([]model.PostDetail, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("post")

	pipeline := mongo.Pipeline{
		bson.D{{"$sort", bson.D{{"createdDatetime", -1}}}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"stringPostID", bson.D{{"$toString", "$_id"}}},
					{"objectUserID", bson.D{{"$toObjectId", "$userID"}}},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "like"},
					{"localField", "stringPostID"},
					{"foreignField", "postID"},
					{"as", "likeResult"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"objectUserID", "$objectUserID"},
					{"stringPostID", "$stringPostID"},
					{"totalLikes", bson.D{{"$size", "$likeResult"}}},
					{"isLike",
						bson.D{
							{"$ifNull",
								bson.A{
									bson.D{
										{"$in",
											bson.A{
												userID,
												"$likeResult.userID",
											},
										},
									},
									false,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "comment"},
					{"localField", "stringPostID"},
					{"foreignField", "postID"},
					{"as", "commentResult"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"objectUserID", "$objectUserID"},
					{"stringPostID", "$stringPostID"},
					{"totalLikes", "$totalLikes"},
					{"totalComments", bson.D{{"$size", "$commentResult"}}},
					{"isLike", "$isLike"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "user"},
					{"localField", "objectUserID"},
					{"foreignField", "_id"},
					{"as", "userResult"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"objectUserID", "$objectUserID"},
					{"stringPostID", "$stringPostID"},
					{"totalLikes", "$totalLikes"},
					{"totalComments", "$totalComments"},
					{"username", bson.D{{"$first", "$userResult.username"}}},
					{"profileImage", bson.D{{"$first", "$userResult.profileImage"}}},
					{"isLike", "$isLike"},
				},
			},
		},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		fmt.Println("Error creating cursor:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []model.PostDetail
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *contentRepository) GetPostByID(userID, postID string) (*model.PostDetail, error) {
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("post")

	pipeline := mongo.Pipeline{
		bson.D{
			{"$project",
				bson.D{
					{"stringPostId", bson.D{{"$toString", "$_id"}}},
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
				},
			},
		},
		bson.D{{"$match", bson.D{{"stringPostId", postID}}}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"stringPostID", bson.D{{"$toString", "$_id"}}},
					{"objectUserID", bson.D{{"$toObjectId", "$userID"}}},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "like"},
					{"localField", "stringPostID"},
					{"foreignField", "postID"},
					{"as", "likeResult"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"objectUserID", "$objectUserID"},
					{"stringPostID", "$stringPostID"},
					{"totalLikes", bson.D{{"$size", "$likeResult"}}},
					{"isLike",
						bson.D{
							{"$ifNull",
								bson.A{
									bson.D{
										{"$in",
											bson.A{
												userID,
												"$likeResult.userID",
											},
										},
									},
									false,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "comment"},
					{"localField", "stringPostID"},
					{"foreignField", "postID"},
					{"as", "commentResult"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"objectUserID", "$objectUserID"},
					{"stringPostID", "$stringPostID"},
					{"totalLikes", "$totalLikes"},
					{"totalComments", bson.D{{"$size", "$commentResult"}}},
					{"isLike", "$isLike"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "user"},
					{"localField", "objectUserID"},
					{"foreignField", "_id"},
					{"as", "userResult"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", "$_id"},
					{"content", "$content"},
					{"userID", "$userID"},
					{"createdDatetime", "$createdDatetime"},
					{"imageUrl", "$imageUrl"},
					{"objectUserID", "$objectUserID"},
					{"stringPostID", "$stringPostID"},
					{"totalLikes", "$totalLikes"},
					{"totalComments", "$totalComments"},
					{"username", bson.D{{"$first", "$userResult.username"}}},
					{"profileImage", bson.D{{"$first", "$userResult.profileImage"}}},
					{"isLike", "$isLike"},
				},
			},
		},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		fmt.Println("Error creating cursor:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []model.PostDetail
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	if len(results) <= 0 {
		return nil, errors.New("couldn't find a post")
	}
	return &results[0], nil
}
