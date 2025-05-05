<script setup lang="ts">
import axios from 'axios';
import { reactive, ref } from 'vue'

const formData = reactive({
  name: <string>'',
  email: '',
  password: '',
  confirmPassword: ''
});

const statusMessage = ref<string>()

const handleSubmit = async () => {
  statusMessage.value = ''

  if (formData.password !== formData.confirmPassword) {
    statusMessage.value = 'Пароли не совпадают';
    return;
  }

  if (!formData.email.includes('@')) {
    statusMessage.value = 'Введите корректный email';
    return;
  }

  if (formData.password.length < 6) {
    statusMessage.value = 'Пароль должен содержать минимум 6 символов';
    return;
  }

  try {
    const response = await axios.post('/api/register/', {
      name: formData.name,
      email: formData.email,
      password: formData.password
    })

    statusMessage.value = response.data

  } catch (error) {

    if (axios.isAxiosError(error)) {
      statusMessage.value = error.response?.data
      console.error(error.response?.data)
    } else {
      statusMessage.value = `Unknown error: ${error}`
      console.error(error)
    }

    return
  }
}



</script>

<template>
  <div>
    <h2>Хопа, регистрируемся тута!</h2>
    <form @submit.prevent="handleSubmit">
      <div>
        <div>
          <label for="name">Name</label>
          <input type="text" id="name" v-model="formData.name" required>
        </div>

        <div>
          <label for="email">Email</label>
          <input type="text" id="email" v-model="formData.email" required>
        </div>

        <div>
          <label for="password">Password</label>
          <input type="password" id="password" v-model="formData.password" required />
        </div>

        <div>
          <label for="confirmPassword">Confirm password</label>
          <input type="password" id="confirmPassword" v-model="formData.confirmPassword" required />
        </div>

        <div v-if="statusMessage" class="error">
          {{ statusMessage }}
        </div>
      </div>

      <button type="submit">Register</button>
    </form>
  </div>
</template>

<style lang="scss" scoped></style>
