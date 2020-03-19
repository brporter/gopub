package storage

import (
	"context"
	"time"

	"github.com/brporter/gopub/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPostRepo struct {
	client     *mongo.Client
	database   *string
	collection *string
	isOpen     bool
}

func NewMongoRepo(open bool, configuration models.IConfiguration) (IPostRepo, error) {
	retVal := new(MongoPostRepo)

	client, err := getClient(configuration)

	if err != nil {
		return nil, err
	}

	retVal.client = client

	retVal.database, err = configuration.GetSecret("database")

	if err != nil {
		return nil, err
	}

	retVal.collection, err = configuration.GetSecret("collection")

	if err != nil {
		return nil, err
	}

	if open {
		err = retVal.Open()
	}

	return retVal, err
}

func getClient(configuration models.IConfiguration) (*mongo.Client, error) {
	connectionString, err := configuration.GetSecret("storage")

	if err != nil {
		panic(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(*connectionString))

	return client, err
}

func makeContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout*time.Second)
}

func (r *MongoPostRepo) Open() error {
	if !r.isOpen {
		ctx, cancel := makeContext(5 * time.Second)
		defer cancel()

		return r.client.Connect(ctx)
	} else {
		return nil
	}
}

func (r *MongoPostRepo) Close() error {
	ctx, cancel := makeContext(5 * time.Second)
	defer cancel()

	return r.client.Disconnect(ctx)
}

func (r *MongoPostRepo) Save(p *models.Post) error {
	ctx, cancel := makeContext(5 * time.Second)
	defer cancel()

	var opts options.ReplaceOptions
	opts.SetUpsert(true)

	collection := r.client.Database(*r.database).Collection(*r.collection)
	filter := bson.M{"_id": primitive.Binary{Subtype: 0x00, Data: (*p).PostId[:]}}
	_, err := collection.ReplaceOne(ctx, filter, *p, &opts)

	return err
}

func (r *MongoPostRepo) FetchOne(id *uuid.UUID) (*models.Post, error) {
	ctx, cancel := makeContext(5 * time.Second)
	defer cancel()

	var result models.Post
	collection := r.client.Database(*r.database).Collection(*r.collection)
	filter := bson.M{"_id": primitive.Binary{Subtype: 0x00, Data: (*id)[:]}}

	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func (r *MongoPostRepo) FetchMany(publishDate time.Time, pageSize int) ([]*models.Post, error) {
	ctx, cancel := makeContext(5 * time.Second)
	defer cancel()

	collection := r.client.Database(*r.database).Collection(*r.collection)
	filter := bson.M{"publishDate": bson.M{"$lt": publishDate.UTC()}}
	projection := bson.D{{"body", 0}} // don't fetch body fields for multi-doc fetches

	fopts := options.Find()
	fopts.SetLimit(int64(pageSize))
	fopts.SetProjection(projection)

	results, err := collection.Find(ctx, filter, fopts)

	if err != nil {
		return nil, err
	}

	returnValue := make([]*models.Post, 0, pageSize) // type, len, cap
	var returnError error

	for results.Next(ctx) {
		var p models.Post
		err = results.Decode(&p)

		if err != nil {
			returnError = err
		}

		returnValue = append(returnValue, &p)
	}

	return returnValue, returnError
}

func (r *MongoPostRepo) Remove(id *uuid.UUID) error {
	ctx, cancel := makeContext(5 * time.Second)
	defer cancel()

	collection := r.client.Database(*r.database).Collection(*r.collection)
	filter := bson.M{"_id": primitive.Binary{Subtype: 0x00, Data: (*id)[:]}}

	_, err := collection.DeleteOne(ctx, filter, nil)

	if err != nil {
		return err
	}

	return nil
}
