package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user schema
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name" binding:"required"`       // Kullanıcının adı
	Surname           string             `bson:"surname" json:"surname" binding:"required"` // Kullanıcının soyadı
	FullName          string             `bson:"full_name" json:"full_name"`                // Otomatik oluşturulan tam ad
	Email             string             `bson:"email" json:"email" binding:"required,email"`
	PreferredLanguage string             `bson:"preferred_language" json:"preferred_language"` // Kullanıcı tercihi
	Username          string             `bson:"username" json:"username" binding:"required"`
	Password          string             `bson:"password" json:"password" binding:"required"`
	Roles             []string           `bson:"roles" json:"roles" binding:"required"` // ["admin", "editor", "user"]
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" example:"newpassword123"`
}

type RequestPasswordReset struct {
	Email string `json:"email" example:"abc@mail.com"`
}

type Login struct {
	Username string `json:"username" example:"mustafakemal"`
	Password string `json:"password" example:"ADsdsasWDD!!!8"`
}
