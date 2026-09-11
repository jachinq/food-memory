<template>
  <div v-if="loading" class="skeleton" aria-hidden="true"></div>
  <div v-else class="home-page">
    <p class="stats-line" aria-label="记忆概览">
      <span>{{ data.stats.total_dishes }} 道菜</span>
      <span>{{ data.stats.cooked_count }} 已做</span>
      <span>{{ data.stats.success_count }} 成功</span>
      <span>本月 {{ data.stats.month_cook_count }} 次</span>
    </p>

    <section class="board" data-tone="chili" aria-labelledby="home-draw-title">
      <header class="board-head">
        <div>
          <p class="section-kicker">Today</p>
          <h2 id="home-draw-title" class="section-title">今日推荐</h2>
        </div>
        <p class="board-meta">不知道吃什么时看这里</p>
      </header>
      <article v-if="data.random_old" class="card hero-card">
        <div class="visual">
          <img v-if="data.random_old.cover_image_url" :src="data.random_old.cover_image_url" :alt="data.random_old.name" />
          <div v-else class="placeholder"></div>
        </div>
        <div class="copy">
          <p class="section-kicker">今日推荐</p>
          <h3 class="section-title">{{ data.random_old.name }}</h3>
          <p class="muted">还记得味道的那一道，今天可以再做。</p>
          <router-link class="btn btn-primary" :to="`/dishes/${data.random_old.id}`">看这道</router-link>
        </div>
      </article>
      <EmptyState v-else title="高评分旧菜还不多" text="再记几次，这里就会替你抽出一道推荐。" />
    </section>

    <section class="board" data-tone="moss" aria-labelledby="home-overdue-title">
      <header class="board-head">
        <div>
          <p class="section-kicker">Signature</p>
          <h2 id="home-overdue-title" class="section-title">招牌</h2>
        </div>
        <p class="board-meta">{{ data.overdue_high_rating.length }} 道高分久未做</p>
      </header>
      <DishGrid :items="data.overdue_high_rating" empty-title="暂时没有招牌" empty-text="评分 4 分以上且 30 天没做的菜会出现在这里。" />
    </section>

    <section class="board home-block" data-tone="copper" aria-labelledby="home-wish-title">
      <header class="board-head">
        <div>
          <p class="section-kicker">Want</p>
          <h2 id="home-wish-title" class="section-title">想做</h2>
        </div>
        <p class="board-meta">{{ data.want_to_cook.length }} 道</p>
      </header>
      <DishGrid :items="data.want_to_cook" empty-title="想做是空的" empty-text="刷到想复刻的菜，先记进菜牌。" />
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { fetchHomeSummary } from '../api/home'
import type { HomeSummary } from '../types'
import DishGrid from '../components/DishGrid.vue'
import EmptyState from '../components/EmptyState.vue'

const loading = ref(true)
const data = ref<HomeSummary>({
  stats: { total_dishes: 0, cooked_count: 0, success_count: 0, month_cook_count: 0 },
  recent: [],
  overdue_high_rating: [],
  random_old: null,
  want_to_cook: [],
})

onMounted(async () => {
  data.value = await fetchHomeSummary()
  loading.value = false
})
</script>
