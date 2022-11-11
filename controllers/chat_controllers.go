package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/achintya-7/go-fiber-chat/configs"
	"github.com/achintya-7/go-fiber-chat/models"
	"github.com/achintya-7/go-fiber-chat/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var chatCollection *mongo.Collection = configs.GetCollection(configs.DB, "chats")
var messageCollection *mongo.Collection = configs.GetCollection(configs.DB, "messages")

func CreateChat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chat models.CreateChatReq
	var resChat models.CreateChatRes

	if err := c.BodyParser(&chat); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    &fiber.Map{},
		})
	}

	chatFilter := bson.D{{Key: "users", Value: bson.A{chat.UserId, chat.SecondUserId}}}
	isChat := chatCollection.FindOne(ctx, chatFilter)

	// no document was found
	if isChat.Err() != nil {

		chatNew := models.CreateChatRes{
			ChatId:          primitive.NewObjectID(),
			IsGroup:         false,
			Users:           []primitive.ObjectID{chat.UserId, chat.SecondUserId},
			LatestMessage:   "",
			LatestMessageId: "",
			UserId:          chat.UserId,
			ChatName:        "",
		}

		result, err := chatCollection.InsertOne(ctx, chatNew)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data: &fiber.Map{
					"data": result,
				},
			})
		}

		return c.Status(http.StatusOK).JSON(responses.UserResponse{
			Status:  200,
			Message: "Chat Room Created",
			Data: &fiber.Map{
				"data": chatNew,
			},
		})
	}

	// a chat was found
	isChat.Decode(&resChat)
	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  200,
		Message: "A Chat Room already exist",
		Data: &fiber.Map{
			"data": resChat,
		},
	})
}

func AddToGroup(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req models.AddToGroupReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Unable to parse JSON",
				Data: &fiber.Map{
					"data": &fiber.Map{},
				},
			})
	}

	filter := bson.D{{Key: "chatid", Value: req.ChatId}, {Key: "isgroup", Value: true}}
	update := bson.D{
		{
			Key: "$addToSet",
			Value: bson.D{
				{
					Key:   "users",
					Value: bson.D{
						{
							Key: "$each", 
							Value: req.Users,
						},
					},
				},
			},
		},
	}

	if err := chatCollection.FindOneAndUpdate(ctx, filter, update); err.Err() != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Unable to find Group Chat",
				Data: &fiber.Map{
					"data": &fiber.Map{},
				},
			})
	}

	return c.Status(200).JSON(
		responses.UserResponse{
			Status:  200,
			Message: "User Added",
			Data: &fiber.Map{
				"data": req,
			},
		})
}

func DeleteFromGroup(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req models.DeleteFromGroupReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Unable to parse JSON",
				Data: &fiber.Map{
					"data": &fiber.Map{},
				},
			})
	}

	filter := bson.D{{Key: "chatid", Value: req.ChatId}, {Key: "isgroup", Value: true}}
	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{
					Key:   "users",
					Value: req.UserId,
				},
			},
		},
	}

	err := chatCollection.FindOneAndUpdate(ctx, filter, update)
	if err.Err() != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Chat Group not found",
				Data: &fiber.Map{
					"data": &fiber.Map{},
				},
			})
	}

	return c.Status(200).JSON(
		responses.UserResponse{
			Status:  200,
			Message: "User Removed",
			Data: &fiber.Map{
				"data": req.UserId,
			},
		})
}

func GetAllMessages(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatId := c.Params("chatId")
	objId, _ := primitive.ObjectIDFromHex(chatId)

	var messages []models.Message

	filter := bson.D{{Key: "roomid", Value: objId}}
	cursor, err := messageCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Unable to find Messages",
				Data: &fiber.Map{
					"data": &fiber.Map{
						"error": err.Error(),
					},
				},
			})
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var singleMessage models.Message
		if err = cursor.Decode(&singleMessage); err != nil {
			continue
		}
		messages = append(messages, singleMessage)
	}

	return c.Status(200).JSON(
		responses.UserResponse{
			Status:  200,
			Message: "Messages Found",
			Data: &fiber.Map{
				"data": messages,
			},
		})
}

func GetAllChats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	var chats []models.CreateChatRes

	objId, _ := primitive.ObjectIDFromHex(userId)

	results, err := chatCollection.Find(ctx, bson.D{{Key: "users", Value: objId}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Database Error",
				Data: &fiber.Map{
					"data": &fiber.Map{},
				},
			})
	}

	defer results.Close(ctx)

	i := 0

	for results.Next(ctx) {
		i++
		var singleChat models.CreateChatRes
		if err = results.Decode(&singleChat); err != nil {
			fmt.Println(err)
			continue
		}
		chats = append(chats, singleChat)
	}

	if len(chats) < 1 {
		return c.Status(401).JSON(
			responses.UserResponse{
				Status:  401,
				Message: "No Chats were found",
				Data: &fiber.Map{
					"data": chats,
				},
			})

	}

	messageRes := fmt.Sprintf("%d Chats were found", len(chats))

	return c.Status(200).JSON(
		responses.UserResponse{
			Status:  200,
			Message: messageRes,
			Data: &fiber.Map{
				"data":  chats,
				"index": i,
			},
		})
}

func CreateGroupChat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req models.CreateGroupChatReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Unable to parse JSON",
				Data: &fiber.Map{
					"data": &fiber.Map{},
				},
			})
	}

	if len(req.Users) < 1 {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Cant make a group with 1 user",
			Data: &fiber.Map{
				"data": &fiber.Map{},
			},
		})
	}

	req.Users = append(req.Users, req.UserId)

	chatNew := models.CreateGroupChatRes{
		ChatId:          primitive.NewObjectID(),
		IsGroup:         true,
		Users:           req.Users,
		UserId:          req.UserId,
		LatestMessage:   "",
		LatestMessageId: "",
		ChatName:        req.ChatName,
	}

	result, err := chatCollection.InsertOne(ctx, chatNew)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data: &fiber.Map{
				"data": result,
			},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  200,
		Message: "Chat Room Created",
		Data: &fiber.Map{
			"data": chatNew,
		},
	})

}

// func CreateGroupChat(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// 	defer cancel()

// 	var req models.DeleteFromGroupReq
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"Status":  err,
// 			"Message": err.Error(),
// 			"Data":    err,
// 		})
// 	}
// }
