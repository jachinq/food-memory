<template>
  <main class="unlock-page">
    <div class="unlock-card">
      <h1 class="brand">食忆<small>PRIVATE KITCHEN</small></h1>
      <p class="muted">这本账上了锁。输入口令，掀开封面。</p>
      <form @submit.prevent="unlock">
        <label class="field"><span>访问口令</span><input v-model="token" type="password" autocomplete="current-password" /></label>
        <p v-if="error" class="muted">{{ error }}</p>
        <button class="btn btn-primary" type="submit">掀开</button>
      </form>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { setToken } from '../api/client'
import { fetchHomeSummary } from '../api/home'

const router = useRouter()
const token = ref('')
const error = ref('')

async function unlock() {
  setToken(token.value.trim())
  try {
    await fetchHomeSummary()
    router.replace('/')
  } catch {
    error.value = '口令不正确'
    setToken('')
  }
}
</script>
