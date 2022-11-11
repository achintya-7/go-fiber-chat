package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	UserId      primitive.ObjectID `json:"userId"`
	RoomId      primitive.ObjectID `json:"roomId"`
	Content     string             `json:"content"`
	ContentType string             `json:"contentType"`
	MessageId   string             `json:"messageId"`
	Timestamp   int64              `json:"timestamp"`
}
