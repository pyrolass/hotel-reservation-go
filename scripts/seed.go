package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Seeding database...")

	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }

	// uri := os.Getenv("MONGODB_URI")

	// if uri == "" {
	// 	log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	// }

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://pyro:LAS99ovi@cluster0.syyndk3.mongodb.net/hotel-reservation?retryWrites=true&w=majority&appName=Cluster0"))

	if err != nil {
		panic(err)
	}

	println("Connected to MongoDB!")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client)

	hotel := entities.Hotel{
		Name:     "Hilton",
		Location: "France",
	}

	room := entities.Room{
		Type:      entities.SinglePerson,
		BasePrice: 100,
		Occupied:  false,
	}

	ctx := context.Background()

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	room.HotelID = insertedHotel.ID

	insertedRoom, err := roomStore.InsertRoom(ctx, &room)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertedHotel)
	fmt.Println(insertedRoom)

}
