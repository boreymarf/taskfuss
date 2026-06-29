<script setup lang="ts">
import { useAuthStore } from '@/modules/auth';
import BaseButton from '@/modules/base/components/BaseButton.vue';
import { RouterLink } from 'vue-router';

const authStore = useAuthStore()
</script>

<template>
  <nav class="flex justify-between items-center bg-green-950 text-white p-4">
    <div class="flex gap-4 items-center">
      <RouterLink to="/" class="hover:underline">Front</RouterLink>
    </div>

    <div class="flex gap-4 items-center">
      <div v-if="authStore.isAuthenticated">
        <BaseButton @click="authStore.logout">Logout</BaseButton>
      </div>
      <div v-else class="flex gap-4 items-center">
        <RouterLink to="/sign-in" class="hover:underline">Sign In</RouterLink>
        <RouterLink to="/sign-up" class="hover:underline">Sign Up</RouterLink>
      </div>
    </div>
  </nav>
</template>
