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
  in_play: Card[][];
  players: Player[];
}

type Move = {
  cards: Card[];
  player_id: number;
}
