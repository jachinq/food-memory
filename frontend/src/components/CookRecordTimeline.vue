<template>
  <div class="timeline">
    <article v-for="item in items" :key="item.id" class="timeline-item">
      <div>
        <img
          v-if="item.photos?.[0]"
          :src="item.photos[0].thumbnail_url || item.photos[0].file_url"
          alt=""
          style="width:100%;height:72px;object-fit:cover;border-radius:10px;cursor:pointer"
          @click="$emit('preview', item.photos![0].file_url)"
        />
        <div v-else class="placeholder" style="height:72px;border-radius:10px">记</div>
      </div>
      <div>
        <div class="row">
          <strong>{{ formatDate(item.cooked_at) }}</strong>
          <span class="badge" :class="item.result">{{ RESULT_LABEL[item.result] || item.result }}</span>
          <span v-if="item.rating" class="rating">{{ item.rating }}</span>
        </div>
        <p v-if="item.notes" class="muted">{{ item.notes }}</p>
        <p v-if="item.next_improvement"><b>下次注意：</b>{{ item.next_improvement }}</p>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import type { CookRecord } from '../types'
import { RESULT_LABEL, formatDate } from '../types'
defineProps<{ items: CookRecord[] }>()
defineEmits<{ preview: [string] }>()
</script>
