import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { installSurface, isIosDevice } from '../pwa/installSurface'

type PromptEvent = Event & { prompt: () => Promise<void> }

export const useInstall = defineStore('install', () => {
  const canPrompt = ref(false)
  const standalone = ref(false)
  const ios = ref(false)
  let deferred: PromptEvent | null = null

  const surface = computed(() =>
    installSurface({
      standalone: standalone.value,
      ios: ios.value,
      canPrompt: canPrompt.value,
    }),
  )

  function hydrate() {
    standalone.value =
      window.matchMedia('(display-mode: standalone)').matches ||
      Boolean((navigator as Navigator & { standalone?: boolean }).standalone)
    ios.value = isIosDevice(navigator.userAgent, navigator.maxTouchPoints ?? 0)
  }

  function onBeforeInstall(e: Event) {
    e.preventDefault()
    deferred = e as PromptEvent
    canPrompt.value = true
  }

  async function install() {
    if (!deferred) return
    await deferred.prompt()
    deferred = null
    canPrompt.value = false
  }

  function markInstalled() {
    deferred = null
    canPrompt.value = false
    standalone.value = true
  }

  return { surface, hydrate, onBeforeInstall, install, markInstalled }
})
