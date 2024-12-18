package models

import "time"

type Slider struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	Translations map[string]struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `bson:"translations" json:"translations"`
	MediaID   string            `bson:"media_id" json:"media_id"`
	Link      map[string]string `bson:"link" json:"link"`
	Order     int               `bson:"order" json:"order"`
	Status    string            `bson:"status" json:"status"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}
