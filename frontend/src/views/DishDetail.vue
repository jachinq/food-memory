<template>
  <div v-if="dish">
    <div class="toolbar">
      <router-link class="muted" to="/dishes">← 回到菜牌</router-link>
      <div style="margin-left:auto" class="row">
        <router-link class="btn btn-ghost" :to="`/dishes/${dish.id}/edit`">编辑</router-link>
        <router-link class="btn btn-primary" :to="`/dishes/${dish.id}/records/new`">记下一次</router-link>
        <button class="btn btn-ghost" type="button" @click="markRecook">待做</button>
        <button class="btn btn-danger" type="button" @click="remove">删除</button>
      </div>
    </div>
    <section class="card hero-card">
      <div class="visual">
        <img v-if="dish.cover_image_url" :src="dish.cover_image_url" :alt="dish.name" @click="preview = dish.cover_image_url" />
        <div v-else class="placeholder"></div>
      </div>
      <div class="copy">
        <div class="row">
          <StatusBadge :status="dish.status" />
          <span v-if="dish.rating" class="rating">{{ dish.rating }} 分</span>
        </div>
        <h1 class="section-title">{{ dish.name }}</h1>
        <p v-if="dish.difficulty" class="muted">难度 {{ dish.difficulty }}</p>
        <p v-if="dish.description" class="dish-method">{{ dish.description }}</p>
        <p class="muted">做过 {{ dish.cook_count }} 次 · 最近 {{ formatDate(dish.last_cooked_at) || '还没有' }}</p>
        <div class="row" style="margin-top:8px">
          <span v-for="tag in menuSectionTags(dish.tags)" :key="tag.id" class="tag">{{ tag.name }}</span>
        </div>
      </div>
    </section>
    <p class="section-kicker">Served</p>
    <h2 class="section-title">过往出品</h2>
    <EmptyState v-if="!dish.records?.length" title="这道菜还没有制作记录" text="做完后记得回来补一条。" />
    <CookRecordTimeline
      v-else
      :items="dish.records"
      :dish-id="dish.id"
      @preview="preview = $event"
      @remove="removeRecord"
    />
    <p v-if="dish.source_url" class="source-footer">
      <a :href="dish.source_url" target="_blank" rel="noreferrer">来源 {{ dish.source_platform || '链接' }}</a>
    </p>
    <div v-if="preview" class="lightbox" @click="preview = ''">
      <img :src="preview" alt="" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { deleteDish, getDish } from '../api/dishes'
import { deleteRecord } from '../api/records'
import { createRecookPlan } from '../api/recook'
import { ApiError } from '../api/client'
import type { Dish } from '../types'
import { menuSectionTags, formatDate } from '../types'
import StatusBadge from '../components/StatusBadge.vue'
import CookRecordTimeline from '../components/CookRecordTimeline.vue'
import EmptyState from '../components/EmptyState.vue'

const route = useRoute()
const router = useRouter()
const dish = ref<Dish | null>(null)
const preview = ref('')

onMounted(async () => {
  dish.value = await getDish(Number(route.params.id))
})

async function markRecook() {
  if (!dish.value) return
  try {
    await createRecookPlan({ dish_id: dish.value.id })
    alert('已加入待做清单')
    dish.value = await getDish(dish.value.id)
  } catch (e) {
    if (e instanceof ApiError && e.message === '已在清单') {
      alert('已在清单')
      return
    }
    alert(e instanceof Error ? e.message : '加入失败')
  }
}

async function removeRecord(id: number) {
  if (!dish.value) return
  if (!confirm('删除这条制作记录？次数和评分会按剩下的记录重算。整道菜不会被删。')) return
  try {
    await deleteRecord(id)
    dish.value = await getDish(dish.value.id)
  } catch (e) {
    alert(e instanceof Error ? e.message : '删除失败')
  }
}

async function remove() {
  if (!dish.value) return
  if (!confirm('删除后菜品和制作记录都会进入回收状态，确定吗？')) return
  await deleteDish(dish.value.id)
  router.push('/dishes')
}
</script>
