package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gameWithID struct {
	ID          primitive.ObjectID `bson:"_id" json:"gameId"`
	White       string             `json:"white"`
	Black       string             `json:"black"`
	History     []string           `json:"history"`
	GameStarted bool               `bson:"gameStarted" json:"gameStarted"`
}

type game struct {
	White       string
	Black       string
	History     []string
	GameStarted bool
}

type move struct {
	ID   string
	Move string
}

type gameReq struct {
	Username string `json:"username"`
}

type changeEvent struct {
	OperationType string     `bson:"operationType"`
	FullDocument  gameWithID `bson:"fullDocument" json:"fullDocument"`
}

var clients = make(map[string][]*websocket.Conn)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	eventChan := make(chan changeEvent)

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

	go watchForChanges(coll, eventChan)
	//go handleEvents(eventChan)

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))
	http.HandleFunc("/ws", makeHandleConnections(eventChan, coll))
	http.HandleFunc("/start", makeHandleStart(coll))

	log.Fatal(http.ListenAndServe("localhost:9001", nil))
}

// func handleEvents(eventChan chan changeEvent) {
// 	for event := range eventChan {
// 		fmt.Println("Event received", event)
// 		fmt.Println("id", event.FullDocument.ID.Hex())
// 		for _, client := range clients[event.FullDocument.ID.Hex()] {
// 			fmt.Println("client: ", client)
// 			err := client.WriteJSON(event)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				client.Close()
// 			}
// 		}
// 	}
// }

func makeHandleConnections(eventChan chan changeEvent, coll *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {

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
		fmt.Println("WS!!")

		go readMessages(ws, coll)

		for event := range eventChan {
			fmt.Println("Event received", event)
			fmt.Println("id", event.FullDocument.ID.Hex())
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
		fmt.Println(string(message))
		fmt.Println(m)

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
