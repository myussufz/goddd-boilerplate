package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"goddd-boilerplate/app/model"
	"goddd-boilerplate/app/repository/kit"
)

// User :
type User struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewUser :
func NewUser(client *mongo.Client, db *mongo.Database) User {
	return User{client: client, db: db}
}

// Create :
func (r User) Create(user *model.User) error {
	_, err := r.db.Collection(model.CollectionUser).InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// FindByPhoneNumber :
func (r User) FindByPhoneNumber(countryCode, phoneNumber string) (*model.User, error) {
	user := new(model.User)

	if err := r.db.Collection(model.CollectionUser).FindOne(
		context.Background(),
		bson.M{
			"country_code": countryCode,
			"phone_number": phoneNumber,
		},
	).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UniqueIndex :
func (r User) UniqueIndex(key, value string) error {
	index := mongo.IndexModel{}
	index.Keys = bsonx.Doc{{Key: key, Value: bsonx.String(value)}}
	index.Options = &options.IndexOptions{}
	index.Options.SetUnique(true)

	_, err := r.db.Collection(model.CollectionUser).Indexes().
		CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	return nil
}

// FindByID :
func (r User) FindByID(id string) (*model.User, error) {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	if err := r.db.Collection(model.CollectionUser).FindOne(
		context.Background(),
		bson.M{"_id": hexID},
	).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Paginate :
func (r User) Paginate(p kit.Paginate) ([]*model.User, string, error) {
	users := make([]*model.User, 0)

	ctx := context.Background()

	query := bson.M{}

	if p.Cursor != "" {
		objectID, err := primitive.ObjectIDFromHex(p.Cursor)
		if err != nil {
			return nil, "", err
		}
		query["_id"] = bson.M{"$gte": objectID}
	}

	for _, each := range p.FilterColumns {
		if val, isExist := p.Filters[each]; isExist {
			query[each] = val
		}
	}

	nextCursor, err := r.db.Collection(model.CollectionUser).Find(
		ctx,
		query,
		options.Find().SetLimit(p.Limit+1),
	)

	defer nextCursor.Close(ctx)
	if err != nil {
		return nil, "", err
	}

	for nextCursor.Next(ctx) {
		user := new(model.User)
		if err := nextCursor.Decode(user); err != nil {
			return nil, "", err
		}
		users = append(users, user)
	}

	if err := nextCursor.Err(); err != nil {
		return nil, "", err
	}

	if len(users) > int(p.Limit) {
		return users[:len(users)-1], users[len(users)-1].ID.Hex(), nil
	}

	return users, fmt.Sprintf("%d", nextCursor.ID()), nil
}
