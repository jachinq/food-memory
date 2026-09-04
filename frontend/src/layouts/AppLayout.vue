<template>
  <div class="app-shell">
    <header class="topbar">
      <router-link class="brand" to="/">食忆</router-link>
      <nav class="nav-links">
        <router-link to="/">首页</router-link>
        <router-link to="/dishes">菜品库</router-link>
        <router-link to="/recook">复做清单</router-link>
      </nav>
      <div class="topbar-search">
        <SearchBar v-model="keyword" placeholder="我记得是鸡肉、甜辣、下饭..." @submit="goSearch" />
      </div>
      <router-link class="btn btn-primary" to="/dishes/new">+ 新增菜品</router-link>
    </header>

    <div class="page">
      <div class="mobile-only mobile-head">
        <strong class="brand">食忆</strong>
        <router-link to="/dishes">搜索</router-link>
      </div>
      <router-view />
    </div>

    <nav class="bottom-nav">
      <router-link to="/">首页</router-link>
      <router-link to="/dishes">菜品库</router-link>
      <button class="fab" type="button" aria-label="新增" @click="sheet = true">+</button>
      <router-link to="/recook">复做</router-link>
      <a href="#me" @click.prevent="sheet = true">我的</a>
    </nav>

    <div v-if="sheet" class="sheet" @click.self="sheet = false">
      <div class="sheet-body">
        <button type="button" @click="go('/dishes/new')">新增菜品</button>
        <button type="button" @click="go('/dishes?action=record')">记录制作</button>
        <button type="button" @click="sheet = false">取消</button>
      </div>
    </div>
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
