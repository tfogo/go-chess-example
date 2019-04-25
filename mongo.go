package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo() *mongo.Collection {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("CHESS_MONGO_URI"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	coll := client.Database("chess").Collection("games")

	return coll
}

func watchForChanges(coll *mongo.Collection, eventChan chan changeEvent) {
	p := mongo.Pipeline{}
	csoptions := options.ChangeStream()
	csoptions.SetFullDocument(options.UpdateLookup)

	cs, err := coll.Watch(context.TODO(), p, csoptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	for cs.Next(context.TODO()) {
		var res changeEvent

		fmt.Printf("%v\n", cs.Current)

		err = cs.Decode(&res)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v\n", res)
		eventChan <- res
	}

}
