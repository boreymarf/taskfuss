<script setup lang="ts">
import { reactive, toRefs, watch, type Ref } from 'vue'
import FieldSelector from './FieldSelector.vue'
import type { Field } from '../types/fields'

const props = defineProps<{
  fields: Record<string, Field>
  modelValue?: Record<string, any>   // <-- это убирает предупреждение
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Record<string, any>): void
}>()

// Вся остальная логика без изменений
const formData = reactive<Record<string, any>>(
  Object.keys(props.fields).reduce((acc, key) => {
    acc[key] = null
    return acc
  }, {} as Record<string, any>)
)

const fieldRefs = toRefs(formData) as Record<string, Ref<any>>

watch(
  formData,
  (newValue) => {
    emit('update:modelValue', newValue)
  },
  { deep: true, flush: 'post' }
)
</script>

<template>
  <FieldSelector v-for="(field, key) in fields" :key="key" :field-ref="fieldRefs[key]!" :field="field" />
</template>
