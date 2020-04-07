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

// Context contains necessary clients throughout the application
type Context struct {
	Game *game.Game
}

func main() {
	context := Context{}

	http.Handle("/", http.FileServer(http.Dir("./client/build")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("env") == "development" {
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Websocket error: %s", err)
			return
		}

		// Reset game when there are no connections
		if len(clients) == 0 {
			context.Game = nil
		}

		// Register client
		clients[conn] = true

		fmt.Println(clients)

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

			var res []byte
			switch message.Type {
			case "join":
				context.handleJoin(message, conn)
				res, _ = json.Marshal(context.Game)
			case "move":
				context.handleMove(message, conn)
				res, _ = json.Marshal(context.Game)
			case "request":
				context.handleRequest(conn)
				res, _ = json.Marshal(context.Game)
			case "new":
				context.handleNew()
				res, _ = json.Marshal(context.Game)
			default:
				res = []byte("Unrecognized message type" + message.Type)
				log.Println("Unrecognized message type" + message.Type)
			}

			// Send to every client that is currently connected
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, res)
				if err != nil {
					// Remove player
					for i, v := range context.Game.Players {
						if v.Conn == client {
							fmt.Println(clients)
							fmt.Println(*v.Name)
							context.Game.Players = append(context.Game.Players[:i], context.Game.Players[i+1:]...)
						}
					}

					// Remove connection
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	})

	port := ":" + os.Getenv("PORT")
	fmt.Println("Running server on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (c *Context) handleJoin(message message.Message, conn *websocket.Conn) *int64 {
	if c.Game == nil || c.Game.GameOver != nil && *c.Game.GameOver {
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

	end, err := c.Game.Play(&move, conn)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if end {
		c.Game.GameOver = ptr.Bool(true)
	}

	return nil
}

func (c *Context) handleRequest(conn *websocket.Conn) {
	for i, v := range c.Game.Players {
		if v.Conn == conn {
			c.Game.Players[i].Request = ptr.Bool(true)
		}
	}

	request := true
	for _, v := range c.Game.Players {
		if !*v.Request {
			request = false
		}
	}

	if request && c.Game.Deck.Cards != nil && len(c.Game.Deck.Cards) > 0 {
		c.Game.AddCards()
		for i := range c.Game.Players {
			c.Game.Players[i].Request = ptr.Bool(false)
		}
	}
}

func (c *Context) handleNew() {
	players := c.Game.Players

	for i := range players {
		players[i].Score = ptr.Int64(0)
	}

	c.Game = game.New()
	c.Game.Players = players
}
