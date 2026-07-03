<script setup lang="ts">
import { computed, onMounted } from 'vue'
import type { Field } from '../types/fields'
import BaseInput from '@/modules/base/components/BaseInput.vue'

const props = defineProps<{
  modelValue: any
  field: Field
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>()

/**
 * Если при создании значение пустое, выставляем дефолт.
 * Используем field.default, если он задан, иначе пустую строку.
 */
onMounted(() => {
  if (props.modelValue === null || props.modelValue === undefined) {
    const defaultValue = props.field.default ?? ''
    emit('update:modelValue', defaultValue)
  }
})

/**
 * Прокси для v-model: отдаём текущее значение (или '', чтобы не было null в input)
 * и пробрасываем изменения наверх.
 */
const modelProxy = computed({
  get: () => props.modelValue ?? '',
  set: (val) => emit('update:modelValue', val)
})
</script>

<template>
  <div class="w-full min-h-12.5 border-2 border-gray-600 flex">
    <div v-if="field.label" class="w-[30%] flex items-center px-2">
      {{ field.label }}
    </div>
    <BaseInput v-model="modelProxy" class="flex-1" />
  </div>
</template>
