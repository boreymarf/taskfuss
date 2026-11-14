<script setup lang="ts">
import { onMounted, ref } from "vue";
import { invoke } from "@tauri-apps/api/core";

const greetMsg = ref("");
const dataPath = ref<string>("");
const name = ref("");


onMounted(async () => {
  dataPath.value = await invoke<string>("debug");
});

async function debug_something() {
  // Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
  await invoke("debug_something");
}



</script>

<template>
  <main class="container">
    <h1>This is a test page</h1>
    <button @click="debug_something">Do something</button>
  </main>
</template>

<style scoped>

</style>

