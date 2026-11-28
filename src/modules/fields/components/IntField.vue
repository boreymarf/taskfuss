<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import BaseInput from '@/modules/base/components/BaseInput.vue';
import { formatFieldError } from '@/i18n';
import type { FieldValueError, IntFieldSchema } from '@generated/bindings';

const props = defineProps<{
  field: IntFieldSchema;
  error?: FieldValueError;
}>();
const model = defineModel<number>();

const { t } = useI18n();

const errorText = computed(() =>
  props.error ? formatFieldError(props.error, t) : null,
);
</script>

<template>
  <div>
    <div class="field-header">
      <div class="field-label">
        {{ props.field.label }}
      </div>
      <div v-if="errorText" class="field-error">
        {{ errorText }}
      </div>
    </div>
    <BaseInput
      v-bind="$attrs"
      type="number"
      v-model="model"
      :required="props.field.required"
      :min="props.field.min ?? undefined"
      :max="props.field.max ?? undefined"
    />
  </div>
</template>

<style scoped lang="scss">
.field-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.25rem;
  gap: 0.5rem;
}

.field-label {
  font-size: 0.875rem;
}

.field-error {
  color: #c00;
  font-size: 0.875rem;
  text-align: right;
  margin-left: auto;
}
</style>