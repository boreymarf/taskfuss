<script setup lang="ts">
import type { StrField } from '@/api/generated';
import StrFieldComponent from './StrFieldComponent.vue';


const example: StrField = {
  discriminator: "str",
  value_type: "str",
  label: "Quest name (optional)",
  description: null,
  required: false,
  default: null,
  max_length: null,
  min_length: null
}

</script>

<template>
  <Story>
    <Variant>
      <StrFieldComponent v-bind="example" />
    </Variant>
  </Story>
</template>
