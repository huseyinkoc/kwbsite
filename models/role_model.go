package models

type Role struct {
	ID          string              `bson:"_id" json:"id"`
	Permissions map[string][]string `bson:"permissions" json:"permissions"` // Modül bazlı izinler
}
