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

// Context contains necessary clients throught the application
type Context struct {
	Game *game.Game
}

func main() {
	context := Context{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			message := message.Message{}
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			err = json.Unmarshal(msg, &message)

			if err != nil {
				fmt.Println(err.Error())
			}

			var res []byte
			switch message.Type {
			case "player":
				context.handlePlayer(message)
				res, _ = json.Marshal(context.Game)
			case "move":
				context.handleMove(message)
				res, _ = json.Marshal(context.Game)
			default:
				res = []byte("Unrecognized message type" + message.Type)
				fmt.Println("Unrecognized message type" + message.Type)
			}

			err = conn.WriteMessage(msgType, res)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	port := ":" + os.Getenv("PORT")
	fmt.Println("Running server on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (c *Context) handlePlayer(message message.Message) {
	if c.Game == nil {
		c.Game = game.New()
	}

	player := game.Player{}
	player.ID = ptr.Int64(int64(len(c.Game.Players)) + 1)
	player.Request = ptr.Bool(false)
	player.Score = ptr.Int64(0)

	if name, ok := message.Payload["name"].(string); ok {
		player.Name = &name
	}

	c.Game.Players = append(c.Game.Players, player)
	return
}

func (c *Context) handleMove(message message.Message) error {
	move := game.Move{}
	j, _ := json.Marshal(message.Payload)
	json.Unmarshal(j, &move)

	err := c.Game.Play(&move)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
