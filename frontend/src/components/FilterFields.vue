<template>
  <div class="field">
    <span>状态</span>
    <select :value="status" @change="emit('update:status', ($event.target as HTMLSelectElement).value)">
      <option value="">全部</option>
      <option v-for="(label, key) in STATUS_LABEL" :key="key" :value="key">{{ label }}</option>
    </select>
  </div>
  <div class="field">
    <span>标签</span>
    <input :value="tag" placeholder="例如 鸡肉" @input="emit('update:tag', ($event.target as HTMLInputElement).value)" />
  </div>
  <div class="field">
    <span>是否做过</span>
    <select :value="cooked" @change="emit('update:cooked', ($event.target as HTMLSelectElement).value)">
      <option value="">全部</option>
      <option value="true">已做</option>
      <option value="false">未做</option>
    </select>
  </div>
  <div class="field">
    <span>是否成功</span>
    <select :value="success" @change="emit('update:success', ($event.target as HTMLSelectElement).value)">
      <option value="">全部</option>
      <option value="true">做成功</option>
      <option value="false">其他</option>
    </select>
  </div>
  <button type="button" class="btn btn-ghost" @click="clear">一键清空</button>
</template>

<script setup lang="ts">
import { STATUS_LABEL } from '../types'

defineProps<{ status: string; tag: string; cooked: string; success: string }>()
const emit = defineEmits<{
  'update:status': [string]
  'update:tag': [string]
  'update:cooked': [string]
  'update:success': [string]
}>()

function clear() {
  emit('update:status', '')
  emit('update:tag', '')
  emit('update:cooked', '')
  emit('update:success', '')
}
</script>
