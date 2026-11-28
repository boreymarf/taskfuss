<script setup lang="ts">
import { computed, type Component } from 'vue';
import type { Field, FieldValueError } from '@generated/bindings';
import StrField from './StrField.vue';
import IntField from './IntField.vue';

// Stop Vue from passing attrs automatically
defineOptions({ inheritAttrs: false });

const props = defineProps<{ field: Field; error?: FieldValueError }>();
const model = defineModel<unknown>();

const map: Record<Field['type'], Component | null> = {
  str: StrField,
  int: IntField,
};

const component = computed(() => map[props.field.type]);

// Remove the 'type' field from props
const fieldProps = computed(() => {
  const { type: _, ...rest } = props.field;
  return rest;
});
</script>

<template>
  <component
    :is="component"
    v-if="component"
    v-model="model"
    v-bind="$attrs"
    :field="fieldProps"
    :error="props.error"
  />
  <div v-else>
    No component for schema "{{ field.type }}"
  </div>
</template>