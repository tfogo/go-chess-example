package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func makeHandleWebsockets(eventChan chan changeEvent, coll *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		// Make sure we close the connection when the function returns
		defer ws.Close()

		// Register our new client
		id := r.URL.Query().Get("id")
		fmt.Println(id)
		clients[id] = append(clients[id], ws)
		fmt.Printf("New client registered for game ID %v\n", id)

		go readMessages(ws, coll)

		for event := range eventChan {
			fmt.Println("Event received: ", event)
			for _, client := range clients[event.FullDocument.ID.Hex()] {
				err := client.WriteJSON(event.FullDocument)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
				}
			}
		}
	}

}

func readMessages(conn *websocket.Conn, coll *mongo.Collection) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		var m move
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Fatal(err)
		}

		id, err := primitive.ObjectIDFromHex(m.ID)
		if err != nil {
			log.Fatal(err)
		}

		filter := bson.D{{
			"_id", id,
		}}

		update := bson.D{{
			"$push", bson.D{{
				"history", m.Move,
			}},
		}}

		updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	}
}
