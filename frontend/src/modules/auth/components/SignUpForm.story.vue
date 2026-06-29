<script setup lang="ts">
import SignUpForm from './SignUpForm.vue'
import { logEvent } from 'histoire/client'
</script>

<template>
  <Story :layout="{ type: 'single' }">
    <Variant title="default">
      <SignUpForm @submit="logEvent('Submit', $event)" />
    </Variant>
  </Story>
</template>
