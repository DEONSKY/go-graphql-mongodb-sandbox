package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DEONSKY/go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

const uri = "mongodb://root:running-away-222@localhost:27017/?maxPoolSize=20&w=majority"

func Connect() *DB {

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	return &DB{
		client: client,
	}

}

func (db *DB) Save(input *model.NewTodo) *model.Todo {
	collection := db.client.Database("newbie").Collection("todos")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
	}

	return &model.Todo{
		ID:   res.InsertedID.(primitive.ObjectID).Hex(),
		Text: input.Text,
	}
}

func (db *DB) SaveUser(input *model.NewUser) *model.User {
	collection := db.client.Database("newbie").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
	}

	return &model.User{
		ID:   res.InsertedID.(primitive.ObjectID).Hex(),
		Name: input.Name,
	}
}

func (db *DB) All() []*model.Todo {
	collection := db.client.Database("newbie").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var todos []*model.Todo
	for cur.Next(ctx) {
		var todo *model.Todo
		err := cur.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	return todos
}

func (db *DB) AllUsers() []*model.User {
	collection := db.client.Database("newbie").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var users []*model.User
	for cur.Next(ctx) {
		var user *model.User

		err := cur.Decode(&user)

		fmt.Printf("user.ID: %v\n", user.ID)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}
