package main

import "github.com/peterzernia/set/game"

// Context keeps tracks of the game, registered connections, and the message queue.
type Context struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// The game of set
	Game *game.Game
}

func newContext() *Context {
	return &Context{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (c *Context) run() {
	// Initialize game
	if c.Game == nil {
		c.Game = game.New()
	}

	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)

				// remove player from game
				for i, v := range c.Game.Players {
					if v.Conn == client.conn {
						c.Game.Players = append(c.Game.Players[:i], c.Game.Players[i+1:]...)
					}
				}

				if len(c.Game.Players) == 0 {
					c.Game = game.New()
				}
			}
		case message := <-c.broadcast:
			for client := range c.clients {

				// only send mesages to clients that have
				// joined the game with a username
				isPlaying := false
				for _, player := range c.Game.Players {
					if client.conn == player.Conn {
						isPlaying = true
					}
				}

				if isPlaying {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(c.clients, client)

						// remove player from game
						for i, v := range c.Game.Players {
							if v.Conn == client.conn {
								c.Game.Players = append(c.Game.Players[:i], c.Game.Players[i+1:]...)
							}
						}

						if len(c.Game.Players) == 0 {
							c.Game = game.New()
						}
					}
				}
			}
		}
	}
}
