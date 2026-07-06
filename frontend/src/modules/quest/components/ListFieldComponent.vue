<script setup lang="ts">
import { ref, watch, type Ref } from "vue";
import type { ListField } from "@/api/generated";
import FieldSelector from "./FieldSelector.vue";

const props = defineProps<{
  field: ListField;
}>();

const model = defineModel<any[]>({ default: () => [] });

function arrayToRefs(arr: any[]): Ref<any>[] {
  return arr.map((item) => ref(item));
}

const itemRefs = ref<Ref<any>[]>(
  Array.isArray(model.value) ? arrayToRefs(model.value) : arrayToRefs([]),
);

// good enough I guess
watch(
  () => model.value,
  (newVal) => {
    if (!Array.isArray(newVal)) return;
    const needsUpdate =
      newVal.length !== itemRefs.value.length ||
      newVal.some((item, i) => itemRefs.value[i]?.value !== item);
    if (needsUpdate) {
      itemRefs.value = arrayToRefs(newVal);
    }
  },
);

watch(
  itemRefs,
  () => {
    model.value = itemRefs.value.map((r) => r.value);
  },
  { deep: true, flush: "post" },
);

function addItem() {
  itemRefs.value = [...itemRefs.value, ref(null)];
}

function removeItem(index: number) {
  itemRefs.value = itemRefs.value.filter((_, i) => i !== index);
}

function onItemUpdate(index: number, value: any) {
  const target = itemRefs.value[index];
  if (target) {
    target.value = value;
  }
}
</script>

<template>
  <div class="list-field">
    <div v-if="field.label" class="font-medium mb-1">
      {{ field.label }}
    </div>

    <div
      v-for="(itemRef, idx) in itemRefs"
      :key="idx"
      class="flex items-start gap-2 mb-2"
    >
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

    <button
      type="button"
      class="mt-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
      @click="addItem"
    >
      + Добавить
    </button>
  </div>
</template>
