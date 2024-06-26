package main

import (
	"backend/controllers"
	"backend/routes"
	"backend/tools"
	"backend/tools/db"
	"backend/tools/validation"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	// *Middleware for logging
	app.Use(logger.New())

	// !Landing Page Routes
	app.Get("/", routes.HandleRoot)

	// Workout builder
	app.Put("/updateWorkout", routes.UpdateCreateWorkout)
	app.Get("/exercises", routes.GetExerciseList)
	app.Get("/tags", routes.GetTagList)
	app.Get("/workouts", routes.GetWorkoutList)
	app.Get("/workout/detail", routes.GetWorkoutDetail)

	// FollowAlong
	app.Get("/trainingPlans", routes.GetTrainingPlan)

	// Chat
	app.Get("/contacts", routes.GetContacts)
	app.Get("/chat", websocket.New(routes.ChatHandler))
	app.Put("/shareWorkout", routes.ShareWorkout)

	//auth
	app.Post("/login", controllers.AttemptLogin)
	app.Post("/register", controllers.Register)
	app.Post("/continue-login", controllers.ContinueUserSession)
	app.Post("/logout", controllers.Logout)

	//BaseUser
	app.Get("/user", controllers.GetBaseUser)
	app.Patch("/user", controllers.UpdateBaseUser)
	app.Patch("/user/photo", controllers.UpdateBaseUserPhoto)
	app.Delete("/user", controllers.DeleteBaseUser)

	//Trainers
	app.Get("/trainers", controllers.GetAvailableTrainers)

	//Specialties
	app.Get("/specialties", controllers.GetAvailableSpecialties)
	app.Get("/specialties/selected", controllers.GetSelectedSpecialties)
	app.Post("/specialties", controllers.SetSpecialties)

	//Athletes
	app.Get("/measurements", controllers.GetMeasurements)
	app.Patch("/measurements", controllers.SetMeasurements)

	//Requests
	app.Get("/athleteRequests", controllers.GetAthleteRequests)
	app.Get("/trainerRequests", controllers.GetTrainerRequests)
	app.Post("/requestTraining", controllers.RequestTraining)
	app.Patch("/resolveRequest", controllers.ResolveRequest)

	//Forum
	app.Get("/forum", controllers.ListForums)
	app.Post("/forum", controllers.CreateForum)
	app.Get("/forum/:id", controllers.GetForum)

	app.Post("/post", controllers.CreatePost)
	app.Get("/post/:id", controllers.GetPost)
	app.Get("/post/:id/comments", controllers.GetComments)

	app.Get("/post/:id/like", controllers.GetLike)
	app.Post("/post/:id/like", controllers.LikePost)
	//app.Get("/post/:id/likes", controllers.GetTotalLikes)

	app.Post("/comment", controllers.CreateComment)

	db.Init()
	validation.Init()
	tools.InitSessions()
	app.Listen(":9090")
}
