<script setup lang="ts">
import { onMounted, ref } from "vue";
import { commands } from "./generated/bindings";

const result = ref();

async function debug_something() {
  // Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
  result.value = await commands.getAllPlans()
}
</script>

<template>
  <main class="container">
    <h1>This is a test page</h1>
    <button @click="debug_something">Do something</button>
    <pre>{{ JSON.stringify(result, null, 2) }}</pre>
  </main>
</template>

<style scoped></style>
