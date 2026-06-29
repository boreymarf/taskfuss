<script setup lang="ts">
import BaseButton from '@/modules/base/components/BaseButton.vue'
import BaseInput from '@/modules/base/components/BaseInput.vue'
import { ref } from 'vue'

export type UserLogin = {
  login: string
  password: string
}

const formData = ref<UserLogin>({
  login: '',
  password: ''
})

const emit = defineEmits<{
  (e: 'submit', data: UserLogin): void
}>()

const handleSubmit = () => {
  emit('submit', formData.value)
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="signin-form">
    <BaseInput
      v-model="formData.login"
      placeholder="Login"
      autofocus
    />

    <BaseInput
      v-model="formData.password"
      type="password"
      placeholder="Password"
    />

    <BaseButton type="submit">
      Sign In
    </BaseButton>
  </form>
</template>

<style scoped>
.signin-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
