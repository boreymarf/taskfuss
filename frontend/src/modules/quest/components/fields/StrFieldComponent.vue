<script setup lang="ts">
import { watch } from "vue";
import BaseInput from "@/modules/base/components/BaseInput.vue";
import type { StrField } from "../../types/fields";

const props = defineProps<{
  field: StrField;
}>();

const model = defineModel<string>();

// We need to use watch to set the default value because
// list field component needs to pass a ref first
watch(
  model,
  (newVal) => {
    if (newVal === null || newVal === undefined) {
      model.value = props.field.default ?? "";
    }
  },
  { immediate: true },
);
</script>

<template>
  <div class="w-full min-h-12.5 flex">
    <div v-if="field.label" class="w-[30%] flex items-center px-2">
      {{ field.label }}
    </div>
    <BaseInput v-model="model" class="flex-1" />
  </div>
</template>
