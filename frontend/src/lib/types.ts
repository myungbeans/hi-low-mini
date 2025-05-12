import type { Game as GameProto, Hand as HandProto, PlayHandResponse as PlayHandResponseProto, Card as CardProto } from '$lib/gen/game_engine/v1/game_engine_pb';

export interface Card {
  id: string;  // Unique identifier for the card
  value: number | string;  // number for numeric cards, string for operators
  type: 'number' | 'operator';
}

export type Game = GameProto;
export type Hand = HandProto;
export type PlayHandResponse = PlayHandResponseProto & {
  error?: string;  // Add optional error field
};

export interface GameState {
  game: Game | null;
  selectedCards: Card[];
  availableCards: Card[];
  isPlaying: boolean;
  currentScore: number;
  scoreHistory: number[];
  highlightIndex: number;
} 