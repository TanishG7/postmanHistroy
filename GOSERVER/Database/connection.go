package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connection *mongo.Client

func connectDatabase() error {
	connectionString := os.Getenv("DB_URL")
	fmt.Println("connection String:", connectionString)

	clientOptions := options.Client().ApplyURI(connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		fmt.Println("err: ", err.Error())
	}
	connection = client
	return err
}

func GetDatabaseConnection() (*mongo.Client, error) {
	var connectError error
	if connection != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := connection.Ping(ctx, readpref.Primary())

		if err != nil {
			fmt.Println("Fetching new connection")
			connectError = connectDatabase()
		}
	} else {
		connectError = connectDatabase()
	}
	return connection, connectError
}

func CloseConnection() {
	if connection != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := connection.Disconnect(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}
