package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user schema
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" binding:"required"`       // Kullanıcının adı
	Surname  string             `bson:"surname" json:"surname" binding:"required"` // Kullanıcının soyadı
	FullName string             `bson:"full_name" json:"full_name"`                // Otomatik oluşturulan tam ad
	Email    string             `bson:"email" json:"email" binding:"required,email"`
	Username string             `bson:"username" json:"username" binding:"required"`
	Password string             `bson:"password" json:"password" binding:"required"`
	Roles    []string           `bson:"roles" json:"roles" binding:"required"` // ["admin", "editor", "user"]
}
