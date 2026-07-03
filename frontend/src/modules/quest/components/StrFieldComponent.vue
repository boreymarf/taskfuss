<script setup lang="ts">
import type { Ref } from 'vue'
import type { Field } from '../types/fields'
import BaseInput from '@/modules/base/components/BaseInput.vue';

const props = defineProps<{
  fieldRef: Ref<any>
  field: Field
}>()
</script>

<template>
  <div class="w-full min-h-12.5 border-2 border-gray-600 flex">
    <div v-if="field.label" class="w-[30%] flex items-center px-2">
      {{ field.label }}
    </div>
    <BaseInput v-model="fieldRef.value" class="flex-1" />
  </div>
</template>
