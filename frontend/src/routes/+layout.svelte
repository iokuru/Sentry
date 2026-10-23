<script lang="ts">
  import '../app.css';
  let { children } = $props();
</script>

<div class="app-container">
  {@render children()}
</div>

<style>
  .app-container {
    max-width: 640px;
    margin: 0 auto;
    padding: 24px 20px 80px;
    min-height: 100vh;
  }
</style>
