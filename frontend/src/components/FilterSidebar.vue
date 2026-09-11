<template>
  <aside class="filter-side pc-only">
        <h3>筛选</h3>
    <FilterFields
      :status="status"
      :tag="tag"
      :cooked="cooked"
      :success="success"
      @update:status="$emit('update:status', $event)"
      @update:tag="$emit('update:tag', $event)"
      @update:cooked="$emit('update:cooked', $event)"
      @update:success="$emit('update:success', $event)"
    />
  </aside>
</template>

<script setup lang="ts">
import FilterFields from './FilterFields.vue'
defineProps<{ status: string; tag: string; cooked: string; success: string }>()
defineEmits<{
  'update:status': [string]
  'update:tag': [string]
  'update:cooked': [string]
  'update:success': [string]
}>()
</script>
