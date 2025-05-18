<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '../stores/auth';
import type { RegisterForm } from '../types/forms';
import { isValidEmail, isValidPassword } from '../utils/validation';
import { RegisterError } from '../types/stores';

const authStore = useAuthStore()
const statusMessage = ref<string>()

const formData = reactive<RegisterForm>({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
});


const handleSubmit = async () => {
  statusMessage.value = ''

  // Client sided validation
  if (!isValidEmail(formData.email)) {
    statusMessage.value = "Please enter valid email"
    return
  }

  if (!isValidPassword(formData.password)) {
    statusMessage.value = "Please enter valid password"
    return
  }

  if (formData.password !== formData.confirmPassword) {
    statusMessage.value = "Passwords do not match"
    return
  }

  try {

    await authStore.register(formData)
    statusMessage.value = "Success!"

  } catch (error) {

    // TODO: Make this error field look better
    if (error instanceof RegisterError) {
      const messages = [
        error.message,
        ...error.details.map(e => e.message)
      ];

      statusMessage.value = messages.join('\n');

    } else if (error instanceof Error) {
      statusMessage.value = error.message
    } else {
      statusMessage.value = "Unknown error"
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

        <div v-if="statusMessage" class="status-message">
          {{ statusMessage }}
        </div>
      </div>

      <button type="submit">Register</button>
    </form>
  </div>
</template>

<style lang="scss" scoped>
.status-message {
  white-space: pre-line
}
</style>
