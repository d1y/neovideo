<template>
  <div class="portal">
    <Header />
    <Menu />
    <router-view />
  </div>
  <Footer />
</template>

<script setup lang="ts">
import Header from '@/views/components/header.vue'
import Footer from '@/views/components/footer.vue'
import Menu from '@/views/components/menu.vue'
</script>

<style scoped lang="less">
.portal {
  width: 100%;
  max-width: 1024px;
  min-height: calc(100vh - 180px);
  margin: 0 auto;
  padding: 0 10px 16px;
}

body {
  font-family: -apple-system, system-ui, BlinkMacSystemFont, Segoe UI, Helvetica Neue, PingFang SC, Helvetica, Microsoft Yahei, Arial, sans-serif;
}
</style>
