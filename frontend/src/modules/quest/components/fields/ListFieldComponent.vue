<script setup lang="ts">
import { ref, watch, type Ref } from 'vue'
import FieldSelector from './FieldSelector.vue'
import type { ListField } from '../../types/fields'
import BaseButton from '@/modules/base/components/BaseButton.vue'

const props = defineProps<{
  field: ListField
}>()

const model = defineModel<any[]>({ default: () => [] })

function arrayToRefs(arr: any[]): Ref<any>[] {
  return arr.map((item) => ref(item))
}

const itemRefs = ref<Ref<any>[]>(
  Array.isArray(model.value) ? arrayToRefs(model.value) : arrayToRefs([]),
)

// good enough I guess
watch(
  () => model.value,
  (newVal) => {
    if (!Array.isArray(newVal)) return
    const needsUpdate =
      newVal.length !== itemRefs.value.length ||
      newVal.some((item, i) => itemRefs.value[i]?.value !== item)
    if (needsUpdate) {
      itemRefs.value = arrayToRefs(newVal)
    }
  },
)

watch(
  itemRefs,
  () => {
    model.value = itemRefs.value.map((r) => r.value)
  },
  { deep: true, flush: 'post' },
)

function addItem() {
  itemRefs.value = [...itemRefs.value, ref(null)]
}

function removeItem(index: number) {
  itemRefs.value = itemRefs.value.filter((_, i) => i !== index)
}

function onItemUpdate(index: number, value: any) {
  const target = itemRefs.value[index]
  if (target) {
    target.value = value
  }
}
</script>

<template>
  <div class="list-field">
    <div v-if="field.label" class="font-medium mb-1">
      {{ field.label }}
    </div>

    <div v-for="(itemRef, idx) in itemRefs" :key="idx" class="flex items-start gap-2 mb-2">
      <FieldSelector
        :field="field.item_field"
        :model-value="itemRef.value"
        @update:model-value="onItemUpdate(idx, $event)"
        class="flex-1"
      />
      <BaseButton
        class="bg-red-500 hover:bg-red-700 px-2 py-1"
        @click="removeItem(idx)"
      >
        ✕
      </BaseButton>
    </div>

    <BaseButton @click="addItem"> Add new item </BaseButton>
  </div>
</template>
