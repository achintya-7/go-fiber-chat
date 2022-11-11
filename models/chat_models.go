package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateChatReq struct {
	UserId       primitive.ObjectID `json:"userId"`
	SecondUserId primitive.ObjectID `json:"secondUserId"`
}

type CreateChatRes struct {
	ChatId          primitive.ObjectID   `json:"chatId"`
	Users           []primitive.ObjectID `json:"users"`
	IsGroup         bool                 `json:"isGroup"`
	LatestMessage   string               `json:"latestMessage"`
	LatestMessageId string               `json:"latestMessageId"`
	UserId          primitive.ObjectID   `json:"userId"`
	ChatName        string               `json:"chatName"`
}

type GetAllChatsReq struct {
	UserId primitive.ObjectID `json:"userId"`
}

type GetAllChatsRes struct {
	UserId primitive.ObjectID `json:"userId"`
	Chats  []CreateChatRes    `json:"chats"`
}

type AddToGroupReq struct {
	Users  []primitive.ObjectID `json:"users"`
	ChatId primitive.ObjectID   `json:"chatId"`
}

type DeleteFromGroupReq struct {
	UserId primitive.ObjectID `json:"userId"`
	ChatId primitive.ObjectID `json:"chatId"`
}

type CreateGroupChatReq struct {
	UserId   primitive.ObjectID   `json:"userId"`
	Users    []primitive.ObjectID `json:"users"`
	ChatName string               `json:"chatName"`
}

type CreateGroupChatRes struct {
	ChatId          primitive.ObjectID   `json:"chatId"`
	Users           []primitive.ObjectID `json:"users"`
	IsGroup         bool                 `json:"isGroup"`
	LatestMessage   string               `json:"latestMessage"`
	LatestMessageId string               `json:"latestMessageId"`
	UserId          primitive.ObjectID   `json:"userId"`
	ChatName        string               `json:"chatName"`
}
