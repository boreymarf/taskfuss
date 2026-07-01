<script setup lang="ts">
import type { Field } from '../types/fields';
import StrFieldComponent from './StrFieldComponent.vue';
import ListFieldComponent from './ListFieldComponent.vue';

const props = defineProps<{
  field: Field
  modelValue: any
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>();
</script>

<template>
  <StrFieldComponent
    v-if="field.discriminator === 'str'"
    :field="field"
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
  />
  <ListFieldComponent
    v-else-if="field.discriminator === 'list'"
    :field="field"
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
  />
  <div v-else>Unknown field!</div>
</template>
