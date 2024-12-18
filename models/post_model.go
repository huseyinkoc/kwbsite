package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" binding:"required"`
	Content   string             `bson:"content" json:"content" binding:"required"`
	Category  string             `bson:"category" json:"category"`
	Tags      []string           `bson:"tags" json:"tags"`
	Status    string             `bson:"status" json:"status"` // published, draft
	Author    string             `bson:"author" json:"author"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}
