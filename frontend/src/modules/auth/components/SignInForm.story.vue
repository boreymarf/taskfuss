<script setup lang="ts">
import SignInForm from "./SignInForm.vue";
import { logEvent } from "histoire/client";
</script>

<template>
  <Story :layout="{ type: 'single' }">
    <Variant title="default">
      <SignInForm @submit="logEvent('Submit', $event)" />
    </Variant>
  </Story>
</template>
