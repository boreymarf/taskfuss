<script setup lang="ts">
import { computed } from 'vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)

</script>

<template>
  <div class="header">
    <div class="logo">TASKFUSS</div>
    <div class="nav">
      <router-link v-if="isAuthenticated" to="/dashboard">Dashboard</router-link>
      <router-link v-if="!isAuthenticated" to="/login">Login</router-link>
      <router-link v-if="!isAuthenticated" to="/register">Register</router-link>
      <router-link to="/about">About</router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.header {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 1rem;
  border-bottom: 1px solid #ddd;

  .logo {
    margin-right: auto;
  }

  .nav {
    display: flex;
    gap: 1.5rem;

    a {
      text-decoration: none;
      color: inherit;
    }
  }
}
</style>
