<template>
  <div class="uploader" :class="{ 'is-busy': busy, 'has-camera': showCamera }">
    <label class="uploader-preview" :for="albumId">
      <img v-if="preview" :src="preview" alt="预览" />
      <span v-else class="muted">{{ placeholder }}</span>
    </label>
    <div class="uploader-actions">
      <label v-if="showCamera" class="btn btn-primary">
        拍照
        <input
          :id="cameraId"
          ref="cameraInput"
          v-camera-capture
          class="uploader-hit"
          type="file"
          accept="image/*"
          capture="environment"
          :disabled="busy"
          @click="armCameraInput"
          @change="onFile"
        />
      </label>
      <label class="btn btn-ghost">
        {{ showCamera ? '相册' : '上传图片' }}
        <input
          :id="albumId"
          class="uploader-hit"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          :disabled="busy"
          @change="onFile"
        />
      </label>
    </div>
    <p v-if="busy" class="muted uploader-status">正在上传…</p>
    <p v-if="error" class="uploader-error">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, useId } from 'vue'
import type { Directive } from 'vue'
import { uploadImage } from '../api/upload'
import type { Attachment } from '../types'

const ALLOWED_MIME = new Set(['image/jpeg', 'image/png', 'image/webp', 'image/jpg'])

const props = defineProps<{
  modelValue?: string
  bizType?: string
  bizId?: number
}>()
const emit = defineEmits<{
  'update:modelValue': [string]
  uploaded: [Attachment]
}>()

const uid = useId()
const cameraId = `${uid}-camera`
const albumId = `${uid}-album`
const cameraInput = ref<HTMLInputElement | null>(null)

function isMobileCaptureEnv() {
  if (typeof navigator === 'undefined') return false
  const ua = navigator.userAgent
  if (/Android|webOS|BlackBerry|IEMobile|Opera Mini/i.test(ua)) return true
  if (/iPhone|iPod/i.test(ua)) return true
  if (/iPad/i.test(ua)) return true
  if (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1) return true
  return false
}

const localUrl = ref('')
const busy = ref(false)
const error = ref('')
const showCamera = ref(isMobileCaptureEnv())
const preview = computed(() => localUrl.value || props.modelValue || '')
const placeholder = computed(() =>
  showCamera.value ? '贴一张照片，或直接拍照' : '从电脑选择一张图片',
)

/** Android Chrome 读的是 HTML 属性，不是 input.capture 这个 JS 属性。 */
function applyCaptureAttr(el: HTMLInputElement) {
  el.setAttribute('accept', 'image/*')
  el.setAttribute('capture', 'environment')
}

const vCameraCapture: Directive<HTMLInputElement> = {
  mounted: applyCaptureAttr,
  updated: applyCaptureAttr,
}

function armCameraInput(e: Event) {
  applyCaptureAttr(e.currentTarget as HTMLInputElement)
}

function refreshCameraUi() {
  showCamera.value = isMobileCaptureEnv()
}

onMounted(() => {
  refreshCameraUi()
  if (cameraInput.value) applyCaptureAttr(cameraInput.value)
  window.addEventListener('orientationchange', refreshCameraUi)
})
onUnmounted(() => {
  window.removeEventListener('orientationchange', refreshCameraUi)
})

function isAllowedImage(file: File): boolean {
  if (ALLOWED_MIME.has(file.type)) return true
  if (!file.type) {
    const ext = file.name.split('.').pop()?.toLowerCase()
    return ext === 'jpg' || ext === 'jpeg' || ext === 'png' || ext === 'webp'
  }
  return false
}

async function onFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  error.value = ''
  if (!isAllowedImage(file)) {
    error.value = '仅支持 JPEG、PNG、WebP，当前格式无法上传'
    return
  }

  if (localUrl.value) URL.revokeObjectURL(localUrl.value)
  localUrl.value = URL.createObjectURL(file)
  busy.value = true
  try {
    const att = await uploadImage(file, props.bizType, props.bizId)
    emit('update:modelValue', att.file_url)
    emit('uploaded', att)
  } catch (err) {
    error.value = err instanceof Error ? err.message : '上传失败'
  } finally {
    busy.value = false
  }
}
</script>
