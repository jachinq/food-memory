<template>
  <div class="install-actions">
    <button
      v-if="install.surface === 'install'"
      class="btn btn-primary"
      type="button"
      @click="install.install()"
    >
      安装食忆
    </button>
    <p v-else-if="install.surface === 'ios-hint'" class="muted">
      在 Safari 点底部分享，再选「添加到主屏幕」。
    </p>
    <p v-else-if="install.surface === 'installed'" class="muted">已在本机打开。</p>
  </div>
</template>

<script setup lang="ts">
import { useInstall } from '../stores/install'

const install = useInstall()
</script>
