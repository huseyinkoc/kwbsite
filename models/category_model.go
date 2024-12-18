package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID            primitive.ObjectID          `bson:"_id,omitempty"`
	Slug          string                      `bson:"slug"`                // SEO dostu URL
	Localizations map[string]LocalizedContent `bson:"localizations"`       // Dil kodu ile eşleştirilmiş kategori adı ve açıklaması
	ParentID      *primitive.ObjectID         `bson:"parent_id,omitempty"` // Alt kategori desteği
	CreatedAt     time.Time                   `bson:"created_at"`
	UpdatedAt     time.Time                   `bson:"updated_at"`
}
