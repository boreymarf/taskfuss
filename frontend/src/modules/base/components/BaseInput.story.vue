<script setup lang="ts">
import BaseInput from './BaseInput.vue';
import { logEvent } from 'histoire/client'
import { ref } from 'vue'

const value = ref('')
const text = ref('Hello')
const number = ref(42)
const password = ref('')
</script>

<template>
  <Story :layout="{ type: 'grid', width: 400 }">
    <Variant title="text">
      <BaseInput v-model="value" placeholder="Enter text..." @input="logEvent('Input', $event)"
        @focus="logEvent('Focus', $event)" @blur="logEvent('Blur', $event)" />
    </Variant>

    <Variant title="with default value">
      <BaseInput v-model="text" />
    </Variant>

    <Variant title="number">
      <BaseInput v-model="number" type="number" />
    </Variant>

    <Variant title="password">
      <BaseInput v-model="password" type="password" placeholder="Enter password..." />
    </Variant>

    <Variant title="disabled">
      <BaseInput :disabled="true" value="Can't touch this" />
    </Variant>

    <Variant title="readonly">
      <BaseInput :readonly="true" value="Read only text" />
    </Variant>

    <Variant title="with placeholder">
      <BaseInput placeholder="Type something here..." />
    </Variant>

    <Variant title="email">
      <BaseInput type="email" placeholder="email@example.com" />
    </Variant>
  </Story>
</template>
