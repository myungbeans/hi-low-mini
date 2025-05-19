<script lang="ts">
  import '../app.css';
  import Navbar from '$lib/components/Navbar.svelte';
  import { darkMode } from '$lib/stores/theme';
</script>
 
<div class="min-h-screen bg-white dark:bg-gray-900 transition-colors">
  <Navbar />
  <slot />
</div> 