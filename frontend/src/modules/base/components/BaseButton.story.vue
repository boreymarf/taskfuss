<script setup lang="ts">
import BaseButton from './BaseButton.vue';
import { logEvent } from 'histoire/client'
</script>

<template>
  <Story :layout="{ type: 'grid', width: 400 }">
    <Variant title="normal">
      <BaseButton @click="logEvent('Click', $event)"> Button </BaseButton>
    </Variant>
    <Variant title="disabled">
      <BaseButton :disabled="true" @click="logEvent('Click', $event)"> Disabled </BaseButton>
    </Variant>
  </Story>
</template>
