<template>
  <main class="page" style="max-width:420px;padding-top:12vh">
    <h1 class="brand">食忆</h1>
    <p class="muted">此部署开启了访问口令。</p>
    <form class="group" @submit.prevent="unlock">
      <label class="field"><span>访问口令</span><input v-model="token" type="password" /></label>
      <p v-if="error">{{ error }}</p>
      <button class="btn btn-primary" type="submit">进入</button>
    </form>
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
