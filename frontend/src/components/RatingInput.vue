<template>
  <div class="stars" role="radiogroup">
    <button v-for="n in 5" :key="n" type="button" @click="$emit('update:modelValue', n)">
      {{ n <= (modelValue || 0) ? '★' : '☆' }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ modelValue: number | null }>()
defineEmits<{ 'update:modelValue': [number] }>()
</script>
