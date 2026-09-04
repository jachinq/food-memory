<template>
  <div class="stars" role="radiogroup" aria-label="评分">
    <button
      v-for="n in 5"
      :key="n"
      type="button"
      class="icon-btn"
      :class="{ on: n <= (modelValue || 0) }"
      :aria-label="`${n} 分`"
      @click="$emit('update:modelValue', n)"
    >
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 3.2 14.6 9l6.4.6-4.9 4.1 1.5 6.3L12 16.8 6.4 20l1.5-6.3L3 9.6 9.4 9z" fill="currentColor"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ modelValue: number | null }>()
defineEmits<{ 'update:modelValue': [number] }>()
</script>
