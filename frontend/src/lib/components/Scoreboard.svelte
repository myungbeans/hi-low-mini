<script lang="ts">
  export let currentScore: number;
  export let scoreHistory: number[];
</script>

<div class="bg-white rounded-lg p-4 shadow-lg">
  <div class="text-center">
    <h2 class="text-2xl font-bold text-primary">Current Score</h2>
    <p class="text-4xl font-bold text-secondary">{currentScore}</p>
  </div>
  
  {#if scoreHistory.length > 0}
    <div class="mt-4">
      <h3 class="text-lg font-semibold mb-2">History</h3>
      <div class="flex flex-wrap gap-2">
        {#each scoreHistory as score}
          <div class="px-3 py-1 bg-gray-100 rounded-full text-sm">
            {score}
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div> 