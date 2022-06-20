package database

import (
	"context"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client

func Connect() error {
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(uint64(15))
	clientlocal, err := mongo.Connect(context.TODO(), clientOptions)
	client = *clientlocal
	if err != nil {
		return err
	}
	err = clientlocal.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	openlog.Info("Connected to Mongodb")
	return nil
}

func GetClient() *mongo.Client { return &client }
