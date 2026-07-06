<!-- Form.story.vue -->
<script setup lang="ts">
import { ref } from "vue";
import Form from "./Form.vue";
import type { Field } from "../types/fields";

// Описание полей
const exampleSettings = ref<Record<string, Field>>({
  name: {
    discriminator: "str",
    value_type: "str",
    label: "Quest name (optional)",
    description: null,
    required: false,
    default: null,
    max_length: null,
    min_length: null,
  },
  task_list: {
    discriminator: "list",
    value_type: "list[object]",
    label: null,
    description: null,
    item_field: {
      discriminator: "str",
      value_type: "str",
      label: null,
      description: null,
      required: false,
      default: null,
      max_length: null,
      min_length: null,
    },
    min_items: null,
    max_items: null,
  },
});

const formResult = ref<Record<string, any>>({});
</script>

<template>
  <Story>
    <Variant>
      <!-- v-model автоматически подпишет formResult на update:modelValue -->
      <Form :fields="exampleSettings" v-model="formResult" />

      <!-- Для отладки: покажем результат -->
      <pre>{{ formResult }}</pre>
    </Variant>
  </Story>
</template>
