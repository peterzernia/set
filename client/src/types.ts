export type Player = {
  id: number;
  name: string;
  request: boolean;
  score: number;
}

export type Card = {
  color: number;
  shape: number;
  number: number;
  shading: number;
}

export type Data = {
  game_over?: boolean;
  in_play: Card[][];
  players: Player[];
  remaining: number;
}

export type Move = {
  cards: Card[];
}
