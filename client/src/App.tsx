import * as React from 'react'
import { Data, Move, Player } from 'types'
import Board from 'components/Board'
import Join from 'components/Join'
import './App.css'

const ws = new WebSocket('ws://localhost:8002/ws')

export default function App(): React.ReactElement {
  const [data, setData] = React.useState<Data>(undefined)

  React.useEffect(() => {
    ws.onopen = (): void => {}

    ws.onmessage = (msg): void => {
      console.log(JSON.parse(msg.data)) // eslint-disable-line
      setData(JSON.parse(msg.data))
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
        <div>Game Over</div>
        <div>{`Winner: ${winners.join('& ')}`}</div>
        <button type="button" onClick={(): void => handleNew()}>Play again?</button>
      </div>
    )
  }

  return (
    <div className="app">
      {
        data
          ? <Board data={data} handleMove={handleMove} />
          : <Join handleJoin={handleJoin} />
      }
    </div>
  )
}
