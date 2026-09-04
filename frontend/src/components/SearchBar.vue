<template>
  <form class="search-bar" @submit.prevent="$emit('submit')">
    <input :value="modelValue" :placeholder="placeholder" @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)" />
  </form>
</template>

<script setup lang="ts">
defineProps<{ modelValue: string; placeholder?: string }>()
defineEmits<{ 'update:modelValue': [string]; submit: [] }>()
</script>
