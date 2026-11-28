<script setup lang="ts">
import { computed, ref, watch } from "vue";
import {
  commands,
  type AppError,
  type Field_Deserialize,
  type FieldValueError,
} from "@/generated/bindings";
import { extractSpectaError } from "@/utils/errors";
import { unwrapAppError } from "@/utils/errors";
import FieldForm from "@modules/fields/components/FieldForm.vue";
import { useDebounce } from "@modules/base/composables/debounce";
import BaseCombobox from "./modules/base/components/BaseCombobox.vue";

const debug_result = ref();
const result = ref();
const plan_id = ref("");
const settings_fields = ref<{ [key: string]: Field_Deserialize }>();
const settings_values = ref<Record<string, unknown>>({});
const validation_error = ref<AppError | null>(null);

// Narrow to kind "form" via unwrapAppError (logs on mismatch), then extract
// the active tag of FormError via Specta helper. Only the "Invalid" variant
// holds a per-field dictionary; other form-level variants (e.g. "Required")
// yield undefined so no field gets an error.
const field_errors = computed<Record<string, FieldValueError> | undefined>(() => {
  const e = validation_error.value;
  if (!e) return undefined;

  const formError = unwrapAppError(e, "form");
  if (!formError) return undefined;

  const entry = extractSpectaError(formError);
  return entry.tag === "Invalid" ? entry.payload : undefined;
});

async function debugSomething() {
  debug_result.value = await commands.debugSomething()
}

async function getAllPlans() {
  result.value = await commands.getAllPlans()
}

async function getPlanSettingsFields() {
  settings_fields.value = await commands.getPlanSettingsFields(plan_id.value)
  settings_values.value = {}
}

async function validateSettingsForm() {
  try {
    validation_error.value = null;
    await commands.validateSettingsForm(plan_id.value, settings_values.value as any);
  } catch (e) {
    validation_error.value = e as AppError;
    console.error("validate_settings_form failed:", e);
  }
}

const { debounced: debouncedValidate } = useDebounce(validateSettingsForm, 300);

watch(settings_values, debouncedValidate, { deep: true });
</script>

<template>
  <main class="container">
    <h1>Component test</h1>
    <BaseCombobox/>

    <h1>This is a test page</h1>
    <button @click="debugSomething">Debug something</button>
    <pre>{{ JSON.stringify(debug_result, null, 2) }}</pre>
    <button @click="getAllPlans">Get all plans</button>
    <pre>{{ JSON.stringify(result, null, 2) }}</pre>

    <input v-model="plan_id" placeholder="plan id" />
    <button @click="getPlanSettingsFields">Get plan settings fields</button>
    <pre>{{ JSON.stringify(settings_fields, null, 2) }}</pre>

    <FieldForm
      v-if="settings_fields"
      v-model="settings_values"
      :fields="settings_fields"
      :errors="field_errors"
    />

    <pre v-if="validation_error" style="color: red">{{ JSON.stringify(validation_error, null, 2) }}</pre>

    <pre>{{ JSON.stringify(settings_values, null, 2) }}</pre>
  </main>
</template>

<style scoped></style>