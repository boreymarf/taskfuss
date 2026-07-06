<script setup lang="ts">
import { ref, watch } from "vue";

interface Props {
  modelValue?: boolean;
  disabled?: boolean;
  indeterminate?: boolean;
  value?: any;
  name?: string;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  disabled: false,
  indeterminate: false,
  value: undefined,
  name: undefined,
});

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "change", event: Event): void;
}>();

const inputRef = ref<HTMLInputElement>();

watch(
  () => props.indeterminate,
  (val) => {
    if (inputRef.value) {
      inputRef.value.indeterminate = val;
    }
  },
  { immediate: true },
);

function handleChange(event: Event) {
  const target = event.target as HTMLInputElement;
  emit("update:modelValue", target.checked);
  emit("change", event);
}
</script>

<template>
  <input
    ref="inputRef"
    type="checkbox"
    :checked="modelValue"
    :disabled="disabled"
    :value="value"
    :name="name"
    @change="handleChange"
  />
</template>
