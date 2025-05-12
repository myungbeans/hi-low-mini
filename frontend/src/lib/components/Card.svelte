<script lang="ts">
  import type { Card as CardType } from '$lib/types';
  
  export let card: CardType;
  export let isSelected = false;
  export let isHighlighted = false;
  export let onClick: () => void;

  $: isOperator = card.type === 'operator';
  
  function handleClick() {
    onClick();
  }
</script>

<button
  class="w-12 h-20 rounded-lg shadow-md transition-all duration-200 flex items-center justify-center text-lg font-bold"
  class:bg-primary={isSelected}
  class:bg-secondary={isOperator && !isSelected}
  class:text-white={isSelected || (isOperator && !isSelected)}
  class:bg-white={!isSelected && !isOperator}
  class:text-primary={!isSelected && !isOperator}
  class:ring-4={isHighlighted}
  class:ring-secondary={isHighlighted}
  on:click={handleClick}
>
  {card.value}
</button>

<style>
  button {
    transform-style: preserve-3d;
    backface-visibility: hidden;
  }

  button:hover:not(:disabled) {
    transform: translateY(-2px);
  }
</style> 