<template>
  <div class="app-shell">
    <header class="topbar">
      <router-link class="brand" to="/">
        食忆
        <small>KITCHEN LEDGER</small>
      </router-link>
      <nav class="nav-links">
        <router-link to="/" active-class="is-prefix-active" exact-active-class="router-link-active">首页</router-link>
        <router-link to="/dishes">菜品库</router-link>
        <router-link to="/recook">复做清单</router-link>
      </nav>
      <div class="topbar-search">
        <SearchBar v-model="keyword" placeholder="鸡肉、甜辣、下饭…" @submit="goSearch" />
      </div>
      <router-link class="btn btn-primary" to="/dishes/new">记一道菜</router-link>
    </header>

    <div class="page">
      <div class="mobile-only mobile-head">
        <router-link class="brand" to="/">食忆</router-link>
        <router-link class="muted" to="/dishes">检索</router-link>
      </div>
      <router-view v-slot="{ Component }">
        <component :is="Component" />
      </router-view>
    </div>

    <nav class="bottom-nav">
      <router-link to="/" active-class="is-prefix-active" exact-active-class="router-link-active">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M4 10.5 12 4l8 6.5V20H4z"/><path d="M9 20v-6h6v6"/></svg>
        首页
      </router-link>
      <router-link to="/dishes">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7"><rect x="4" y="5" width="16" height="14"/><path d="M4 10h16M9 5v14"/></svg>
        菜品
      </router-link>
      <button class="fab" type="button" aria-label="新增" @click="sheet = true">+</button>
      <router-link to="/recook">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7"><path d="M7 7h13v13H7z"/><path d="M4 4h13v3"/><path d="M10 12h6M10 16h4"/></svg>
        复做
      </router-link>
      <a href="#me" @click.prevent="sheet = true">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7"><circle cx="12" cy="8" r="3"/><path d="M5 19c1.4-3 3.8-4.5 7-4.5S17.6 16 19 19"/></svg>
        我的
      </a>
    </nav>

    <Transition name="sheet">
      <div v-if="sheet" class="sheet" @click.self="sheet = false">
        <div class="sheet-body">
          <button type="button" class="sheet-action" @click="go('/dishes/new')">记一道新菜</button>
          <button type="button" class="sheet-action" @click="go('/dishes?action=record')">记下一次制作</button>
          <button type="button" class="sheet-action" @click="sheet = false">先放一放</button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import SearchBar from '../components/SearchBar.vue'

const router = useRouter()
const keyword = ref('')
const sheet = ref(false)

function goSearch() {
  router.push({ path: '/dishes', query: { keyword: keyword.value } })
}
function go(path: string) {
  sheet.value = false
  router.push(path)
}
</script>
