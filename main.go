package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/set/deck"
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
	context.Game = game.New()

	context.Game.Refresh()

	cards := context.Game.InPlay[0][0:3]
	move := game.Move{Cards: cards}
	err := context.Game.Play(move)
	if err != nil {
		fmt.Println(err)
	}
	context.Game.Refresh()
	context.Game.AddCards()

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
				player := context.handlePlayer(message)
				context.Game.Players = append(context.Game.Players, player)
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
	http.ListenAndServe(":3000", nil)
}

func (c Context) handlePlayer(message message.Message) game.Player {
	player := game.Player{}
	player.ID = ptr.Int64(int64(len(c.Game.Players)) + 1)
	player.Request = ptr.Bool(false)
	player.Score = ptr.Int64(0)

	if name, ok := message.Payload["name"].(string); ok {
		player.Name = &name
	}

	c.Game.Players = append(c.Game.Players, player)
	return player
}

func (c *Context) handleMove(message message.Message) error {
	if cards, ok := message.Payload["cards"].([]deck.Card); ok {
		valid, err := c.Game.Deck.CheckSet(cards)
		if !valid || err != nil {
			return errors.New("Invalid set")
		}
	}
	return nil
}
