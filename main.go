package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/pyrolass/hotel-reservation-go/common"
	"github.com/pyrolass/hotel-reservation-go/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if apiError, ok := err.(common.ApiError); ok {
			return c.Status(apiError.Code).JSON(
				map[string]any{
					"error": apiError.Message,
				},
			)
		}

		return c.Status(500).JSON(
			map[string]any{
				"error": err.Error(),
			},
		)
	},
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	println("Connected to MongoDB!")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := fiber.New(config)

	port := ":3000"

	router := app.Group("api/v1")

	router.Get(
		"/ping",
		func(c *fiber.Ctx) error {
			return c.JSON(
				map[string]any{
					"message": "pong",
				},
			)
		},
	)

	routes.UserRoutes(router, client)
	routes.HotelRoutes(router, client)

	app.Listen(port)
}
