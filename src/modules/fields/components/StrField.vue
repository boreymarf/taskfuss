<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import BaseInput from '@/modules/base/components/BaseInput.vue';
import { formatFieldError } from '@/i18n';
import type { FieldValueError, StrFieldSchema } from '@generated/bindings';

const props = defineProps<{
  field: StrFieldSchema;
  error?: FieldValueError;
}>();
const model = defineModel<string>();

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
      type="text"
      v-model="model"
      :required="props.field.required"
      :maxlength="props.field.max_length ?? undefined"
      :minlength="props.field.min_length ?? undefined"
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
  color: red;
  font-size: 0.875rem;
  text-align: right;
  margin-left: auto;
}
</style>