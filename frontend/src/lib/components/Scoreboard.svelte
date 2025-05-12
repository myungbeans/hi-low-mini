<script lang="ts">
  export let currentScore: number;
  export let scoreHistory: number[];
  export let shouldVibrate = false;
  export let vibrateIntensity = 0;

  $: vibrateClass = `text-4xl font-bold text-secondary ${shouldVibrate ? `vibrate-${Math.min(vibrateIntensity, 3)}` : ''}`;

  function formatNumber(num: number): string {
    return num.toFixed(2);
  }
</script>

<div class="bg-white rounded-lg p-4 shadow-lg">
  <div class="text-center">
    <h2 class="text-2xl font-bold text-primary">Current Score</h2>
    <p class={vibrateClass}>{formatNumber(currentScore)}</p>
  </div>
</div>

<style>
  @keyframes vibrate-1 {
    0% { transform: translate(0); }
    25% { transform: translate(-0.5px, 0.5px); }
    50% { transform: translate(0.5px, -0.5px); }
    75% { transform: translate(-0.5px, -0.5px); }
    100% { transform: translate(0); }
  }

  @keyframes vibrate-2 {
    0% { transform: translate(0); }
    25% { transform: translate(-2px, 2px); }
    50% { transform: translate(2px, -2px); }
    75% { transform: translate(-2px, -2px); }
    100% { transform: translate(0); }
  }

  @keyframes vibrate-3 {
    0% { transform: translate(0); }
    25% { transform: translate(-4px, 4px); }
    50% { transform: translate(4px, -4px); }
    75% { transform: translate(-4px, -4px); }
    100% { transform: translate(0); }
  }

  .vibrate-1 { animation: vibrate-1 0.2s linear; }
  .vibrate-2 { animation: vibrate-2 0.2s linear; }
  .vibrate-3 { animation: vibrate-3 0.2s linear; }
</style> 