# Realtime Chess Website Example

This is very simple chess website built using React, Go, and MongoDB. Players will be matched with another player when one is available. They will then be able to play chess in realtime. The frontend is built with React and Redux. It uses the chess.js and chessboardjsx libraries which are used by the likes of chess.com and lichess.org. The underlying databases is MongoDB, using the new MongoDB Driver for Go. Clients communicate with the server through websockets. When a player makes a move, the game is updated in the database. A change stream takes any changes to games and sends them to clients via websockets.

![Two browsers side by side. The browsers show chess boards. Moving a piece in one board, moves the piece in the other so two people can play together.](./example.gif)

## Running the application

1. Running the React application:

    ```
    cd frontend
    npm i 
    npm run dev
    ```

    This runs webpack dev server. If you build the frontend and serve it with the go app, the chess board will flicker when moving pieces. I haven't figured out the cause of this yet, so best results are when using `npm run dev` to run the react app separatey from the go app. 

2. Running the Go backend:

    ```
    cd ..
    go run .
    ```

## Displaying the application in use

1. Open `localhost:8080` in two browser windows side by side.
2. Enter a name in one window and click "Start Game".
3. Enter another name in the other window and click "Start Game".
4. Each window will now have a chess board open. White moves first.
5. Each move made will be broadcast to the other player. A list of moves in the game is shown on the side.
6. Quick gaem to show off (Scholar's Mate): `1. e4 e5 2. Bc4 Nc6 3. Qh5 Nf6?? 4. Qxf7#`

## How does it work?

There are 4 files:

- `main.go`: `main()` does three things:
    1. calls a function to connect to MongoDB, 
    2. calls a function to start a change stream to watch the `chess.games` collection,
    3. Defines three routes and starts an http server.

    The server handles two routes:

    1. `/ws`: The websocket which the frontend connects to. The frontend will send moves made to this websocket.
    2. `/start`: Clients will POST here to start a game. 

- `mongo.go`:
    - `connectMongo()`: Connects to MongoDB using `mongo.Connect()`.
    - `watchForChanges()`: Will watch a MongoDB collection and will send the change events to a channel.

- `websockets.go`
    - `makeHandleWebsockets()`: Returns the handler function for the `/ws` endpoint. It creates a websocket using the `gorilla` websocket library. A couple of things happen here:

        1. The client registers with the websocket.
        2. Then in the `readMessages` goroutine, we wait for any messages from the client. The client sends moves. These are added to the game document in the database.
        3. At the same time, we wait for any change events sent through a channel. When we receive these change events, we broadcast that to the players so the state of their game is updated.


- `find_game.go`
    - `makeHandleStart()`: Returns the handler function for the `/start` endpoint. This does the matchmaking. The matchmaking is very simple. If there's a document in the games collection where `black` is an empty string, that means there's no black player, and the player is added to that game as black. If there isn't a document where `black` is an empty string, a new game is made where the player is set to be white. They will then have to wait for another player to come along and join their game as black.

## Demoing

It's worth just showing off a couple of things:

1. Show the functionality of the chess app, by playing a few moves.
2. Show `mongo.go`. This shows how to connect to the database in `connectMongo()` and how it watches for changes in `watchForChanges()`
3. You can then show how in `websockets.go`, and changes to a game is sent to the players via websockets. The key code block is:

    ```go
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
    ```