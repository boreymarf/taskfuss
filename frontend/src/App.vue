<script setup lang="ts">
import Navbar from "@/modules/layout/components/Navbar.vue";
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <Navbar />
    <main class="flex-1 flex flex-col">
      <RouterView />
    </main>
  </div>
</template>

<style scoped></style>
