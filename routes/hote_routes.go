package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

func HotelRoutes(router fiber.Router, client *mongo.Client) {

	hotelHandler := handlers.NewHotelHandler(db.NewMongoHotelStore(client), db.NewMongoRoomStore(client))

	router.Get("/hotels", hotelHandler.HandleGetHotels)

}
