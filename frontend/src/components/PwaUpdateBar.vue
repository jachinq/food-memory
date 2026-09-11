<template>
  <div v-if="needRefresh" class="pwa-update" role="status">
    <span>食忆有更新，刷新后使用</span>
    <button class="btn btn-primary" type="button" @click="updateServiceWorker()">刷新</button>
  </div>
</template>

<script setup lang="ts">
import { useRegisterSW } from 'virtual:pwa-register/vue'

const { needRefresh, updateServiceWorker } = useRegisterSW({ immediate: true })
</script>
