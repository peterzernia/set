type Player = {
  id: number;
  name: string;
  request: boolean;
  score: number;
}

type Card = {
  color: number;
  shape: number;
  number: number;
  shading: number;
}

type Data = {
  game_over?: boolean;
  in_play: Card[][];
  players: Player[];
  remaining: number;
}

type Move = {
  cards: Card[];
}
