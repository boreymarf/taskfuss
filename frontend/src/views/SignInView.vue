<script setup lang="ts">
import type { UserLogin } from '@/api/generated'
import { SignInForm, useSignIn } from '@/modules/auth'
import { ref } from 'vue'

const { mutateAsync: signIn, isPending } = useSignIn()
const errorMessage = ref<string | null>(null)

const onFormSubmit = async (values: UserLogin) => {
  errorMessage.value = null

  try {
    const token = await signIn(values)
    console.log('Успех:', token)
  } catch (error: any) {
    errorMessage.value = error?.detail || error?.message || 'Sign in failed'
  }
}
</script>

<template>
  <div class="flex-1 flex items-center justify-center">
    <div class="flex flex-col w-96">
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 text-red-700 rounded">
        {{ errorMessage }}
      </div>
      <SignInForm @submit="onFormSubmit" :isPending="isPending" />
    </div>
  </div>
</template>
