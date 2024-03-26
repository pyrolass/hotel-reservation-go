package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/handlers"
	"github.com/pyrolass/hotel-reservation-go/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(router fiber.Router, client *mongo.Client) {

	userHandle := handlers.NewUserHandler(db.NewMongoUserStore(client))

	router.Get("/user", middleware.JWTAuthentication, userHandle.HandleGetUsers)
	router.Get("/user/:id", userHandle.HandleGetUser)
	router.Post("/user", userHandle.HandlePostUser)
	router.Delete("/user/:id", userHandle.HandleDeleteUser)
	router.Put("/user/:id", userHandle.HandlePutUser)
}
