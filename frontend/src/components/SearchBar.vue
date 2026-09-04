<template>
  <form class="search-bar" @submit.prevent="$emit('submit')">
    <input
      :value="modelValue"
      :placeholder="placeholder"
      aria-label="搜索菜品"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <button type="submit">找</button>
  </form>
</template>

<script setup lang="ts">
defineProps<{ modelValue: string; placeholder?: string }>()
defineEmits<{ 'update:modelValue': [string]; submit: [] }>()
</script>
