package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents the user schema
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName string             `bson:"full_name" json:"full_name" binding:"required"`
	Email    string             `bson:"email" json:"email" binding:"required,email"`
	Username string             `bson:"username" json:"username" binding:"required"`
	Password string             `bson:"password" json:"password" binding:"required"`
	Role     string             `bson:"role" json:"role" binding:"required"` // admin, editor, etc.
}
