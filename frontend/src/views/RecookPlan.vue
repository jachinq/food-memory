<template>
  <div>
    <p class="section-kicker">Replay</p>
    <h1 class="section-title">复做清单</h1>
    <section>
      <h2 class="section-title">计划中</h2>
      <div v-if="active.length" class="dish-grid">
        <article v-for="plan in active" :key="plan.id" class="card plan-card">
          <h3>{{ plan.dish?.name || '菜品' }}</h3>
          <p class="muted">计划日期 {{ plan.planned_date?.slice(0, 10) || '未定' }}</p>
          <div class="row">
            <input type="date" :value="plan.planned_date?.slice(0,10)" @change="onDate(plan.id, ($event.target as HTMLInputElement).value)" />
            <router-link class="btn btn-primary" :to="`/dishes/${plan.dish_id}/records/new`">去记录</router-link>
            <button class="btn btn-ghost" type="button" @click="done(plan.id)">完成</button>
            <button class="btn btn-ghost" type="button" @click="cancel(plan.id)">取消</button>
          </div>
        </article>
      </div>
      <EmptyState v-else title="还没有复做计划" text="在菜品详情点「想再做」，就会出现在这里。" />
    </section>
    <section>
      <h2 class="section-title">本周建议</h2>
      <DishGrid :items="suggest" empty-title="暂时没有建议" empty-text="高评分且很久没做的菜会作为本周建议。" />
    </section>
    <section>
      <h2 class="section-title">已完成</h2>
      <EmptyState v-if="!completed.length" title="还没有完成的计划" text="做完后点完成，并顺手记一次制作。" />
      <ul v-else class="completed-list">
        <li v-for="plan in completed" :key="plan.id">{{ plan.dish?.name }} · {{ plan.completed_at?.slice(0,10) }}</li>
      </ul>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { cancelRecookPlan, completeRecookPlan, fetchRecookPlans, updateRecookPlan } from '../api/recook'
import { fetchHomeSummary } from '../api/home'
import type { Dish, RecookPlan } from '../types'
import DishGrid from '../components/DishGrid.vue'
import EmptyState from '../components/EmptyState.vue'

const plans = ref<RecookPlan[]>([])
const suggest = ref<Dish[]>([])
const active = computed(() => plans.value.filter((p) => p.status === 'active'))
const completed = computed(() => plans.value.filter((p) => p.status === 'completed'))

async function load() {
  plans.value = await fetchRecookPlans()
  const home = await fetchHomeSummary()
  suggest.value = home.overdue_high_rating
}

onMounted(load)
async function done(id: number) { await completeRecookPlan(id); await load() }
async function cancel(id: number) { await cancelRecookPlan(id); await load() }
async function onDate(id: number, date: string) { await updateRecookPlan(id, { planned_date: date }); await load() }
</script>
