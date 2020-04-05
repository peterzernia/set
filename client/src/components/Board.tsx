import * as React from 'react'
import { Data, Card as CardType } from 'types'
import Card from './Card'
import Empty from './Empty'

type Props = {
  data: Data;
  handleMove: (cards: CardType[]) => void;
  handleRequest: () => void;
}

export default function Board(props: Props): React.ReactElement {
  const [selected, setSelected] = React.useState<Card[]>([])
  const { data, handleMove, handleRequest } = props
  const { in_play, players } = data

  const handleClick = (card: Card): void => {
    const i = selected.indexOf(card)
    if (i === -1) {
      const slctd = [...selected, card]
      if (slctd.length === 3) {
        handleMove(slctd)
        setSelected([])
      } else {
        setSelected(slctd)
      }
    } else {
      setSelected(selected.filter((c) => c !== card))
    }
  }

  return (
    <div className="board">
      {in_play.map((cards: CardType[], i: number) => (
        <div className="row" key={i} /* eslint-disable-line */>
          {
            cards.map((card: Card) => {
              if (card.color === null) {
                return <Empty />
              }
              return (
                <Card
                  selected={selected.indexOf(card) !== -1}
                  key={`${card.color}${card.shape}${card.number}${card.shading}`}
                  onClick={(): void => handleClick(card)}
                  card={card}
                />
              )
            })
          }
        </div>
      ))}
      <button
        type="button"
        onClick={(): void => handleRequest()}
      >
        Request more cards
      </button>
      {players.map((player: Player) => (
        <div key={player.id}>
          {`${player.name}: ${player.score} ${player.request ? 'Requested more cards' : ''}`}
        </div>
      ))}
    </div>
  )
}
