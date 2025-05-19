import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import type { Game } from '$lib/types/game';

interface CachedGame {
  game: Game;
  timestamp: number;
}

const CACHE_KEY = 'cachedGame';
const ONE_DAY_MS = 24 * 60 * 60 * 1000;

// Initialize from localStorage if available
const storedGame = browser ? localStorage.getItem(CACHE_KEY) : null;
const initialGame: CachedGame | null = storedGame ? JSON.parse(storedGame) : null;

// Check if the cached game is from today
const isGameFromToday = (cachedGame: CachedGame | null): boolean => {
  if (!cachedGame) return false;
  const now = new Date();
  const gameDate = new Date(cachedGame.timestamp);
  return (
    now.getFullYear() === gameDate.getFullYear() &&
    now.getMonth() === gameDate.getMonth() &&
    now.getDate() === gameDate.getDate()
  );
};

// Create the store with initial value
export const gameStore = writable<Game | null>(isGameFromToday(initialGame) ? initialGame.game : null);

// Derived store to check if we need to fetch a new game
export const needsNewGame = derived(gameStore, ($game) => !$game);

// Function to update the game store and cache
export const updateGame = (game: Game) => {
  const cachedGame: CachedGame = {
    game,
    timestamp: Date.now()
  };
  
  if (browser) {
    localStorage.setItem(CACHE_KEY, JSON.stringify(cachedGame));
  }
  
  gameStore.set(game);
};

// Function to clear the game cache
export const clearGameCache = () => {
  if (browser) {
    localStorage.removeItem(CACHE_KEY);
  }
  gameStore.set(null);
}; 