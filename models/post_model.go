package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID            primitive.ObjectID        `bson:"_id,omitempty" json:"id"`
	Localizations map[string]LocalizedField `bson:"localizations" json:"localizations"` // Dil koduna göre içerik
	Status        string                    `bson:"status" json:"status"`               // draft, published, scheduled
	CategoryIDs   []primitive.ObjectID      `bson:"category_ids" json:"category_ids"`
	TagIDs        []primitive.ObjectID      `bson:"tag_ids" json:"tag_ids"`
	PublishDate   *primitive.DateTime       `bson:"publish_date,omitempty" json:"publish_date,omitempty"` // Yayınlanma tarihi
	AuthorID      primitive.ObjectID        `bson:"author_id" json:"author_id"`
	MetaTags      map[string]MetaTag        `bson:"meta_tags" json:"meta_tags"` // Dil kodu ve SEO bilgileri
	CreatedAt     primitive.DateTime        `bson:"created_at" json:"created_at"`
	UpdatedAt     primitive.DateTime        `bson:"updated_at" json:"updated_at"`
	CreatedBy     string                    `bson:"created_by" json:"created_by"`
	UpdatedBy     string                    `bson:"updated_by" json:"updated_by"`
}
