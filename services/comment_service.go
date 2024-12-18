package services

import (
	"admin-panel/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var commentCollection *mongo.Collection

func InitCommentService(client *mongo.Client) {
	commentCollection = client.Database("admin_panel").Collection("comments")
	log.Println("Comment service initialized with collection:", commentCollection.Name())
}

func CreateComment(ctx context.Context, comment *models.Comment) (*mongo.InsertOneResult, error) {
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return commentCollection.InsertOne(ctx, comment)
}

func GetCommentsByPostID(ctx context.Context, postID primitive.ObjectID) ([]models.Comment, error) {
	filter := bson.M{"post_id": postID}
	cursor, err := commentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	for cursor.Next(ctx) {
		var comment models.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func AddReply(ctx context.Context, commentID, replyID primitive.ObjectID) error {
	filter := bson.M{"_id": commentID}
	update := bson.M{"$push": bson.M{"replies": replyID}, "$set": bson.M{"updated_at": time.Now()}}
	_, err := commentCollection.UpdateOne(ctx, filter, update)
	return err
}

func LikeComment(ctx context.Context, commentID primitive.ObjectID) error {
	filter := bson.M{"_id": commentID}
	update := bson.M{"$inc": bson.M{"likes": 1}, "$set": bson.M{"updated_at": time.Now()}}
	_, err := commentCollection.UpdateOne(ctx, filter, update)
	return err
}
