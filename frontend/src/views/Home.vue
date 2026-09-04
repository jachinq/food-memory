<template>
  <div v-if="loading" class="muted">正在翻开记忆...</div>
  <template v-else>
    <section class="stats">
      <div class="stat"><span>总菜品</span><b>{{ data.stats.total_dishes }}</b></div>
      <div class="stat"><span>已做</span><b>{{ data.stats.cooked_count }}</b></div>
      <div class="stat"><span>成功</span><b>{{ data.stats.success_count }}</b></div>
      <div class="stat"><span>本月制作</span><b>{{ data.stats.month_cook_count }}</b></div>
    </section>

    <article v-if="data.random_old" class="card hero-card" style="margin-bottom:24px">
      <img v-if="data.random_old.cover_image_url" :src="data.random_old.cover_image_url" :alt="data.random_old.name" />
      <div v-else class="placeholder">忆</div>
      <div class="copy">
        <p class="muted">今天想起这道菜</p>
        <h2 class="section-title">{{ data.random_old.name }}</h2>
        <p>不知道吃什么时，从记忆库里抽一道旧菜。</p>
        <router-link class="btn btn-primary" :to="`/dishes/${data.random_old.id}`">看看它</router-link>
      </div>
    </article>
    <EmptyState v-else title="高评分旧菜还不多" text="继续记录几次后我就能帮你回忆啦。" />

    <div class="home-split">
      <section>
        <h2 class="section-title">最近做过</h2>
        <DishGrid :items="data.recent" empty-title="还没有做过菜" empty-text="做完后回来记一笔，首页就会亮起来。" />
      </section>
      <section>
        <h2 class="section-title">很久没做但评分高</h2>
        <DishGrid :items="data.overdue_high_rating" empty-title="暂时没有推荐" empty-text="评分 4 分以上且 30 天没做的菜会出现在这里。" />
      </section>
    </div>

    <section>
      <h2 class="section-title">想做清单</h2>
      <DishGrid :items="data.want_to_cook" empty-title="想做清单是空的" empty-text="刷到想复刻的菜，先存进来。" />
    </section>
  </template>
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
