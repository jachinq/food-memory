<template>
  <div class="me-page">
    <p class="section-kicker">Me</p>
    <h1 class="section-title">我的</h1>
    <p class="muted">菜品和图片需要联网。安装只把食忆钉到本机，不是另一套离线菜牌。</p>
    <InstallActions />
  </div>
</template>

<script setup lang="ts">
import InstallActions from '../components/InstallActions.vue'
</script>
