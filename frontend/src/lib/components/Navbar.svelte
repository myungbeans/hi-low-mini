<script lang="ts">
  import { browser } from '$app/environment';
  import { fade } from 'svelte/transition';
  import { darkMode } from '$lib/stores/theme';
  import { clearGameCache } from '$lib/stores/game';
  
  let showInstructions = false;
  let showSettings = false;
  let showClearCacheConfirm = false;

  function handleClearCache() {
    clearGameCache();
    showClearCacheConfirm = false;
    showSettings = false;
  }
</script>

{#if browser}
<nav class="w-full bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-4 py-2.5">
  <div class="max-w-screen-xl mx-auto flex justify-between items-center">
    {#if showInstructions}
      <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" 
           on:click|self={() => showInstructions = false}
           transition:fade>
        <div class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-md w-full mx-4">
          <h2 class="text-xl font-bold mb-4 dark:text-white">How to Play</h2>
          <div class="space-y-2 dark:text-gray-300">
            <p>• You receive 10 cards daily (Numbers and Operators)</p>
            <p>• Select and arrange 7 cards to create your hand</p>
            <p>• Cards are evaluated left to right</p>
            <p>• Try to achieve the highest score possible!</p>
          </div>
          <button 
            class="mt-4 px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors dark:text-white"
            on:click={() => showInstructions = false}>
            Close
          </button>
        </div>
      </div>
    {/if}

    {#if showSettings}
      <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" 
           on:click|self={() => showSettings = false}
           transition:fade>
        <div class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-md w-full mx-4">
          <h2 class="text-xl font-bold mb-4 dark:text-white">Settings</h2>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <span class="dark:text-white">Dark Mode</span>
              <button 
                class="w-12 h-6 {$darkMode ? 'bg-blue-600' : 'bg-gray-200'} rounded-full relative transition-colors"
                on:click={() => darkMode.set(!$darkMode)}
                aria-label="Toggle dark mode">
                <span 
                  class="absolute top-1 w-4 h-4 bg-white rounded-full transition-transform duration-200"
                  style="left: {$darkMode ? '1.75rem' : '0.25rem'}">
                </span>
              </button>
            </div>
            <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
              <button 
                class="w-full px-4 py-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                on:click={() => showClearCacheConfirm = true}>
                Clear Cached Data
              </button>
            </div>
          </div>
          <button 
            class="mt-4 px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors dark:text-white"
            on:click={() => showSettings = false}>
            Close
          </button>
        </div>
      </div>
    {/if}

    {#if showClearCacheConfirm}
      <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" 
           on:click|self={() => showClearCacheConfirm = false}
           transition:fade>
        <div class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-md w-full mx-4">
          <h2 class="text-xl font-bold mb-4 dark:text-white">Clear Cache</h2>
          <p class="text-gray-600 dark:text-gray-300 mb-6">
            Are you sure you want to clear all cached data? This will reset your current game and you'll need to fetch a new one.
          </p>
          <div class="flex space-x-4">
            <button 
              class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
              on:click={handleClearCache}>
              Clear
            </button>
            <button 
              class="flex-1 px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors dark:text-white"
              on:click={() => showClearCacheConfirm = false}>
              Cancel
            </button>
          </div>
        </div>
      </div>
    {/if}

    <div class="text-xl font-bold dark:text-white">Hi-Lo Mini</div>
    
    <div class="flex items-center space-x-4">
      <!-- Leaderboard Icon -->
      <button class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full" aria-label="Leaderboard">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 dark:stroke-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </button>

      <!-- Help Icon -->
      <button 
        class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full" 
        aria-label="Help"
        on:click={() => showInstructions = true}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 dark:stroke-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </button>

      <!-- Settings Icon -->
      <button 
        class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full" 
        aria-label="Settings"
        on:click={() => showSettings = true}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 dark:stroke-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      </button>
    </div>
  </div>
</nav>
{:else}
<!-- Render a placeholder during SSR -->
<nav class="w-full bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-4 py-2.5">
  <div class="max-w-screen-xl mx-auto flex justify-between items-center">
    <div class="text-xl font-bold dark:text-white">Hi-Lo Mini</div>
    <div class="flex items-center space-x-4">
      <!-- Placeholder icons with same dimensions but no interactivity -->
      <div class="w-6 h-6"></div>
      <div class="w-6 h-6"></div>
      <div class="w-6 h-6"></div>
    </div>
  </div>
</nav>
{/if} 