package router

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"goddd-boilerplate/app/config"
	"goddd-boilerplate/app/repository"
)

// New :
func New(e *echo.Echo) *echo.Echo {
	e = versionOne(e, repository.NewMongoDB(getMongoClient()))

	return e
}

func getMongoClient() (*mongo.Client, string) {
	connStr := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s",
		config.MongoDBUsername,
		config.MongoDBPassword,
		config.MongoDBHost,
		config.MongoDBName,
	)
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	os.Stdout.WriteString("Connected to MongoDB!")

	return client, config.MongoDBName
}
