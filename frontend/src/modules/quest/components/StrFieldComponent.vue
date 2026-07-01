<script setup lang="ts">
import BaseInput from '@/modules/base/components/BaseInput.vue';
import type { StrField } from '../types/fields';
import { computed } from 'vue';

const props = defineProps<{
  field: StrField
  modelValue: string | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const safeValue = computed(() => props.modelValue ?? '')

function updateValue(value: string | number) {
  emit('update:modelValue', String(value))
}
</script>

<template>
  <div class="w-full min-h-12.5 border-2 border-gray-600 flex">
    <div v-if="field.label" class="w-[30%] flex items-center px-2">
      {{ field.label }}
    </div>
    <BaseInput class="flex-1" :model-value="safeValue" :required="field.required" @update:model-value="updateValue" />
  </div>
</template>
