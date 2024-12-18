package models

type LocalizedContent struct {
	Title   string `bson:"title"`
	Content string `bson:"content"`
}
