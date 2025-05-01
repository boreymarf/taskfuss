<script setup lang="ts">
import { ref } from 'vue';
import { API_CONFIG } from '../api/config';

const username = ref<string>('');
const password = ref<string>('');
const error = ref<string>('');

const handleLogin = async () => {
  try {
    const response = await fetch(`${API_CONFIG.SERVER_URL}/api/auth`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value
      })
    });

    const data = await response.json();

    if (!response.ok) {
      // Handle server-sent error messages
      throw new Error(data.message || 'Login failed');
    }

    localStorage.setItem('token', data.token);

    error.value = "Success"
    //window.location.href = '/dashboard'; // Redirect after login
  } catch (err) {

    if (err instanceof Error) {
      error.value = err.message;
    } else if (typeof err === 'string') {
      error.value = err;
    } else {
      error.value = 'An unknown error occurred';
    }

    password.value = '';
  }
};
</script>

<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <div>
        <input placeholder="username" type="text" v-model="username" required>
      </div>
      <div>
        <input placeholder="password" type="password" v-model="password" required>
      </div>
      <button type="submit">Login</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style lang="scss"></style>
