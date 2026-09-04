<template>
  <div>
    <div class="field">
      <span>状态</span>
      <div class="chip-row">
        <button type="button" class="chip" :class="{ on: !status }" @click="emit('update:status', '')">全部</button>
        <button
          v-for="(label, key) in STATUS_LABEL"
          :key="key"
          type="button"
          class="chip"
          :class="{ on: status === key }"
          @click="emit('update:status', String(key))"
        >{{ label }}</button>
      </div>
    </div>
    <div class="field">
      <span>标签</span>
      <input :value="tag" placeholder="例如 鸡肉" @input="emit('update:tag', ($event.target as HTMLInputElement).value)" />
    </div>
    <div class="field">
      <span>是否做过</span>
      <div class="chip-row">
        <button type="button" class="chip" :class="{ on: !cooked }" @click="emit('update:cooked', '')">全部</button>
        <button type="button" class="chip" :class="{ on: cooked === 'true' }" @click="emit('update:cooked', 'true')">已做</button>
        <button type="button" class="chip" :class="{ on: cooked === 'false' }" @click="emit('update:cooked', 'false')">未做</button>
      </div>
    </div>
    <div class="field">
      <span>是否成功</span>
      <div class="chip-row">
        <button type="button" class="chip" :class="{ on: !success }" @click="emit('update:success', '')">全部</button>
        <button type="button" class="chip" :class="{ on: success === 'true' }" @click="emit('update:success', 'true')">做成功</button>
        <button type="button" class="chip" :class="{ on: success === 'false' }" @click="emit('update:success', 'false')">其他</button>
      </div>
    </div>
    <button type="button" class="btn btn-ghost" @click="clear">一键清空</button>
  </div>
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
