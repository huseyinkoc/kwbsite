package services

import (
	"admin-panel/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection

func InitPostService(client *mongo.Client) {
	postCollection = client.Database("admin_panel").Collection("posts")
}

// CreatePost creates a new post
func CreatePost(ctx context.Context, post *models.Post) error {
	post.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	post.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := postCollection.InsertOne(ctx, post)
	return err
}

// GetAllPosts retrieves all posts
func GetAllPosts(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post
	cursor, err := postCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostByID retrieves a single post by its ID
func GetPostByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
	var post models.Post
	err := postCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	return &post, err
}

// GetFilteredPosts retrieves posts based on filters
func GetFilteredPosts(ctx context.Context, filter bson.M) ([]models.Post, error) {
	var posts []models.Post
	cursor, err := postCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost updates an existing post
func UpdatePost(ctx context.Context, post *models.Post) error {
	post.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := postCollection.UpdateOne(
		ctx,
		bson.M{"_id": post.ID},
		bson.M{"$set": post},
	)
	return err
}
