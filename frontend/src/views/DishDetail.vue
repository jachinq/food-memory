<template>
  <div v-if="dish">
    <div class="toolbar">
      <router-link to="/dishes">返回</router-link>
      <div style="margin-left:auto" class="row">
        <router-link class="btn btn-ghost" :to="`/dishes/${dish.id}/edit`">编辑</router-link>
        <router-link class="btn btn-primary" :to="`/dishes/${dish.id}/records/new`">记录制作</router-link>
        <button class="btn btn-ghost" type="button" @click="markRecook">想再做</button>
        <button class="btn btn-danger" type="button" @click="remove">删除</button>
      </div>
    </div>
    <section class="card hero-card" style="margin-bottom:20px">
      <img v-if="dish.cover_image_url" :src="dish.cover_image_url" :alt="dish.name" @click="preview = dish.cover_image_url" />
      <div v-else class="placeholder">菜</div>
      <div class="copy">
        <div class="row">
          <StatusBadge :status="dish.status" />
          <span v-if="dish.rating" class="rating">{{ dish.rating }} 分</span>
        </div>
        <h1 class="section-title">{{ dish.name }}</h1>
        <p class="muted">做过 {{ dish.cook_count }} 次 · 最近 {{ formatDate(dish.last_cooked_at) || '还没有' }}</p>
        <div class="row" style="margin-top:8px">
          <span v-for="tag in dish.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
        </div>
        <p v-if="dish.source_url"><a :href="dish.source_url" target="_blank">来源链接</a> {{ dish.source_platform }}</p>
        <p v-if="dish.note">{{ dish.note }}</p>
      </div>
    </section>
    <h2 class="section-title">制作记录</h2>
    <EmptyState v-if="!dish.records?.length" title="这道菜还没有制作记录" text="做完后记得回来补一条。" />
    <CookRecordTimeline v-else :items="dish.records" @preview="preview = $event" />
    <div v-if="preview" class="lightbox" @click="preview = ''">
      <img :src="preview" alt="" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { deleteDish, getDish } from '../api/dishes'
import { createRecookPlan } from '../api/recook'
import type { Dish } from '../types'
import { formatDate } from '../types'
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
  await createRecookPlan({ dish_id: dish.value.id })
  alert('已加入复做清单')
  dish.value = await getDish(dish.value.id)
}

async function remove() {
  if (!dish.value) return
  if (!confirm('删除后菜品和制作记录都会进入回收状态，确定吗？')) return
  await deleteDish(dish.value.id)
  router.push('/dishes')
}
</script>
