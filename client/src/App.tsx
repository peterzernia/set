import * as React from 'react'
import { Data } from 'types'
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

  return (
    <div className="app">
      {
        data
          ? <Board data={data} />
          : <Join handleJoin={handleJoin} />
      }
    </div>
  )
}
