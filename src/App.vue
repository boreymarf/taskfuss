<script setup lang="ts">
import { ref, watch } from "vue";
import { commands, type Field_Deserialize } from "./generated/bindings";
import FieldForm from "@modules/fields/components/FieldForm.vue";
import { useDebounce } from "@modules/base/composables/debounce";

const debug_result = ref();
const result = ref();
const plan_id = ref("");
const settings_fields = ref<{ [key: string]: Field_Deserialize }>();
const settings_values = ref<Record<string, unknown>>({});
const validation_error = ref<unknown>(null);

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
    validation_error.value = e;
    console.error("validate_settings_form failed:", e);
  }
}

const { debounced: debouncedValidate } = useDebounce(validateSettingsForm, 300);

watch(settings_values, debouncedValidate, { deep: true });
</script>

<template>
  <main class="container">
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
    />

    <pre v-if="validation_error" style="color: red">{{ JSON.stringify(validation_error, null, 2) }}</pre>

    <pre>{{ JSON.stringify(settings_values, null, 2) }}</pre>
  </main>
</template>

<style scoped></style>