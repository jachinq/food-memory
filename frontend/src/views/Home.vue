<template>
  <div v-if="loading" class="skeleton" aria-hidden="true"></div>
  <div v-else class="home-page">
    <section class="ledger" aria-label="记忆概览">
      <div class="stat"><span>总菜品</span><b>{{ data.stats.total_dishes }}</b></div>
      <div class="stat"><span>已做</span><b>{{ data.stats.cooked_count }}</b></div>
      <div class="stat"><span>成功</span><b>{{ data.stats.success_count }}</b></div>
      <div class="stat"><span>本月</span><b>{{ data.stats.month_cook_count }}</b></div>
    </section>

    <section class="board" data-tone="chili" aria-labelledby="home-draw-title">
      <header class="board-head">
        <span class="board-mark" aria-hidden="true">壹</span>
        <div>
          <p class="section-kicker">Draw</p>
          <h2 id="home-draw-title" class="section-title">今日抽菜</h2>
        </div>
        <p class="board-meta">不知道吃什么时翻一页</p>
      </header>
      <article v-if="data.random_old" class="card hero-card">
        <div class="visual">
          <span class="stamp">今日抽菜</span>
          <img v-if="data.random_old.cover_image_url" :src="data.random_old.cover_image_url" :alt="data.random_old.name" />
          <div v-else class="placeholder">忆</div>
        </div>
        <div class="copy">
          <p class="section-kicker">从旧账抽出</p>
          <h3 class="section-title">{{ data.random_old.name }}</h3>
          <p class="muted">还记得味道的那一道，今天可以再做。</p>
          <router-link class="btn btn-primary" :to="`/dishes/${data.random_old.id}`">翻开这页</router-link>
        </div>
      </article>
      <EmptyState v-else title="高评分旧菜还不多" text="再记几次，这里就会替你翻出一道旧菜。" />
    </section>

    <div class="home-split">
      <section class="board" data-tone="ink" aria-labelledby="home-recent-title">
        <header class="board-head">
          <span class="board-mark" aria-hidden="true">贰</span>
          <div>
            <p class="section-kicker">Recent</p>
            <h2 id="home-recent-title" class="section-title">最近做过</h2>
          </div>
          <p class="board-meta">{{ data.recent.length }} 道</p>
        </header>
        <DishGrid :items="data.recent" empty-title="还没有做过菜" empty-text="做完后回来记一笔，首页就会亮起来。" />
      </section>
      <section class="board" data-tone="moss" aria-labelledby="home-overdue-title">
        <header class="board-head">
          <span class="board-mark" aria-hidden="true">叁</span>
          <div>
            <p class="section-kicker">Overdue</p>
            <h2 id="home-overdue-title" class="section-title">很久没做</h2>
          </div>
          <p class="board-meta">{{ data.overdue_high_rating.length }} 道</p>
        </header>
        <DishGrid :items="data.overdue_high_rating" empty-title="暂时没有推荐" empty-text="评分 4 分以上且 30 天没做的菜会出现在这里。" />
      </section>
    </div>

    <section class="board home-block" data-tone="copper" aria-labelledby="home-wish-title">
      <header class="board-head">
        <span class="board-mark" aria-hidden="true">肆</span>
        <div>
          <p class="section-kicker">Wishlist</p>
          <h2 id="home-wish-title" class="section-title">想做清单</h2>
        </div>
        <p class="board-meta">{{ data.want_to_cook.length }} 道</p>
      </header>
      <DishGrid :items="data.want_to_cook" empty-title="想做清单是空的" empty-text="刷到想复刻的菜，先夹进这本账里。" />
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
