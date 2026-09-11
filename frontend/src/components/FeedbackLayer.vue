<template>
  <div class="notice-host" aria-live="polite">
    <button
      v-if="feedback.flash"
      :key="feedback.flash.id"
      type="button"
      class="notice"
      :data-tone="feedback.flash.kind"
      @click="feedback.dismissNotice()"
    >
      {{ feedback.flash.text }}
    </button>
  </div>

  <Teleport to="body">
    <Transition name="confirm">
      <div
        v-if="feedback.confirmText"
        class="confirm-scrim"
        role="presentation"
        @click.self="feedback.settle(false)"
      >
        <div
          ref="dialog"
          class="confirm-card"
          role="alertdialog"
          aria-modal="true"
          aria-labelledby="confirm-copy"
          tabindex="-1"
        >
          <p id="confirm-copy" class="confirm-copy">{{ feedback.confirmText }}</p>
          <div class="confirm-actions">
            <button class="btn btn-ghost" type="button" @click="feedback.settle(false)">取消</button>
            <button class="btn btn-danger" type="button" @click="feedback.settle(true)">删除</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { nextTick, onUnmounted, ref, watch } from 'vue'
import { useFeedback } from '../stores/feedback'

const feedback = useFeedback()
const dialog = ref<HTMLElement | null>(null)

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && feedback.confirmText) {
    e.preventDefault()
    feedback.settle(false)
  }
}

watch(
  () => feedback.confirmText,
  async (text) => {
    document.removeEventListener('keydown', onKey)
    if (!text) return
    document.addEventListener('keydown', onKey)
    await nextTick()
    dialog.value?.focus()
  },
)

onUnmounted(() => document.removeEventListener('keydown', onKey))
</script>
