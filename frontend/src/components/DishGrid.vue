<template>
  <div v-if="items.length" class="dish-grid">
    <DishCard v-for="dish in items" :key="dish.id" :dish="dish" />
  </div>
  <EmptyState v-else :title="emptyTitle" :text="emptyText" />
</template>

<script setup lang="ts">
import type { Dish } from '../types'
import DishCard from './DishCard.vue'
import EmptyState from './EmptyState.vue'

withDefaults(defineProps<{
  items: Dish[]
  emptyTitle?: string
  emptyText?: string
}>(), {
  emptyTitle: '还没有记录菜',
  emptyText: '先把最近想做的一道菜存进来吧。',
})
</script>
