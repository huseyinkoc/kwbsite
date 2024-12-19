package models

type Role struct {
	ID      string              `bson:"_id"`     // Rol adı
	Modules map[string][]string `bson:"modules"` // Modül bazlı izinler
}
