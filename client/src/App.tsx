import * as React from 'react'
import {
  Card,
  Data,
  Move,
  Player,
} from 'types'
import Board from 'components/Board'
import Join from 'components/Join'
import GameOver from 'assets/game_over.gif'
import './App.css'


const ws = new WebSocket(`${process.env.REACT_APP_API_URL}/ws`)

export default function App(): React.ReactElement {
  const [data, setData] = React.useState<Data | undefined>()

  React.useEffect(() => {
    ws.onopen = (): void => {}

    ws.onmessage = (msg): void => {
      setData(JSON.parse(msg.data))
    }

    ws.onclose = (msg): void => {
      alert('Disconnected from server') // eslint-disable-line
      window.location.reload()
    }
  }, [])

  const handleJoin = (name: string): void => {
    ws.send(JSON.stringify({
      type: 'join',
      payload: {
        name,
      },
    }))
  }

  const handleMove = (cards: Card[]): void => {
    const move: Move = { cards }
    ws.send(JSON.stringify({
      type: 'move',
      payload: move,
    }))
  }

  const handleRequest = (): void => {
    ws.send(JSON.stringify({
      type: 'request',
    }))
  }

  const handleNew = (): void => {
    ws.send(JSON.stringify({
      type: 'new',
    }))
  }

  if (data && data.game_over) {
    const highScore = Math.max(...data.players.map((p: Player) => p.score))
    const winners = data.players.filter(
      (p: Player) => p.score === highScore,
    ).map((p: Player) => p.name)

    return (
      <div className="app">
        <div className="game-over">
          <img src={GameOver} alt="game-over" />
          <div>{`Winner: ${winners.join('& ')}`}</div>
          <button type="button" onClick={(): void => handleNew()}>Play again?</button>
        </div>
      </div>
    )
  }

  return (
    <div className="app">
      { data
        ? <Board data={data} handleMove={handleMove} handleRequest={handleRequest} />
        : <Join handleJoin={handleJoin} /> }
    </div>
  )
}
