package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gameWithID struct {
	ID          primitive.ObjectID `bson:"_id" json:"gameId"`
	White       string             `bson:"white" json:"white"`
	Black       string             `bson:"black" json:"black"`
	History     []string           `bson:"history" json:"history"`
	GameStarted bool               `bson:"gameStarted" json:"gameStarted"`
}

type game struct {
	White       string
	Black       string
	History     []string
	GameStarted bool
}

func makeHandleStart(coll *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			//handle preflight in here
		} else {

			decoder := json.NewDecoder(r.Body)

			var t gameReq
			err := decoder.Decode(&t)

			if err != nil {
				panic(err)
			}

			fmt.Println(t.Username)

			filter := bson.D{{
				"black", "",
			}}

			update := bson.D{{
				"$set", bson.D{{
					"black", t.Username,
				}, {
					"gameStarted", true,
				}},
			}}

			foptions := options.FindOneAndUpdate()
			foptions.SetReturnDocument(options.After)

			var result gameWithID

			err = coll.FindOneAndUpdate(context.TODO(), filter, update, foptions).Decode(&result)
			if err == mongo.ErrNoDocuments {
				newGame := game{
					White:   t.Username,
					History: make([]string, 0),
				}

				insertResult, err := coll.InsertOne(context.TODO(), newGame)
				if err != nil {
					log.Fatal(err)
				}
				newGameWithID := gameWithID{
					White: t.Username,
					ID:    insertResult.InsertedID.(primitive.ObjectID),
				}

				fmt.Println("Inserted a single document: ", insertResult.InsertedID)

				json.NewEncoder(w).Encode(newGameWithID)
			} else if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("Found doc %#v\n", result)

				json.NewEncoder(w).Encode(result)
			}

		}

	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}
