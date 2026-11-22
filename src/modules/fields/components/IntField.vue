<script setup lang="ts">
import BaseInput from '@/modules/base/components/BaseInput.vue';
import type { IntFieldSchema } from '@generated/bindings';

const props = defineProps<IntFieldSchema>();
const model = defineModel<number>();
</script>

<template>
  <BaseInput
    v-bind="$attrs"
    type="number"
    v-model="model"
    :required="props.required"
    :min="props.min ?? undefined"
    :max="props.max ?? undefined"
  />
</template>

<style scoped lang="scss"></style>