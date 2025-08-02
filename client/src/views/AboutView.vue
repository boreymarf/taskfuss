<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../api/client/http'
import type { DtoPongResponse } from '../api/generated'

const response = ref<string>('')
const isLoading = ref(false)

const pingServer = async () => {
  isLoading.value = true
  response.value = ""

  try {
    const result = await api.service.pingGet()
    const responseData = result.data as DtoPongResponse
    response.value = responseData.message || "No message found!"

  } catch (error) {
    console.error('API Error:', error)
    response.value = 'Error connecting to server'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="main">
    <h1>This is an about page.</h1>
    <button @click="pingServer" :disabled="isLoading">
      {{ isLoading ? 'Loading...' : 'Test Connection' }}
    </button>
    <p v-if="response">{{ response }}</p>
  </div>
</template>
