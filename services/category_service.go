package services

import (
	"admin-panel/models"
	"admin-panel/utils"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection

func InitCategoryService(client *mongo.Client) {
	categoryCollection = client.Database("admin_panel").Collection("categories")
}

func CreateCategory(category models.Category) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Slug oluşturma
	if category.Slug == nil {
		category.Slug = make(map[string]string)
	}

	for lang, localization := range category.Localizations {
		if localization.Title != "" {
			category.Slug[lang] = utils.GenerateSlug(localization.Title)
		} else if lang == "en" {
			return nil, errors.New("English name is required for slug generation")
		}
	}

	category.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
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
