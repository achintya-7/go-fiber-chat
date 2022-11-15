package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetEnv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

func GetJWTSecret() string {
	return GetEnv("SECRET_TOKEN")
}

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := GetEnv("DB_NAME")
	if dbName == "" {
		log.Fatal("No database name found in env file")
	}
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}
