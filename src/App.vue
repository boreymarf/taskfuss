<script setup lang="ts">
import { onMounted, ref } from "vue";
import { commands } from "./generated/bindings";

const debug_result = ref();
const result = ref();

async function debugSomething() {
  debug_result.value = await commands.debugSomething()
}

async function getAllPlans() {
  // Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
  result.value = await commands.getAllPlans()
}
</script>

<template>
  <main class="container">
    <h1>This is a test page</h1>
    <button @click="debugSomething">Debug something</button>
    <pre>{{ JSON.stringify(debug_result, null, 2) }}</pre>
    <button @click="getAllPlans">Get all plans</button>
    <pre>{{ JSON.stringify(result, null, 2) }}</pre>
  </main>
</template>

<style scoped></style>
