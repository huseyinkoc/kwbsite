package services

import (
	"admin-panel/models" // models paketinden Post modelini alıyoruz
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection

func InitPostService(client *mongo.Client) {
	postCollection = client.Database("admin_panel").Collection("posts")
}

func CreatePost(post models.Post) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post.CreatedAt = time.Now().Unix()
	return postCollection.InsertOne(ctx, post)
}

func GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func GetFilteredPosts(filter bson.M) ([]models.Post, error) {
	var posts []models.Post
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
