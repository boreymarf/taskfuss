<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios' // или используй fetch

const status = ref<string | null>(null)
const loading = ref<boolean>(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:4000/api/status')
    status.value = response.data.isUp
  } catch (err) {
    error.value = `Ошибка при загрузке статуса ${err}`
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="status-container">
    <div v-if="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <h2>Статус сервиса:</h2>
      <div :class="['status', status ? 'up' : 'down']">
        {{ status ? 'Всё работает ✅' : 'Сервис недоступен ❌' }}
      </div>
    </div>
  </div>
</template>
