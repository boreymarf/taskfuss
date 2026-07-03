<script setup lang="ts">
import { computed } from 'vue'
import type { Field } from '../types/fields'
import StrFieldComponent from './StrFieldComponent.vue'
import ListFieldComponent from './ListFieldComponent.vue'

const props = defineProps<{
  modelValue: any
  field: Field
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>()

/**
 * Прокси для безопасного использования v-model.
 * Не даёт мутировать prop напрямую, а эмитит событие обновления.
 */
const modelProxy = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})
</script>

<template>
  <!-- Выбираем компонент поля по типу, пробрасываем все $attrs (например, class) -->
  <StrFieldComponent
    v-if="field.discriminator === 'str'"
    v-bind="$attrs"
    v-model="modelProxy"
    :field="field"
  />
  <ListFieldComponent
    v-else-if="field.discriminator === 'list'"
    v-bind="$attrs"
    v-model="modelProxy"
    :field="field"
  />
  <p v-else v-bind="$attrs">Unknown field type!</p>
</template>
