package mongoconn

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoClient    *mongo.Client
	postCollection *mongo.Collection
)

func Connect() (*mongo.Database, error) {
	ctx := context.TODO()

	dbUri := "mongodb://localhost:27017"
	connectionOpts := options.Client().ApplyURI(dbUri)
	mongoClient, err := mongo.Connect(ctx, connectionOpts)
	if err != nil {
		fmt.Printf("an error ocurred when connect to mongoDB : %v", err)
		panic(err)
	}

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("an error ocurred when connect to mongoDB : %v", err)
		panic(err)
	}

	return mongoClient.Database("kreditPlus"), nil
}
