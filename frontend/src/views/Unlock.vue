<template>
  <main class="unlock-page">
    <div class="unlock-card">
      <h1 class="brand">食忆<small>HOUSE MENU</small></h1>
      <p class="muted">厨房上了锁。输入口令，看菜牌。</p>
      <form @submit.prevent="unlock">
        <label class="field"><span>访问口令</span><input v-model="token" type="password" autocomplete="current-password" /></label>
        <p v-if="error" class="muted">{{ error }}</p>
        <button class="btn btn-primary" type="submit">看菜牌</button>
      </form>
      <InstallActions />
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { setToken } from '../api/client'
import { fetchHomeSummary } from '../api/home'
import InstallActions from '../components/InstallActions.vue'

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
