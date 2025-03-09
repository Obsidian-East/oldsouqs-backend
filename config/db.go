package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() (*mongo.Database, error) {
	// Set MongoDB API options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI("mongodb+srv://eastobsidian:obsidianjad@clusteroldsouqs.fbnni.mongodb.net/?retryWrites=true&w=majority&appName=ClusterOldSouqs").
		SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Ping to confirm a successful connection
	if err := client.Database("oldsouqs").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client.Database("oldsouqs"), nil
}