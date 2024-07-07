package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
)

type ServerSetting struct {
	GuildID  string          `bson:"guild_id"`
	Settings map[string]bool `bson:"settings"`
}

func initDatabase() {
	db := os.Getenv("MONGO_DB_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(db)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("discord_bot").Collection("server_settings")
	log.Println("Connected to MongoDB")
}

func getServerSettingFromDB(guildID, service string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result ServerSetting
	err := collection.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
		log.Printf("Error fetching server settings: %v", err)
		return true
	}

	return !result.Settings[service]
}

func setServerSettingToDB(guildID, service string, enabled bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$set": bson.M{"settings." + service: !enabled}}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error updating server settings: %v", err)
	}
}
