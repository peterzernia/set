package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/set/game"
	"github.com/peterzernia/set/message"
	"github.com/peterzernia/set/ptr"
)

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

	// Reset request for new cards when a set is successful
	for i := range c.Game.Players {
		c.Game.Players[i].Request = ptr.Bool(false)
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
