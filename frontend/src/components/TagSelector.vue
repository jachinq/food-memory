<template>
  <div class="tag-picker">
    <div class="row" style="margin-bottom:8px">
      <span v-for="tag in modelValue" :key="tag.name + tag.type" class="tag">
        {{ tag.name }}
        <button type="button" class="icon-btn" aria-label="移除标签" @click="remove(tag)">×</button>
      </span>
    </div>
    <input
      :value="draft"
      placeholder="输入标签回车添加，可从已有标签选择"
      @input="draft = ($event.target as HTMLInputElement).value"
      @focus="load"
      @keydown.enter.prevent="addDraft"
    />
    <div v-if="filtered.length" class="row" style="margin-top:8px">
      <button v-for="t in filtered" :key="t.id" type="button" class="btn btn-ghost" @click="pick(t)">{{ t.name }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { fetchTags } from '../api/tags'
import type { Tag } from '../types'

const props = defineProps<{
  modelValue: { name: string; type: string }[]
  type?: string
}>()
const emit = defineEmits<{ 'update:modelValue': [{ name: string; type: string }[]] }>()

const draft = ref('')
const all = ref<Tag[]>([])

const filtered = computed(() => {
  const q = draft.value.trim()
  return all.value.filter((t) => {
    if (props.modelValue.some((x) => x.name === t.name)) return false
    if (props.type && t.type !== props.type) return false
    return !q || t.name.includes(q)
  }).slice(0, 8)
})

async function load() {
  if (all.value.length) return
  all.value = await fetchTags(props.type)
}

function add(name: string, type = props.type || 'custom') {
  name = name.trim()
  if (!name) return
  emit('update:modelValue', [...props.modelValue, { name, type }])
  draft.value = ''
}
function addDraft() { add(draft.value) }
function pick(t: Tag) { add(t.name, t.type) }
function remove(tag: { name: string; type: string }) {
  emit('update:modelValue', props.modelValue.filter((t) => !(t.name === tag.name && t.type === tag.type)))
}
</script>
