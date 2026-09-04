<template>
  <div v-if="open" class="sheet" @click.self="$emit('close')">
    <div class="sheet-body">
      <h3>筛选</h3>
      <FilterFields :status="status" :tag="tag" :cooked="cooked" :success="success" @update:status="$emit('update:status', $event)" @update:tag="$emit('update:tag', $event)" @update:cooked="$emit('update:cooked', $event)" @update:success="$emit('update:success', $event)" />
      <button type="button" class="btn btn-primary" style="margin-top:12px" @click="$emit('close')">完成</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import FilterFields from './FilterFields.vue'
defineProps<{ open: boolean; status: string; tag: string; cooked: string; success: string }>()
defineEmits<{
  close: []
  'update:status': [string]
  'update:tag': [string]
  'update:cooked': [string]
  'update:success': [string]
}>()
</script>
