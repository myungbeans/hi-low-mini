<script lang="ts">
  import { onMount } from 'svelte';
  import Card from '$lib/components/Card.svelte';
  import Scoreboard from '$lib/components/Scoreboard.svelte';
  import type { Card as CardType, Game, PlayHandResponse } from '$lib/types';
  import { HandSchema, CardType as ProtoCardType } from '$lib/gen/game_engine/v1/game_engine_pb';
  import { create } from '@bufbuild/protobuf';

  let game: Game | null = null;
  let selectedCards: CardType[] = [];
  let availableCards: CardType[] = [];
  let isPlaying = false;
  let currentScore = 0;
  let scoreHistory: number[] = [];
  let highlightIndex = -1;
  let error: string | null = null;

  onMount(async () => {
    await fetchGame();
  });

  async function fetchGame() {
    try {
      error = null;
      const response = await fetch('/api/game');
      if (!response.ok) {
        throw new Error('Failed to fetch game');
      }
      const data = await response.json();
      
      if (data.error) {
        throw new Error(data.error);
      }

      game = data.game;
      
      // Convert protobuf cards to our CardType format
      availableCards = (game?.pool?.cards || []).map((card, index) => {
        // Check if it's an operator card (type 2)
        const isOperator = card.type === 2; // CARD_TYPE_OPERATOR = 2
        
        return {
          id: `${index}-${card.value}-${isOperator ? 'op' : 'num'}`,
          value: isOperator ? card.value : Number(card.value), // Convert to number only for numeric cards
          type: isOperator ? 'operator' : 'number'
        };
      });

      shuffleCards();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to fetch game';
    }
  }

  function shuffleCards() {
    // Create new card objects while shuffling to maintain unique IDs
    availableCards = availableCards
      .map(card => ({ 
        card: {
          id: card.id,
          value: card.value,
          type: card.type
        }, 
        sort: Math.random() 
      }))
      .sort((a, b) => a.sort - b.sort)
      .map(({ card }) => card);
  }

  function handleCardClick(card: CardType) {
    if (isPlaying) return;

    const isSelected = selectedCards.some(c => c.id === card.id);

    if (isSelected) {
      // Return the card to available cards
      selectedCards = selectedCards.filter(c => c.id !== card.id);
      // Create a new card object to avoid reference issues
      const newCard: CardType = {
        id: card.id,
        value: card.value,
        type: card.type
      };
      availableCards = [...availableCards, newCard];
    } else if (selectedCards.length < 7) {
      // Add the card to selected cards
      const newCard: CardType = {
        id: card.id,
        value: card.value,
        type: card.type
      };
      selectedCards = [...selectedCards, newCard];
      // Remove exactly one instance of the card from available cards
      const cardIndex = availableCards.findIndex(c => c.id === card.id);
      availableCards = [
        ...availableCards.slice(0, cardIndex),
        ...availableCards.slice(cardIndex + 1)
      ];
    }
  }

  async function playHand() {
    if (selectedCards.length !== 7 || isPlaying || !game) {
      return;
    }

    isPlaying = true;
    highlightIndex = 0;

    try {
      error = null;
      
      // Validate selected cards
      if (!selectedCards.every(card => card && typeof card.value !== 'undefined' && card.type)) {
        throw new Error('Invalid card data in selected cards');
      }
      
      // Create the hand data with proper card types
      const handCards = selectedCards.map(card => ({
        value: card.value.toString(), // Ensure value is a string
        type: card.type === 'operator' ? 2 : 1  // Use numeric values: 2 for CARD_TYPE_OPERATOR, 1 for CARD_TYPE_NUMBER
      }));

      const requestData = {
        "@type": "type.googleapis.com/game_engine.v1.PlayHandRequest",
        "value": {
          game_id: game.id,
          timestamp: {
            seconds: Math.floor(Date.now() / 1000).toString(),
            nanos: 0
          },
          hand: {
            cards: handCards,
            game_id: game.id
          },
          secs_elapsed: 0
        }
      };

      const response = await fetch('/api/play', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to play hand: ${errorText}`);
      }

      const data: PlayHandResponse = await response.json();
      
      if (data.error) {
        throw new Error(data.error);
      }
      
      // Animate through the cards and update score
      for (let i = 0; i < data.scoreCounter.length; i++) {
        await new Promise(resolve => setTimeout(resolve, 1000));
        highlightIndex = i;
        currentScore = data.scoreCounter[i];
      }

      scoreHistory = [...scoreHistory, currentScore];
      
      // Reset for next play
      await new Promise(resolve => setTimeout(resolve, 1000));
      highlightIndex = -1;
      selectedCards = [];
      isPlaying = false;
      
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to play hand';
      isPlaying = false;
      highlightIndex = -1;
    }
  }
</script>

<div class="min-h-screen bg-background p-4">
  <div class="max-w-lg mx-auto space-y-8">
    <h1 class="text-4xl font-bold text-center text-text">Hi-Lo Mini</h1>
    
    {#if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{error}</span>
      </div>
    {/if}

    <Scoreboard {currentScore} {scoreHistory} />
    
    <div class="space-y-6">
      <div class="bg-white rounded-lg p-4 shadow-lg">
        <h2 class="text-xl font-semibold mb-4">Selected Cards ({selectedCards.length}/7)</h2>
        <div class="flex flex-wrap gap-2 justify-center">
          {#each selectedCards as card, i}
            <Card
              {card}
              isSelected={true}
              isHighlighted={highlightIndex === i}
              onClick={() => handleCardClick(card)}
            />
          {/each}
          {#each Array(7 - selectedCards.length) as _}
            <div class="w-16 h-24 border-2 border-dashed border-gray-300 rounded-lg" />
          {/each}
        </div>
      </div>

      <div class="bg-white rounded-lg p-4 shadow-lg">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">Available Cards ({availableCards.length})</h2>
          <button
            class="px-4 py-2 bg-secondary text-white rounded-lg hover:bg-opacity-90 transition-colors"
            on:click={shuffleCards}
          >
            Shuffle
          </button>
        </div>
        <div class="flex flex-wrap gap-2 justify-center">
          {#each availableCards as card}
            <Card
              {card}
              onClick={() => handleCardClick(card)}
            />
          {/each}
        </div>
      </div>

      <button
        class="w-full py-3 bg-primary text-white rounded-lg font-semibold hover:bg-opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={selectedCards.length !== 7 || isPlaying}
        on:click={playHand}
      >
        Play Hand
      </button>
    </div>
  </div>
</div>
