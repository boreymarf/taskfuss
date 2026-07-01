<script setup lang="ts">
import BaseInput from '@/modules/base/components/BaseInput.vue';
import type { StrField } from '../types/fields';

const props = defineProps<StrField>()
</script>

<template>
  <div class="w-full min-h-12.5 border-2 border-gray-600 flex">
    <BaseInput />
  </div>
</template>
