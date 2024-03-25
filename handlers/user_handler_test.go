package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	UserStore db.UserStore
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}

}

func setup() *testdb {
	if err := godotenv.Load("../.env"); err != nil {
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

	return &testdb{UserStore: db.NewMongoUserStore(client)}

}
func TestPostUser(t *testing.T) {
	tdb := setup()
	defer tdb.tearDown(t)

	app := fiber.New()

	UserHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", UserHandler.HandlePostUser)

	params := entities.CreateUserParams{
		Email:     "4@5.com",
		Password:  "12345678",
		FirstName: "teste",
		LastName:  "teste",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)

	var user entities.User

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if user.Email != params.Email {
		t.Fatalf("expected %s, got %s", params.Email, user.Email)
	}

	if user.FirstName != params.FirstName {
		t.Fatalf("expected %s, got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Fatalf("expected %s, got %s", params.LastName, user.LastName)

	}

	// fmt.Println(res.Status)

}
