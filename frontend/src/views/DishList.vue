<template>
  <div>
    <p class="section-kicker">Archive</p>
    <h1 class="section-title">菜品库</h1>
    <SearchBar v-model="keyword" placeholder="我记得是鸡肉、甜辣、下饭…" @submit="reload" />
    <div class="toolbar">
      <button class="btn btn-ghost mobile-only" type="button" @click="drawer = true">筛选</button>
      <button class="btn btn-ghost" type="button" @click="reload">按记忆搜</button>
    </div>
    <div class="list-layout">
      <FilterSidebar v-model:status="status" v-model:tag="tag" v-model:cooked="cooked" v-model:success="success" />
      <div>
        <DishGrid :items="items" empty-title="没找到相关菜品" empty-text="试试搜索食材、口味或标签。" />
        <div v-if="total > items.length" class="toolbar">
          <button class="btn btn-ghost" type="button" @click="more">再翻一页</button>
        </div>
      </div>
    </div>
    <FilterDrawer :open="drawer" v-model:status="status" v-model:tag="tag" v-model:cooked="cooked" v-model:success="success" @close="drawer = false; reload()" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { fetchDishes } from '../api/dishes'
import type { Dish } from '../types'
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
const items = ref<Dish[]>([])
const page = ref(1)
const total = ref(0)
const drawer = ref(false)

async function load(reset = true) {
  if (reset) page.value = 1
  const data = await fetchDishes({
    keyword: keyword.value,
    status: status.value,
    tag: tag.value,
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

onMounted(reload)
watch(() => route.query.keyword, (v) => {
  keyword.value = String(v || '')
  reload()
})
watch([status, cooked, success], reload)
let tagTimer = 0
watch(tag, () => {
  window.clearTimeout(tagTimer)
  tagTimer = window.setTimeout(reload, 300)
})
</script>
