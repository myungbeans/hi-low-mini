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
  let lastPlayResponse: PlayHandResponse | null = null;
  let hasPlayed = false;
  let shouldVibrateScore = false;
  let vibrateIntensity = 0;
  let numberCardCount = 0;
  let lastNumberCardIntensity = 0;

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
      // Prevent removing cards if hand has been played
      if (hasPlayed) {
        error = 'Cannot modify hand after it has been played';
        setTimeout(() => { error = null; }, 2000);
        return;
      }
      
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
      // Prevent adding cards if hand has been played
      if (hasPlayed) {
        error = 'Cannot modify hand after it has been played';
        setTimeout(() => { error = null; }, 2000);
        return;
      }

      // Validate card placement
      const isOperator = card.type === 'operator';
      
      // Check if it's the first card
      if (selectedCards.length === 0 && isOperator) {
        error = 'First card must be a number';
        setTimeout(() => { error = null; }, 2000);
        return;
      }

      // Check for consecutive operators
      if (isOperator && selectedCards.length > 0 && selectedCards[selectedCards.length - 1].type === 'operator') {
        error = 'Cannot place two operators in a row';
        setTimeout(() => { error = null; }, 2000);
        return;
      }

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

  function calculateVibrateIntensity(numberCards: number) {
    // Cap at 3 number cards
    const cappedNumberCards = Math.min(numberCards, 3);
    // More dramatic exponential growth starting from the second card
    if (numberCards === 1) {
      return 1; // Very subtle effect for first card
    }
    // Exponential growth with higher base for more dramatic increase, max at level 3
    return Math.min(Math.floor(Math.pow(2.2, cappedNumberCards - 1)), 3);
  }

  async function animateHand(data: PlayHandResponse) {
    isPlaying = true;
    let scoreCounterIndex = 0;
    numberCardCount = 0;
    lastNumberCardIntensity = 0;
    
    // Reset score when replaying
    if (hasPlayed) {
      currentScore = 0;
    }

    // Start with first card immediately
    highlightIndex = 0;
    if (selectedCards[0].type === 'number') {
      numberCardCount++;
      lastNumberCardIntensity = calculateVibrateIntensity(numberCardCount);
      vibrateIntensity = lastNumberCardIntensity;
      shouldVibrateScore = true;
      const score = Number(data.scoreCounter[scoreCounterIndex]);
      if (!isNaN(score)) {
        currentScore = score;
      }
      scoreCounterIndex++;
    } else {
      // Operator card - use last number card's intensity
      vibrateIntensity = lastNumberCardIntensity;
      shouldVibrateScore = vibrateIntensity > 0;
    }
    // Reset vibration after a short delay if vibration was triggered
    if (shouldVibrateScore) {
      setTimeout(() => {
        shouldVibrateScore = false;
      }, 200);
    }

    // Animate through the remaining cards and update score
    for (let i = 1; i < selectedCards.length; i++) {
      await new Promise(resolve => setTimeout(resolve, 500));
      highlightIndex = i;
      
      if (selectedCards[i].type === 'number') {
        numberCardCount++;
        lastNumberCardIntensity = calculateVibrateIntensity(numberCardCount);
        vibrateIntensity = lastNumberCardIntensity;
        shouldVibrateScore = true;
        const score = Number(data.scoreCounter[scoreCounterIndex]);
        if (!isNaN(score)) {
          currentScore = score;
        }
        scoreCounterIndex++;
      } else {
        // Operator card - use last number card's intensity
        vibrateIntensity = lastNumberCardIntensity;
        shouldVibrateScore = vibrateIntensity > 0;
      }
      // Reset vibration after a short delay if vibration was triggered
      if (shouldVibrateScore) {
        setTimeout(() => {
          shouldVibrateScore = false;
        }, 200);
      }
    }

    // Ensure final score is set if it's valid
    const finalScore = Number(data.scoreCounter[data.scoreCounter.length - 1]);
    if (!isNaN(finalScore)) {
      currentScore = finalScore;
      if (!hasPlayed) {
        scoreHistory = [...scoreHistory, currentScore];
      }
    }
    
    // Reset for next play
    await new Promise(resolve => setTimeout(resolve, 500));
    highlightIndex = -1;
    isPlaying = false;
    hasPlayed = true;
  }

  async function playHand() {
    if (selectedCards.length !== 7 || isPlaying || !game) {
      return;
    }

    try {
      error = null;

      if (!hasPlayed) {
        // First time playing this hand - make the API request
        // Validate selected cards
        if (!selectedCards.every(card => card && typeof card.value !== 'undefined' && card.type)) {
          throw new Error('Invalid card data in selected cards');
        }
        
        // Create the hand data with proper card types
        const handCards = selectedCards.map(card => ({
          value: card.value.toString(),
          type: card.type === 'operator' ? 2 : 1
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

        lastPlayResponse = data;
        await animateHand(data);
      } else {
        // Replaying the hand - use stored response
        if (lastPlayResponse) {
          await animateHand(lastPlayResponse);
        }
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to play hand';
      isPlaying = false;
      highlightIndex = -1;
    }
  }
</script>

<div class="min-h-screen bg-background p-4">
  <div class="max-w-4xl mx-auto space-y-8">
    <h1 class="text-4xl font-bold text-center text-text">Hi-Lo Mini</h1>
    
    {#if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{error}</span>
      </div>
    {/if}

    <Scoreboard {currentScore} {scoreHistory} shouldVibrate={shouldVibrateScore} {vibrateIntensity} />
    
    <div class="space-y-6">
      <div class="bg-white rounded-lg p-4 shadow-lg">
        <h2 class="text-xl font-semibold mb-4">Hand ({selectedCards.length}/7)</h2>
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
            <div class="w-12 h-20 border-2 border-dashed border-gray-300 rounded-lg" />
          {/each}
        </div>
      </div>

      <div class="bg-white rounded-lg p-4 shadow-lg">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">Pool ({availableCards.length})</h2>
          <button
            class="px-4 py-2 bg-secondary text-white rounded-lg hover:bg-opacity-90 transition-colors"
            on:click={shuffleCards}
          >
            Shuffle
          </button>
        </div>
        <div class="flex gap-1.5 justify-center flex-nowrap overflow-x-auto px-2">
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
        {hasPlayed ? 'Replay Hand' : 'Play Hand'}
      </button>
    </div>
  </div>
</div>
 