<script setup lang="ts">
import type { Field } from '../types/fields';
import StrFieldComponent from './StrFieldComponent.vue';
import ListFieldComponent from './ListFieldComponent.vue';
import { computed } from 'vue';

const props = defineProps<{
  field: Field
  modelValue: any
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>();

const value = computed({
  get: () => props.modelValue,
  set: (newValue) => emit('update:modelValue', newValue)
});
</script>

<template>
  <StrFieldComponent
    v-if="field.discriminator === 'str'"
    :field="field"
    v-model="value"
  />
  <ListFieldComponent
    v-else-if="field.discriminator === 'list'"
    :field="field"
    v-model="value"
  />
  <div v-else>Unknown field!</div>
</template>
