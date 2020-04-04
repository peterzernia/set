import * as React from 'react'
import { Card as CardType } from 'types'
import { COLORS, SHAPES, SHADINGS } from 'constants'
import Diamond from './Diamond'
import Oval from './Oval'
import Squiggle from './Squiggle'

type Props = {
  card: CardType;
}

export default function Card(props: Props): React.ReactElement {
  const { card } = props
  const {
    color,
    shape,
    number,
    shading,
  } = card

  let element: React.ReactElement
  switch (shape) {
    case SHAPES.DIAMOND:
      element = <Diamond color={COLORS[color]} shading={SHADINGS[shading]} />
      break
    case SHAPES.OVAL:
      element = <Oval color={COLORS[color]} shading={SHADINGS[shading]} />
      break
    case SHAPES.SQUIGGLE:
      element = <Squiggle color={COLORS[color]} shading={SHADINGS[shading]} />
      break
    default:
      throw new Error('Undefined shape')
  }
  const elements = [...Array(number + 1).keys()].map(() => element)
  return (
    <div className="card">
      {elements}
    </div>
  )
}
