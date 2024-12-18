package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Page struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" binding:"required"`
	Content   string             `bson:"content" json:"content" binding:"required"`
	Status    string             `bson:"status" json:"status" binding:"required"` // "published" or "draft"
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}
