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
	CreatePost(userID string, content string, imageUrl *string, ogTitle *string, ogDescription *string, ogLink *string, ogImage *string, ogDomain *string) (string, error)
	AddComment(userID string, postID string, content string) (string, error)
	FindPost(postID string) (*model.Post, error)
	GetPosts(userID string, limit int, timeFrom *time.Time, postFilter, username string) (*model.PostDetailPagination, error)
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

func (r *contentRepository) CreatePost(userID string, content string, imageUrl *string, ogTitle *string, ogDescription *string, ogLink *string, ogImage *string, ogDomain *string) (string, error) {
	now := time.Now()
	newPost := model.Post{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Content:         content,
		CreatedDatetime: &now,
		ImageUrl:        imageUrl,
		OgTitle:         ogTitle,
		OgLink:          ogLink,
		OgDescription:   ogDescription,
		OgDomain:        ogDomain,
		OgImage:         ogImage,
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

func (r *contentRepository) GetPosts(userID string, limit int, timeFrom *time.Time, postFilter, username string) (*model.PostDetailPagination, error) {
	fmt.Println(postFilter)
	fmt.Println(username)
	fmt.Println(timeFrom)
	collection := r.mongoClient.Database(r.envConfig.DatabaseName).Collection("post")
	// for following posts
	projectCurrentUserAsString := bson.D{
		{"$project",
			bson.D{
				{"userID", "$userID"},
				{"title", "$title"},
				{"content", "$content"},
				{"imageUrl", "$imageUrl"},
				{"createdDatetime", "$createdDatetime"},
				{"currentUserID", userID},
				{"ogTitle", "$ogTitle"},
				{"ogDescription", "$ogDescription"},
				{"ogLink", "$ogLink"},
				{"ogImage", "$ogImage"},
				{"ogDomain", "$ogDomain"},
			},
		},
	}
	mergingFollowStage := bson.D{
		{"$lookup",
			bson.D{
				{"from", "follow"},
				{"localField", "currentUserID"},
				{"foreignField", "userID"},
				{"as", "results"},
			},
		},
	}

	getFollowingListStage := bson.D{
		{"$project",
			bson.D{
				{"userID", "$userID"},
				{"title", "$title"},
				{"content", "$content"},
				{"imageUrl", "$imageUrl"},
				{"createdDatetime", "$createdDatetime"},
				{"currentUserID", "$currentUserID"},
				{"followUserList", "$results.followUserID"},
				{"ogTitle", "$ogTitle"},
				{"ogDescription", "$ogDescription"},
				{"ogLink", "$ogLink"},
				{"ogImage", "$ogImage"},
				{"ogDomain", "$ogDomain"},
			},
		},
	}

	getMatchingPostStage := bson.D{
		{"$match",
			bson.D{
				{"$or",
					bson.A{
						bson.D{
							{"$expr",
								bson.D{
									{"$in",
										bson.A{
											"$userID",
											"$followUserList",
										},
									},
								},
							},
						},
						bson.D{
							{"$expr",
								bson.D{
									{"$eq",
										bson.A{
											"$currentUserID",
											"$userID",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// end for following posts
	sortingStage := bson.D{{"$sort", bson.D{{"createdDatetime", -1}}}}
	projectConversionForSearchingStage := bson.D{
		{"$project",
			bson.D{
				{"_id", "$_id"},
				{"content", "$content"},
				{"userID", "$userID"},
				{"createdDatetime", "$createdDatetime"},
				{"imageUrl", "$imageUrl"},
				{"stringPostID", bson.D{{"$toString", "$_id"}}},
				{"objectUserID", bson.D{{"$toObjectId", "$userID"}}},
				{"ogTitle", "$ogTitle"},
				{"ogDescription", "$ogDescription"},
				{"ogLink", "$ogLink"},
				{"ogImage", "$ogImage"},
				{"ogDomain", "$ogDomain"},
			},
		},
	}
	timeAfterStage := bson.D{{"$match", bson.D{{"createdDatetime", bson.D{{"$lte", timeFrom}}}}}}
	likeMergingStage := bson.D{
		{"$lookup",
			bson.D{
				{"from", "like"},
				{"localField", "stringPostID"},
				{"foreignField", "postID"},
				{"as", "likeResult"},
			},
		},
	}

	projectCountingLikeStage := bson.D{
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
				{"ogTitle", "$ogTitle"},
				{"ogDescription", "$ogDescription"},
				{"ogLink", "$ogLink"},
				{"ogImage", "$ogImage"},
				{"ogDomain", "$ogDomain"},
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
	}

	commentMergingStage := bson.D{
		{"$lookup",
			bson.D{
				{"from", "comment"},
				{"localField", "stringPostID"},
				{"foreignField", "postID"},
				{"as", "commentResult"},
			},
		},
	}

	projectCountingCommentStage := bson.D{
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
				{"ogTitle", "$ogTitle"},
				{"ogDescription", "$ogDescription"},
				{"ogLink", "$ogLink"},
				{"ogImage", "$ogImage"},
				{"ogDomain", "$ogDomain"},
				{"isComment",
					bson.D{
						{"$ifNull",
							bson.A{
								bson.D{
									{"$in",
										bson.A{
											userID,
											"$commentResult.userID",
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
	}

	userMergingStage := bson.D{
		{"$lookup",
			bson.D{
				{"from", "user"},
				{"localField", "objectUserID"},
				{"foreignField", "_id"},
				{"as", "userResult"},
			},
		},
	}

	projectUserMappingStage := bson.D{
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
				{"displayName", bson.D{{"$first", "$userResult.displayName"}}},
				{"profileImage", bson.D{{"$first", "$userResult.profileImage"}}},
				{"isLike", "$isLike"},
				{"isComment", "$isComment"},
				{"ogTitle", "$ogTitle"},
				{"ogDescription", "$ogDescription"},
				{"ogLink", "$ogLink"},
				{"ogImage", "$ogImage"},
				{"ogDomain", "$ogDomain"},
			},
		},
	}

	paginationQueryStage := bson.D{
		{"$facet",
			bson.D{
				{"pagination",
					bson.A{
						bson.D{{"$count", "total"}},
						bson.D{
							{"$addFields",
								bson.D{
									{"limit", limit},
								},
							},
						},
					},
				},
				{"data",
					bson.A{
						bson.D{{"$limit", limit}},
					},
				},
			},
		},
	}

	paginationExtractingstage := bson.D{
		{"$project",
			bson.D{
				{"pagination",
					bson.D{
						{"$arrayElemAt",
							bson.A{
								"$pagination",
								0,
							},
						},
					},
				},
				{"posts", "$data"},
			},
		},
	}

	matchUserStage := bson.D{{"$match", bson.D{{"username", username}}}}

	// DEFAULT FILTER (GLOBAL FEEDS)
	pipeline := mongo.Pipeline{
		sortingStage,
		projectConversionForSearchingStage,
		likeMergingStage,
		projectCountingLikeStage,
		commentMergingStage,
		projectCountingCommentStage,
		userMergingStage,
		projectUserMappingStage,
		paginationQueryStage,
		paginationExtractingstage,
	}

	if timeFrom != nil {
		pipeline = mongo.Pipeline{
			sortingStage,
			timeAfterStage,
			projectConversionForSearchingStage,
			likeMergingStage,
			projectCountingLikeStage,
			commentMergingStage,
			projectCountingCommentStage,
			userMergingStage,
			projectUserMappingStage,
			paginationQueryStage,
			paginationExtractingstage,
		}
	}

	// FOLLOWING FEED FILTER
	if postFilter == "FOLLOWING_POST" {
		pipeline = mongo.Pipeline{
			projectCurrentUserAsString,
			mergingFollowStage,
			getFollowingListStage,
			getMatchingPostStage,
			sortingStage,
			projectConversionForSearchingStage,
			likeMergingStage,
			projectCountingLikeStage,
			commentMergingStage,
			projectCountingCommentStage,
			userMergingStage,
			projectUserMappingStage,
			paginationQueryStage,
			paginationExtractingstage,
		}
		if timeFrom != nil {
			pipeline = mongo.Pipeline{
				projectCurrentUserAsString,
				mergingFollowStage,
				getFollowingListStage,
				getMatchingPostStage,
				sortingStage,
				timeAfterStage,
				projectConversionForSearchingStage,
				likeMergingStage,
				projectCountingLikeStage,
				commentMergingStage,
				projectCountingCommentStage,
				userMergingStage,
				projectUserMappingStage,
				paginationQueryStage,
				paginationExtractingstage,
			}
		}
	} else if postFilter == "USER" {
		if username == "" {
			return nil, errors.New("username cannot be empty")
		}
		pipeline = mongo.Pipeline{
			sortingStage,
			projectConversionForSearchingStage,
			likeMergingStage,
			projectCountingLikeStage,
			commentMergingStage,
			projectCountingCommentStage,
			userMergingStage,
			projectUserMappingStage,
			matchUserStage,
			paginationQueryStage,
			paginationExtractingstage,
		}

		if timeFrom != nil {
			pipeline = mongo.Pipeline{
				sortingStage,
				timeAfterStage,
				projectConversionForSearchingStage,
				likeMergingStage,
				projectCountingLikeStage,
				commentMergingStage,
				projectCountingCommentStage,
				userMergingStage,
				projectUserMappingStage,
				matchUserStage,
				paginationQueryStage,
				paginationExtractingstage,
			}
		}
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		fmt.Println("Error creating cursor:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []model.PostDetailPagination
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	if len(results) <= 0 {
		return nil, errors.New("couldn't find a post")
	}
	fmt.Println(results)

	return &results[0], nil
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
					{"imageUrl", "$imageUrl"},
					{"createdDatetime", "$createdDatetime"},
					{"ogTitle", "$ogTitle"},
					{"ogDescription", "$ogDescription"},
					{"ogLink", "$ogLink"},
					{"ogImage", "$ogImage"},
					{"ogDomain", "$ogDomain"},
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
					{"ogTitle", "$ogTitle"},
					{"ogDescription", "$ogDescription"},
					{"ogLink", "$ogLink"},
					{"ogImage", "$ogImage"},
					{"ogDomain", "$ogDomain"},
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
					{"ogTitle", "$ogTitle"},
					{"ogDescription", "$ogDescription"},
					{"ogLink", "$ogLink"},
					{"ogImage", "$ogImage"},
					{"ogDomain", "$ogDomain"},
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
					{"ogTitle", "$ogTitle"},
					{"ogDescription", "$ogDescription"},
					{"ogLink", "$ogLink"},
					{"ogImage", "$ogImage"},
					{"ogDomain", "$ogDomain"},
					{"isComment",
						bson.D{
							{"$ifNull",
								bson.A{
									bson.D{
										{"$in",
											bson.A{
												userID,
												"$commentResult.userID",
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
					{"displayName", bson.D{{"$first", "$userResult.displayName"}}},
					{"profileImage", bson.D{{"$first", "$userResult.profileImage"}}},
					{"isLike", "$isLike"},
					{"isComment", "$isComment"},
					{"ogTitle", "$ogTitle"},
					{"ogDescription", "$ogDescription"},
					{"ogLink", "$ogLink"},
					{"ogImage", "$ogImage"},
					{"ogDomain", "$ogDomain"},
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
