<script setup lang="ts">
import { RouterView } from 'vue-router'
import MyHeader from './components/MyHeader.vue';
</script>

<template>
  <main>
    <MyHeader class="header" />
    <RouterView class="main" />
  </main>
</template>

<style lang="scss" scoped>
main {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  position: sticky;
  top: 0;
  z-index: 10;
}
</style>
