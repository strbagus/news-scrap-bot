package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func InitMongo() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set the MONGODB_URI environment variable.")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), map[string]int{"ping": 1}).Err(); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}


	fmt.Println("Connected to MongoDB")

	Client = client

	if err := ensureIndexes(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func ensureIndexes(ctx context.Context) error {
	db := Client.Database("ftibot")

	indexes := []struct {
		col   string
		model mongo.IndexModel
	}{
		{
			"users",
			mongo.IndexModel{
				Keys:    bson.D{{Key: "chat_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			"news",
			mongo.IndexModel{
				Keys:    bson.D{{Key: "title", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
	}

	for _, idx := range indexes {
		if _, err := db.Collection(idx.col).Indexes().CreateOne(ctx, idx.model); err != nil {
			return fmt.Errorf("failed to create %s index: %w", idx.col, err)
		}
	}

	return nil
}
func CloseMongo() {
	if Client != nil {
		if err := Client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		} else {
			fmt.Println("MongoDB connection closed.")
		}
	}
}
