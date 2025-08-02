<script setup lang="ts">
import { RouterView } from 'vue-router'
import MyHeader from './components/MyHeader.vue';
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth';

const authStore = useAuthStore()
onMounted(async () => {
  if (authStore.auth_token && !authStore.user) {
    try {
      await authStore.fetchUser()
    } catch (error) {
      console.error('Failed to fetch user:', error)
    }
  }
})
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
