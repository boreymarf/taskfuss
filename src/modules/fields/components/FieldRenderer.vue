<script setup lang="ts">
import { computed, type Component } from 'vue';
import type { PlanField } from '@generated/bindings';
import StrField from './StrField.vue';
import IntField from './IntField.vue';

// Stop Vue from passing attrs automatically
defineOptions({ inheritAttrs: false });

const props = defineProps<{ field: PlanField }>();
const model = defineModel<unknown>();

const map: Record<PlanField['type'], Component | null> = {
  str: StrField,
  int: IntField,
  str_list: null,
  progress: null,
  parent: null,
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
    v-bind="{ ...$attrs, ...fieldProps }"
  />
  <div v-else>
    No component for schema "{{ field.type }}"
  </div>
</template>