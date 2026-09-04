<template>
  <router-link class="card dish-card" :to="`/dishes/${dish.id}`">
    <div class="thumb">
      <img v-if="src" :src="src" :alt="dish.name" />
      <div v-else class="placeholder">菜</div>
    </div>
    <div class="body">
      <div class="row">
        <StatusBadge :status="dish.status" />
        <span v-if="dish.rating" class="rating">{{ dish.rating }}</span>
      </div>
      <h3>{{ dish.name }}</h3>
      <div class="row">
        <span v-for="tag in (dish.tags || []).slice(0, 3)" :key="tag.id" class="tag">{{ tag.name }}</span>
      </div>
      <p v-if="dish.last_cooked_at" class="muted" style="margin:8px 0 0;font-size:12px">最近 {{ formatDate(dish.last_cooked_at) }}</p>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Dish } from '../types'
import { dishCover, formatDate } from '../types'
import StatusBadge from './StatusBadge.vue'

const props = defineProps<{ dish: Dish }>()
const src = computed(() => dishCover(props.dish))
</script>
