<template>
  <router-link class="dish-card" :to="`/dishes/${dish.id}`" :style="{ '--i': index }">
    <div class="thumb">
      <img v-if="src" :src="src" :alt="dish.name" />
      <div v-else class="placeholder"></div>
    </div>
    <div class="body">
      <h3>{{ dish.name }}</h3>
      <div class="row">
        <StatusBadge :status="dish.status" />
        <span v-if="dish.rating" class="rating">{{ dish.rating }}</span>
      </div>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Dish } from '../types'
import { dishCover } from '../types'
import StatusBadge from './StatusBadge.vue'

const props = defineProps<{ dish: Dish; index?: number }>()
const src = computed(() => dishCover(props.dish))
</script>
