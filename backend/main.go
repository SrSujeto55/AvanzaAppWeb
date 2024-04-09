package main

import (
	"backend/controllers"
	"backend/routes"
	"backend/tools"
	"backend/tools/db"
	"backend/tools/validation"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	// *Middleware for logging
	app.Use(logger.New())

	// Allow CORS from our front-end server
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080",
		AllowCredentials: true,
	}))

	// !Landing Page Routes
	app.Get("/", routes.HandleRoot)

	// Chat
	app.Get("/contacts", routes.GetContacts)
	app.Get("/chat", websocket.New(routes.ChatHandler))

	//auth
	app.Post("/login", controllers.AttemptLogin)
	app.Post("/register", controllers.Register)
	app.Post("/continue-login", controllers.ContinueUserSession)

	//BaseUser
	//app.Post("/user") // handled by auth register ??
	app.Get("/user", controllers.GetBaseUser)
	app.Patch("/user", controllers.UpdateBaseUser)
	app.Delete("/user", controllers.DeleteBaseUser)

	db.Init()
	validation.Init()
	tools.InitSessions()
	app.Listen(":9090")
}
