package db

import (
	"context"
	"fmt"

	"github.com/pyrolass/hotel-reservation-go/common"
	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(ctx context.Context) error
}
type UserStore interface {
	GetUserById(ctx context.Context, id string) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, user entities.UpdateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	Dropper
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

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {

	var user entities.User

	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*entities.User, error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, common.InvalidID()
	}

	var user entities.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	var users []*entities.User

	cursor, err := s.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	if users == nil {
		return []*entities.User{}, nil
	}

	return users, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	res, err := s.coll.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, user entities.UpdateUserParams) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": user})

	if err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}

func (s *MongoUserStore) Drop(ctx context.Context) error {

	fmt.Println("----> Dropping users collection <---")
	return s.coll.Drop(ctx)
}
