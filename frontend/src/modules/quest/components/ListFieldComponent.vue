<script setup lang="ts">
import type { ListField } from '@/api/generated'
import FieldSelector from './FieldSelector.vue'

const props = withDefaults(defineProps<{
  field: ListField
  modelValue?: any[]
}>(), {
  modelValue: () => []
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: any[]): void
}>()

const addItem = () => {
  emit('update:modelValue', [...props.modelValue, null])
}

const removeItem = (index: number) => {
  const updated = [...props.modelValue]
  updated.splice(index, 1)
  emit('update:modelValue', updated)
}

const updateItem = (index: number, value: any) => {
  const updated = [...props.modelValue]
  updated[index] = value
  emit('update:modelValue', updated)
}
</script>

<template>
  <div class="list-field">
    <div v-if="field.label" class="font-medium mb-1">{{ field.label }}</div>
    <div v-for="(item, idx) in modelValue" :key="idx" class="flex items-start gap-2 mb-2">
      <FieldSelector
        :field="field.item_field"
        :model-value="item"
        @update:model-value="updateItem(idx, $event)"
        class="flex-1"
      />
      <button
        type="button"
        class="text-red-500 hover:text-red-700 px-2 py-1 border rounded"
        @click="removeItem(idx)"
      >
        ✕
      </button>
    </div>
    <button
      type="button"
      class="mt-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
      @click="addItem"
    >
      + Добавить
    </button>
  </div>
</template>
