package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	app := fiber.New()

	api := app.Group("api/v1")

	port := ":3000"

	api.Get(
		"/ping",
		func(c *fiber.Ctx) error {
			return c.JSON(
				map[string]any{
					"message": "pong",
				},
			)
		},
	)

	userHandle := handlers.NewUserHandler(db.NewMongoUserStore(client))

	api.Get("/user", userHandle.HandleGetUsers)

	api.Get("/user/:id", userHandle.HandleGetUser)

	app.Listen(port)
}
