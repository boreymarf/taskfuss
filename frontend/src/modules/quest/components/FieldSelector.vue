<script setup lang="ts">
import type { Field } from '../types/fields'
import StrFieldComponent from './StrFieldComponent.vue'
import ListFieldComponent from './ListFieldComponent.vue'

const model = defineModel<any>({ required: true })
const props = defineProps<{ field: Field }>()
</script>

<template>
  <StrFieldComponent
    v-if="field.discriminator === 'str'"
    v-bind="$attrs"
    v-model="model"
    :field="field"
  />
  <ListFieldComponent
    v-else-if="field.discriminator === 'list'"
    v-bind="$attrs"
    v-model="model"
    :field="field"
  />
  <p v-else v-bind="$attrs">Unknown field type!</p>
</template>
