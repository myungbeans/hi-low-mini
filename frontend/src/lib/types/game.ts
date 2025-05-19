export interface Game {
  id: string;
  cards: Card[];
  createdAt: string;
}

export interface Card {
  id: string;
  type: 'number' | 'operator';
  value: string;
} 