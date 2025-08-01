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
      <router-link v-if="isAuthenticated" to="/tasks" class="option">Tasks</router-link>
      <router-link to="/about" class="option">About</router-link>
      <router-link v-if="!isAuthenticated" to="/login" class="option">Login</router-link>
      <router-link v-if="!isAuthenticated" to="/register" class="option">Register</router-link>
      <div v-if="isAuthenticated" class="username">{{ authStore.$state.user?.username || "Unknown" }}</div>
      <button v-if="isAuthenticated" v-on:click="authStore.logout" class="option">Logout</button>
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

  .option {
    background-color: transparent;
    padding: 8px 8px;
    border: 1px solid transparent;
    transition: border-color 0.25s;
    border-radius: 8px;

    &:hover {
      border-color: #646cff;
    }
  }

  .username {
    background-color: transparent;
    padding: 8px 8px;
    border: 1px solid transparent;
    transition: border-color 0.25s;
  }
}
</style>
