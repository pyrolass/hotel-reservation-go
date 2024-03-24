package db

import (
	"context"

	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserById(ctx context.Context, id string) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	coll := client.Database(DBNAME).Collection(userColl)

	return &MongoUserStore{
		client: client,
		coll:   coll,
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*entities.User, error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var user entities.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	var users []entities.User

	cursor, err := s.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
