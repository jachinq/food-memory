<template>
  <span class="badge" :class="status">{{ STATUS_LABEL[status] || status }}</span>
</template>

<script setup lang="ts">
import { STATUS_LABEL } from '../types'
defineProps<{ status: string }>()
</script>
