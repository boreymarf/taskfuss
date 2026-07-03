<script setup lang="ts">
import { ref, watch, onMounted, type Ref } from 'vue'
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

// -----------------------------------------------------------------
// 1. Инициализация
// -----------------------------------------------------------------

/** Возвращает массив по умолчанию (из метаданных поля или пустой) */
function getDefaultArray(): any[] {
  return Array.isArray(props.field.default) ? props.field.default : []
}

/** Превращает сырой массив в массив ref-ов для независимой мутации элементов */
function arrayToRefs(arr: any[]): Ref<any>[] {
  return arr.map(item => ref(item))
}

/**
 * Хранилище ref-ов для каждого элемента списка.
 * Сам массив itemRefs меняется только при добавлении/удалении элементов,
 * а при редактировании мутируется значение внутри конкретного ref –
 * это позволяет не терять фокус и не пересоздавать компоненты.
 */
const itemRefs = ref<Ref<any>[]>(
  Array.isArray(props.modelValue) ? arrayToRefs(props.modelValue) : arrayToRefs(getDefaultArray())
)

// При первом рендере, если переданный modelValue не был массивом, эмитим дефолт
onMounted(() => {
  if (!Array.isArray(props.modelValue)) {
    emit('update:modelValue', getDefaultArray())
  }
})

// -----------------------------------------------------------------
// 2. Синхронизация с внешним modelValue
// -----------------------------------------------------------------

/**
 * Если родитель передал новый массив (например, при сбросе формы),
 * мы сравниваем его с текущими значениями ref-ов и при необходимости обновляем ref-ы.
 */
watch(
  () => props.modelValue,
  (newVal) => {
    if (!Array.isArray(newVal)) return
    // Пересоздаём ref-ы только если данные действительно изменились
    const needsUpdate =
      newVal.length !== itemRefs.value.length ||
      newVal.some((item, i) => itemRefs.value[i]?.value !== item)
    if (needsUpdate) {
      itemRefs.value = arrayToRefs(newVal)
    }
  }
)

// -----------------------------------------------------------------
// 3. Отправка изменений наверх (с группировкой)
// -----------------------------------------------------------------

/**
 * Следим за всем массивом ref-ов (глубоко) и эмитим сырой массив значений.
 * flush: 'post' + deep гарантирует, что несколько изменений подряд
 * (например, быстрый ввод) сгенерируют только одно событие update:modelValue.
 */
watch(
  itemRefs,
  () => {
    emit('update:modelValue', itemRefs.value.map(r => r.value))
  },
  { deep: true, flush: 'post' }
)

// -----------------------------------------------------------------
// 4. Методы мутации
// -----------------------------------------------------------------

function addItem() {
  // Добавляем новый элемент со значением null (или по дефолту, если нужно)
  itemRefs.value = [...itemRefs.value, ref(null)]
}

function removeItem(index: number) {
  itemRefs.value = itemRefs.value.filter((_, i) => i !== index)
}

/**
 * Обработчик изменения конкретного элемента.
 * Просто обновляем значение нужного ref – watch сам отправит обновление.
 */
function onItemUpdate(index: number, value: any) {
  itemRefs.value[index].value = value
}
</script>

<template>
  <div class="list-field">
    <div v-if="field.label" class="font-medium mb-1">{{ field.label }}</div>

    <!-- Элементы списка -->
    <div v-for="(itemRef, idx) in itemRefs" :key="idx" class="flex items-start gap-2 mb-2">
      <FieldSelector
        :field="field.item_field"
        :model-value="itemRef.value"
        @update:model-value="onItemUpdate(idx, $event)"
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

    <!-- Кнопка добавления -->
    <button
      type="button"
      class="mt-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
      @click="addItem"
    >
      + Добавить
    </button>
  </div>
</template>
