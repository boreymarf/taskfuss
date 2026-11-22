<script setup lang="ts">
import { reactive, watch } from 'vue';
import type { PlanField } from '@generated/bindings';
import FieldRenderer from './FieldRenderer.vue';

const props = defineProps<{
  fields: Record<string, PlanField>;
}>();

const model = defineModel<Record<string, unknown>>({ default: () => ({}) });

// 1) Dictionary: { [fieldKey]: null } — seeded from any existing model values
const values = reactive<Record<string, unknown>>(
  Object.fromEntries(
    Object.keys(props.fields).map((k) => [k, model.value?.[k] ?? null]),
  ),
);

// Keep dictionary keys in sync if `fields` changes (keys added / removed)
watch(
  () => props.fields,
  (fields) => {
    for (const key of Object.keys(fields)) {
      if (!(key in values)) values[key] = null;
    }
    for (const key of Object.keys(values)) {
      if (!(key in fields)) delete values[key];
    }
  },
);

// 2) Push the whole dictionary up whenever any field updates
watch(
  values,
  (v) => {
    model.value = { ...v };
  },
  { deep: true },
);
</script>

<template>
  <div class="field-form">
    <FieldRenderer
      v-for="(field, key) in props.fields"
      :key="key"
      v-model="values[key]"
      :field="field"
    />
  </div>
</template>

<style scoped lang="scss">
.field-form {
  display: flex;
  flex-direction: column;
  width: 100%;
  gap: 0.75rem;

  > * {
    width: 100%;
  }
}
</style>