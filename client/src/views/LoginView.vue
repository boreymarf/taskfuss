<script setup lang="ts">
import axios from 'axios';
import { reactive, ref } from 'vue';

const statusMessage = ref<string>()

const formData = reactive({
  name: <string>'',
  password: <string>'',
});

const handleSubmit = async () => {

  try {
    // const result = await authStore.register({
    //   name: formData.name,
    //   email: formData.email,
    //   password: formData.password
    // })


    const response = await axios.post('/api/login/', {
      name: formData.name,
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

  return
}

</script>

<template>
  <div>
    <h2>Введите логин и пароль</h2>
    <form @submit.prevent="handleSubmit">
      <div>
        <div>
          <label for="name">Name</label>
          <input type="text" id="name" v-model="formData.name" required>
        </div>

        <div>
          <label for="password">Password</label>
          <input type="password" id="password" v-model="formData.password" required />
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
