package repository

import (
	"goddd-boilerplate/app/repository/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

// Paginate :
type Paginate struct {
	Cursor        string
	FilterColumns []string
	Filters       map[string]map[string]interface{}
	Limit         int64
}

// Repository :
type Repository struct {
	User User
}

// NewMongoDB :
func NewMongoDB(client *mongo.Client, database string) *Repository {
	return &Repository{
		User: mongodb.NewUser(client, client.Database(database)),
	}
}
