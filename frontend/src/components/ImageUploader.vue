<template>
  <label class="uploader">
    <input type="file" accept="image/jpeg,image/png,image/webp" hidden @change="onFile" />
    <img v-if="preview" :src="preview" alt="预览" />
    <span v-else class="muted">贴一张照片，或直接拍照</span>
  </label>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { uploadImage } from '../api/upload'
import type { Attachment } from '../types'

const props = defineProps<{
  modelValue?: string
  bizType?: string
  bizId?: number
}>()
const emit = defineEmits<{
  'update:modelValue': [string]
  uploaded: [Attachment]
}>()

const localUrl = ref('')
const preview = computed(() => localUrl.value || props.modelValue || '')

async function onFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  localUrl.value = URL.createObjectURL(file)
  const att = await uploadImage(file, props.bizType, props.bizId)
  emit('update:modelValue', att.file_url)
  emit('uploaded', att)
}
</script>
