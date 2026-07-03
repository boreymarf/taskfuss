<script setup lang="ts">
import { reactive, watch } from 'vue'
import FieldSelector from './FieldSelector.vue'
import type { Field } from '../types/fields'

const props = defineProps<{
  fields: Record<string, Field>
}>()

// defineModel даёт writable computed для v-model
const model = defineModel<Record<string, any>>({ default: () => ({}) })

/**
 * Инициализируем локальное состояние формы.
 * Все поля получают null — конкретный тип поля сам выставит дефолт ('' для строки, [] для списка).
 */
const formData = reactive<Record<string, any>>(
  Object.keys(props.fields).reduce((acc, key) => {
    acc[key] = model.value?.[key] ?? null
    return acc
  }, {} as Record<string, any>)
)

// Предыдущее значение для сравнения, чтобы не эмитить одинаковые данные повторно
let previousModelValue = JSON.stringify(model.value)

/**
 * Следим за formData и выносим изменения наверх (в родительский v-model).
 * flush: 'post' группирует мутации внутри одного тика, что предотвращает лавину апдейтов.
 * JSON-сравнение гарантирует, что emit произойдёт только при реальном изменении данных.
 */
watch(
  formData,
  () => {
    const newValue = { ...formData }
    const newString = JSON.stringify(newValue)
    if (newString !== previousModelValue) {
      previousModelValue = newString
      model.value = newValue
    }
  },
  { deep: true, flush: 'post' }
)
</script>

<template>
  <FieldSelector
    v-for="(field, key) in fields"
    :key="key"
    v-model="formData[key]"
    :field="field"
  />
</template>
