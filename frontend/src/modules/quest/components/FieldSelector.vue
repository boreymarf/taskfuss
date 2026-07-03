<script setup lang="ts">
import type { Ref } from 'vue'
import type { Field } from '../types/fields'
import StrFieldComponent from './StrFieldComponent.vue';

defineProps<{
  fieldRef: Ref<any>
  field: Field
}>()
</script>

<template>
  <StrFieldComponent v-if="field.discriminator === 'str'" :field-ref="fieldRef" :field="field" />
  <p v-else>Unknown field type!</p>
</template>
