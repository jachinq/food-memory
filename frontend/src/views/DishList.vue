<template>
  <div>
    <p class="section-kicker">Menu</p>
    <h1 class="section-title">菜牌</h1>
    <SearchBar v-model="keyword" placeholder="菜名、食材、标签…" @submit="reload" />
    <div class="menu-cats" role="tablist" aria-label="菜牌栏目">
      <button type="button" class="menu-cat" :class="{ on: section === 'all' }" @click="selectSection('all')">全部</button>
      <button
        v-for="item in sections"
        :key="item.name"
        type="button"
        class="menu-cat"
        :class="{ on: section === item.name }"
        @click="selectSection(item.name)"
      >{{ item.name }}</button>
      <button type="button" class="menu-cat" :class="{ on: section === 'other' }" @click="selectSection('other')">其他</button>
    </div>
    <div class="toolbar">
      <button class="btn btn-ghost mobile-only" type="button" @click="drawer = true">筛选</button>
      <button class="btn btn-ghost" type="button" @click="reload">搜索</button>
    </div>
    <div class="list-layout">
      <FilterSidebar v-model:status="status" v-model:tag="tag" v-model:cooked="cooked" v-model:success="success" />
      <div>
        <DishGrid :items="items" empty-title="没找到相关菜品" empty-text="试试搜索菜名、食材或标签。" />
        <div v-if="total > items.length" class="toolbar">
          <button class="btn btn-ghost" type="button" @click="more">更多</button>
        </div>
      </div>
    </div>
    <FilterDrawer :open="drawer" v-model:status="status" v-model:tag="tag" v-model:cooked="cooked" v-model:success="success" @close="drawer = false; reload()" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { fetchDishes } from '../api/dishes'
import { fetchTags } from '../api/tags'
import type { Dish, Tag } from '../types'
import { menuSectionTags } from '../types'
import SearchBar from '../components/SearchBar.vue'
import DishGrid from '../components/DishGrid.vue'
import FilterSidebar from '../components/FilterSidebar.vue'
import FilterDrawer from '../components/FilterDrawer.vue'

const route = useRoute()
const keyword = ref(String(route.query.keyword || ''))
const status = ref('')
const tag = ref('')
const cooked = ref('')
const success = ref('')
const untagged = ref(false)
const items = ref<Dish[]>([])
const page = ref(1)
const total = ref(0)
const drawer = ref(false)
const sections = ref<Tag[]>([])

const section = computed(() => {
  if (untagged.value) return 'other'
  if (tag.value) return tag.value
  return 'all'
})

function selectSection(name: string) {
  if (name === 'all') {
    tag.value = ''
    untagged.value = false
    return
  }
  if (name === 'other') {
    tag.value = ''
    untagged.value = true
    return
  }
  untagged.value = false
  tag.value = name
}

async function load(reset = true) {
  if (reset) page.value = 1
  const data = await fetchDishes({
    keyword: keyword.value,
    status: status.value,
    tag: untagged.value ? undefined : tag.value,
    untagged: untagged.value || undefined,
    cooked: cooked.value || undefined,
    success: success.value || undefined,
    page: page.value,
    pageSize: 12,
  })
  total.value = data.total
  items.value = reset ? data.items : items.value.concat(data.items)
}
function reload() { load(true) }
function more() { page.value += 1; load(false) }

onMounted(async () => {
  try {
    sections.value = menuSectionTags(await fetchTags())
  } catch {
    sections.value = []
  }
  await reload()
})
watch(() => route.query.keyword, (v) => {
  keyword.value = String(v || '')
  reload()
})
watch([status, cooked, success, untagged], reload)
let tagTimer = 0
watch(tag, () => {
  if (tag.value) untagged.value = false
  window.clearTimeout(tagTimer)
  tagTimer = window.setTimeout(reload, 300)
})
</script>
