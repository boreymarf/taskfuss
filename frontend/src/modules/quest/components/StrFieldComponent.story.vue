<script setup lang="ts">
import type { StrField } from '@/api/generated';
import StrFieldComponent from './StrFieldComponent.vue';
import { ref } from 'vue';


let example = ref<StrField>({
  discriminator: "str",
  value_type: "str",
  label: "Quest name (optional)",
  description: null,
  required: false,
  default: null,
  max_length: null,
  min_length: null
})
let value = ref<string>('')

</script>

<template>
  <Story>
    <Variant>
      <StrFieldComponent :field="example" v-model="value" />
    </Variant>
  </Story>
</template>
