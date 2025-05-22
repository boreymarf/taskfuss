<script setup lang="ts">
import axios from 'axios';
import { reactive, ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import type { LoginForm, RegisterForm } from '../types/forms';
import { StoreGenericError } from '../types/stores';

const statusMessage = ref<string>()
const authStore = useAuthStore()

const formData = reactive<LoginForm>({
  email: '',
  password: '',
});

const handleSubmit = async () => {

  try {

    await authStore.login(formData)

  } catch (error) {

    if (error instanceof StoreGenericError) { 

      statusMessage.value = error.message

    }   

  }

}


</script>

<template>
  <div>
    <h2>Login page</h2>
    <form @submit.prevent="handleSubmit">
      <div>
        <div>
          <label for="name">Email</label>
          <input type="text" id="name" v-model="formData.email" required>
        </div>

        <div>
          <label for="password">Password</label>
          <input type="password" id="password" v-model="formData.password" required />
        </div>

        <div v-if="statusMessage" class="error">
          {{ statusMessage }}
        </div>
      </div>

      <button type="submit">Login</button>
    </form>
  </div>
</template>

<style lang="scss" scoped></style>
