package config

import (
	"context"
	"fmt"
	"os"

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

	// Get database name from env variable
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "oldsouqs" // Default database name
	}

	fmt.Println("Using Database:", dbName)

	// Ping to confirm a successful connection
	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client.Database(dbName), nil
}
