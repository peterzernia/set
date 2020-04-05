package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/set/game"
	"github.com/peterzernia/set/message"
	"github.com/peterzernia/set/ptr"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Keeps track of all connections
var clients = make(map[*websocket.Conn]bool)

// Context contains necessary clients throught the application
type Context struct {
	Game *game.Game
}

func main() {
	context := Context{}

	http.Handle("/", http.FileServer(http.Dir("./client/build")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Websocket error: %s", err)
			return
		}

		// Register client
		clients[conn] = true

		for {
			message := message.Message{}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Websocket error: %s", err)
				return
			}

			err = json.Unmarshal(msg, &message)

			if err != nil {
				log.Printf("Websocket error: %s", err)
			}

			// Reset game when there are no connections
			if len(clients) == 0 {
				context.Game = nil
			}

			var res []byte
			switch message.Type {
			case "join":
				context.handleJoin(message, conn)
				res, _ = json.Marshal(context.Game)
			case "move":
				context.handleMove(message, conn)
				res, _ = json.Marshal(context.Game)
			default:
				res = []byte("Unrecognized message type" + message.Type)
				log.Println("Unrecognized message type" + message.Type)
			}

			// Send to every client that is currently connected
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, res)
				if err != nil {
					// Remove connection
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(clients, client)

					// Remove player
					var index int
					for i, v := range context.Game.Players {
						if v.Conn == conn {
							index = i
						}
					}

					context.Game.Players = append(context.Game.Players[:index], context.Game.Players[index+1:]...)
				}
			}
		}
	})

	port := ":" + os.Getenv("PORT")
	fmt.Println("Running server on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (c *Context) handleJoin(message message.Message, conn *websocket.Conn) *int64 {
	if c.Game == nil {
		c.Game = game.New()
	}

	player := game.Player{}
	player.ID = ptr.Int64(int64(len(c.Game.Players)) + 1)
	player.Conn = conn
	player.Request = ptr.Bool(false)
	player.Score = ptr.Int64(0)

	if name, ok := message.Payload["name"].(string); ok {
		player.Name = &name
	}

	c.Game.Players = append(c.Game.Players, player)
	return player.ID
}

func (c *Context) handleMove(message message.Message, conn *websocket.Conn) error {
	move := game.Move{}
	j, _ := json.Marshal(message.Payload)
	json.Unmarshal(j, &move)

	err := c.Game.Play(&move, conn)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
