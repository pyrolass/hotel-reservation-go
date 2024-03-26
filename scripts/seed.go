package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Seeding database...")

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

	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		panic(err)
	}

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

	rooms := []entities.Room{
		{
			Type:      entities.SinglePerson,
			BasePrice: 100,
			Occupied:  false,
		},
		{
			Type:      entities.DoublePerson,
			BasePrice: 12.99,
			Occupied:  false,
		},
		{
			Type:      entities.TriplePerson,
			BasePrice: 99.99,
			Occupied:  false,
		},
	}

	ctx := context.Background()

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	var roomIds = []primitive.ObjectID{}

	for _, room := range rooms {

		room.HotelID = insertedHotel.ID

		insertedRoom, err := roomStore.InsertRoom(ctx, &room)

		if err != nil {
			log.Fatal(err)
		}

		roomIds = append(roomIds, insertedRoom.ID)

		fmt.Println(insertedRoom)

	}

	hotelStore.UpdateHotelRoomIds(ctx, insertedHotel.ID.Hex(), roomIds)

	fmt.Println(insertedHotel)

}
