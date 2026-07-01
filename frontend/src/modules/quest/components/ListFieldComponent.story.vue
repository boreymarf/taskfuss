<script setup lang="ts">
import type { ListField } from '@/api/generated';
import { ref } from 'vue';
import ListFieldComponent from './ListFieldComponent.vue';


let example = ref<ListField>({
  "discriminator": "list",
  "value_type": "list[object]",
  "label": null,
  "description": null,
  "item_field": {
    "discriminator": "str",
    "value_type": "str",
    "label": null,
    "description": null,
    "required": false,
    "default": null,
    "max_length": null,
    "min_length": null
  },
  "min_items": null,
  "max_items": null
})
let value = ref<any>()

</script>

<template>
  <Story>
    <Variant>
      <ListFieldComponent :field="example" :modelValue="value" @update:modelValue="value = $event"/>
    </Variant>
  </Story>
</template>
