<script setup lang="ts">
import type { UserCreate } from '@/api/generated'
import SignUpForm from '@/modules/auth/components/SignUpForm.vue'
import { useSignUp } from '@/modules/auth/mutations/useSignUp'
import { ref } from 'vue'

const { mutateAsync: registerUser, isPending } = useSignUp()
const errorMessage = ref<string | null>(null)

const onFormSubmit = async (values: UserCreate) => {
  errorMessage.value = null

  try {
    const token = await registerUser(values)
    console.log('Успех:', token)
  } catch (error: any) {
    errorMessage.value = error?.detail || error?.message || 'Registration failed'
  }
}
</script>

<template>
  <div class="flex-1 flex items-center justify-center">
    <div class="flex flex-col w-96">
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 text-red-700 rounded">
        {{ errorMessage }}
      </div>
      <SignUpForm @submit="onFormSubmit" :isPending="isPending" />
    </div>
  </div>
</template>
