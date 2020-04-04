import * as React from 'react'
import { Data } from 'types'
import Board from 'components/Board'
import './App.css'

const ws = new WebSocket('ws://localhost:8002/ws')

export default function App(): React.ReactElement {
  const [value, setValue] = React.useState<string>('')
  const [data, setData] = React.useState<Data>(undefined)

  React.useEffect(() => {
    ws.onopen = (): void => {}

    ws.onmessage = (msg): void => {
      console.log(JSON.parse(msg.data)) // eslint-disable-line
      setData(JSON.parse(msg.data))
    }
  }, [])

  return (
    <div className="app">
      {
        data
          ? (
            <Board data={data} />
          ) : (
            <>
              <input
                onChange={(e): void => setValue(e.target.value)}
                value={value}
              />
              <button
                type="button"
                onClick={(): void => {
                  ws.send(JSON.stringify({
                    type: 'player',
                    payload: {
                      name: value,
                    },
                  }))
                  setValue('')
                }}
              >
                Send
              </button>
            </>
          )
      }
    </div>
  )
}
