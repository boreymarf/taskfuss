<script setup lang="ts">

import { ref, onMounted } from 'vue';
import http from '@/api/http';

const loading = ref(true);
const error = ref('');
const secret = ref<number | null>(null);
const userEmail = ref('');

</script>

<template>
  <div v-if="loading">Loading...</div>
  <div v-else-if="error" class="error">{{ error }}</div>
  <div v-else>
    <h2>Welcome, {{ userEmail }}!</h2>
    <p>Your secret number: {{ secret }}</p>
  </div>
</template>

<style lang="scss" scoped></style>
