package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Media struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileName   string             `bson:"file_name" json:"file_name"`
	FilePath   string             `bson:"file_path" json:"file_path"`
	FileType   string             `bson:"file_type" json:"file_type"`
	UploadedAt int64              `bson:"uploaded_at" json:"uploaded_at"`
}
