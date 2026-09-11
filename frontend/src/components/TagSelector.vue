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
      :placeholder="placeholder"
      @input="draft = ($event.target as HTMLInputElement).value"
      @focus="load"
      @keydown.enter.prevent="onEnter"
    />
    <div v-if="filtered.length" class="row" style="margin-top:8px">
      <button v-for="t in filtered" :key="t.id" type="button" class="btn btn-ghost" @mousedown.prevent @click="pick(t)">{{ t.name }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { fetchTags } from '../api/tags'
import type { Tag } from '../types'

const props = withDefaults(defineProps<{
  modelValue: { name: string; type: string }[]
  type?: string
  excludeTypes?: string[]
  placeholder?: string
}>(), {
  placeholder: '输入标签回车添加，可从已有标签选择',
})
const emit = defineEmits<{ 'update:modelValue': [{ name: string; type: string }[]] }>()

const draft = ref('')
const all = ref<Tag[]>([])

const filtered = computed(() => {
  const q = draft.value.trim()
  return all.value.filter((t) => {
    if (props.modelValue.some((x) => x.name === t.name)) return false
    if (props.type && t.type !== props.type) return false
    if (props.excludeTypes?.includes(String(t.type))) return false
    return !q || t.name.includes(q)
  }).slice(0, 8)
})

async function load() {
  all.value = await fetchTags(props.type)
}

onMounted(load)

function add(name: string, type = props.type || 'custom') {
  name = name.trim()
  if (!name) return
  const current = props.modelValue || []
  if (current.some((t) => t.name === name && t.type === type)) {
    draft.value = ''
    return
  }
  emit('update:modelValue', [...current, { name, type }])
  draft.value = ''
}
function addDraft() { add(draft.value) }
function onEnter(e: KeyboardEvent) {
  if (e.isComposing) return
  addDraft()
}
function pick(t: Tag) { add(t.name, t.type) }
function remove(tag: { name: string; type: string }) {
  emit('update:modelValue', (props.modelValue || []).filter((t) => !(t.name === tag.name && t.type === tag.type)))
}

defineExpose({ commit: addDraft })
</script>
