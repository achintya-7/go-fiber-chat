package routes

import (
	"github.com/achintya-7/go-fiber-chat/controllers"
	"github.com/gofiber/fiber/v2"
)

func ChatRoute(app *fiber.App) {
	app.Post("/create_chat", controllers.CreateChat)
	app.Put("/add_to_group", controllers.AddToGroup)
	app.Delete("/delete_from_group", controllers.DeleteFromGroup)
	app.Get("/get_all_chats/:userId", controllers.GetAllChats)
	app.Get("/get_all_messages/:chatId", controllers.GetAllMessages)
	app.Post("/create_group_chat", controllers.CreateGroupChat)
}
