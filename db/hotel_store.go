package db

import (
	"context"

	"github.com/pyrolass/hotel-reservation-go/common"
	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context context.Context, hotel *entities.Hotel) (*entities.Hotel, error)
	UpdateHotelRoomIds(context context.Context, id string, roomIds []primitive.ObjectID) error
	GetHotels(context context.Context) ([]*entities.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	coll := client.Database(DBNAME).Collection(hotelColl)

	return &MongoHotelStore{
		client: client,
		coll:   coll,
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *entities.Hotel) (*entities.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)

	if err != nil {
		return nil, err
	}

	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil

}

func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*entities.Hotel, error) {

	var hotels []*entities.Hotel

	cur, err := s.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil

}

func (s *MongoHotelStore) UpdateHotelRoomIds(context context.Context, id string, roomIds []primitive.ObjectID) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return common.InvalidID()
	}

	_, err = s.coll.UpdateOne(context,
		bson.M{
			"_id": oid,
		},
		bson.M{
			"$set": bson.M{
				"rooms": roomIds,
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}
