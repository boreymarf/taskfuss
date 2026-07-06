<script setup lang="ts">
import { reactive, watch } from "vue";
import FieldSelector from "./FieldSelector.vue";
import type { Field } from "../types/fields";

const props = defineProps<{
  fields: Record<string, Field>;
}>();

const model = defineModel<Record<string, any>>({ default: () => ({}) });

// We fill formData with nulls, so we can reference them in children
const formData = reactive<Record<string, any>>({});
for (const key of Object.keys(props.fields)) {
  formData[key] = model.value?.[key] ?? null;
}

watch(
  formData,
  () => {
    model.value = { ...formData };
  },
  { deep: true, flush: "post" },
);
</script>

<template>
  <FieldSelector
    v-for="(field, key) in fields"
    :key="key"
    v-model="formData[key]"
    :field="field"
  />
</template>
