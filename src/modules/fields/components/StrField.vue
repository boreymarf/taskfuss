<script setup lang="ts">
import BaseInput from '@/modules/base/components/BaseInput.vue';
import type { StrFieldSchema } from '@generated/bindings';

const props = defineProps<StrFieldSchema>();
const model = defineModel<string>();
</script>

<template>
  <BaseInput
    v-bind="$attrs"
    type="text"
    v-model="model"
    :required="props.required"
    :maxlength="props.max_length ?? undefined"
    :minlength="props.min_length ?? undefined"
  />
</template>

<style scoped lang="scss"></style>