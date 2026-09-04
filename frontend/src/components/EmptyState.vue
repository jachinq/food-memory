<template>
  <div class="empty">
    <h3>{{ title }}</h3>
    <p>{{ text }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; text: string }>()
</script>
