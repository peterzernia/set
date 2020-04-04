import * as React from 'react'
import Icon from 'assets/oval.svg'

type Props = {
  shading: string;
  color: string;
}

export default function Oval(props: Props): React.ReactElement {
  const { shading, color } = props

  return (
    <Icon className={`${color}-${shading}`} />
  )
}
