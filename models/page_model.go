package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID              primitive.ObjectID          `bson:"_id,omitempty"`
	Slug            string                      `bson:"slug"`                                    // SEO dostu URL
	Localizations   map[string]LocalizedContent `bson:"localizations"`                           // Dil kodu ile eşleştirilmiş içerik
	MetaTitle       string                      `bson:"meta_title"`                              // SEO Başlığı
	MetaDescription string                      `bson:"meta_description"`                        // SEO Açıklaması
	MetaKeywords    []string                    `bson:"meta_keywords"`                           // SEO Anahtar Kelimeler
	Status          string                      `bson:"status" json:"status" binding:"required"` // "published" or "draft"
	CreatedAt       time.Time                   `bson:"created_at"`
	UpdatedAt       time.Time                   `bson:"updated_at"`
}
