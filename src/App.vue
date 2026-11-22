<script setup lang="ts">
import { onMounted, ref } from "vue";
import { commands, type PlanField_Serialize } from "./generated/bindings";
import FieldForm from "@/modules/fields/components/FieldForm.vue";

const debug_result = ref();
const result = ref();
const plan_id = ref("");
const settings_fields = ref<{ [key: string]: PlanField_Serialize }>();
const settings_values = ref<Record<string, unknown>>({});

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

    <pre>{{ JSON.stringify(settings_values, null, 2) }}</pre>
  </main>
</template>

<style scoped></style>