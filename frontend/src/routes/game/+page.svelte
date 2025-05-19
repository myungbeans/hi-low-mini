<script lang="ts">
  import { onMount } from 'svelte';
  import { gameStore, needsNewGame, updateGame } from '$lib/stores/game';
  import type { Game } from '$lib/types/game';

  let loading = true;
  let error: string | null = null;

  async function fetchGame() {
    try {
      loading = true;
      error = null;
      const response = await fetch('/api/game');
      if (!response.ok) {
        throw new Error('Failed to fetch game');
      }
      const game: Game = await response.json();
      updateGame(game);
    } catch (e) {
      error = e instanceof Error ? e.message : 'An error occurred';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    if ($needsNewGame) {
      fetchGame();
    } else {
      loading = false;
    }
  });
</script>

{#if loading}
  <div class="flex items-center justify-center min-h-screen">
    <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
  </div>
{:else if error}
  <div class="flex items-center justify-center min-h-screen">
    <div class="text-red-500">Error: {error}</div>
  </div>
{:else if $gameStore}
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6 dark:text-white">Today's Game</h1>
    <div class="grid grid-cols-2 gap-4">
      {#each $gameStore.cards as card}
        <div class="p-4 bg-white dark:bg-gray-800 rounded-lg shadow dark:text-white">
          <div class="text-lg font-semibold">{card.value}</div>
          <div class="text-sm text-gray-500 dark:text-gray-400">{card.type}</div>
        </div>
      {/each}
    </div>
  </div>
{/if} 