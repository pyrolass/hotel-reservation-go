package db

import (
	"context"

	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context context.Context, room *entities.Room) (*entities.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client) *MongoRoomStore {
	coll := client.Database(DBNAME).Collection(roomColl)

	return &MongoRoomStore{
		client: client,
		coll:   coll,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *entities.Room) (*entities.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)

	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)

	return room, nil

}
