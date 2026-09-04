<template>
  <div
    class="stars"
    role="radiogroup"
    :aria-label="label"
    @mouseleave="hover = 0"
  >
    <button
      v-for="n in 5"
      :key="n"
      type="button"
      class="icon-btn"
      :class="{ on: n <= filled }"
      :aria-label="`${n} 星`"
      :aria-checked="n === (modelValue || 0)"
      role="radio"
      @mouseenter="hover = n"
      @click="$emit('update:modelValue', n)"
    >
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 3.2 14.6 9l6.4.6-4.9 4.1 1.5 6.3L12 16.8 6.4 20l1.5-6.3L3 9.6 9.4 9z" fill="currentColor"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: number | null
  label?: string
}>(), {
  label: '评分',
})
defineEmits<{ 'update:modelValue': [number] }>()

const hover = ref(0)
const filled = computed(() => hover.value || props.modelValue || 0)
</script>
