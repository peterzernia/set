import * as React from 'react'
import { Data, Card as CardType } from 'types'
import Card from './Card'

type Props = {
  data: Data;
}

export default function Board(props: Props): React.ReactElement {
  const { data } = props
  const { in_play, players } = data

  return (
    <div className="board">
      {in_play.map((cards: CardType[], i: number) => (
        <div className="row" key={i} /* eslint-disable-line */>
          {
            cards.map((card) => (
              <Card
                key={`${card.color}${card.shape}${card.number}${card.shading}`}
                card={card}
              />
            ))
          }
        </div>
      ))}
      {players.map((player: Player) => (
        <div key={player.id}>
          {`${player.name}: ${player.score}`}
        </div>
      ))}
    </div>
  )
}
