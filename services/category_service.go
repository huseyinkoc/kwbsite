package services

import (
	"admin-panel/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection

func InitCategoryService(client *mongo.Client) {
	categoryCollection = client.Database("admin_panel").Collection("categories")
}

func CreateCategory(category models.Category) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	category.CreatedAt = time.Now().Unix()
	return categoryCollection.InsertOne(ctx, category)
}

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := categoryCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
