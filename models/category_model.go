package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" binding:"required"`
	Slug      string             `bson:"slug" json:"slug"` // URL dostu isim
	CreatedAt int64              `bson:"created_at" json:"created_at"`
}
