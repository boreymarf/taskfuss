<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})
</script>

<template>
  <div class="container">{{ task.title }}</div>
</template>

<style lang="scss" scoped>
.container {
  height: 30px;
  width: 100%;
  background-color: green;
  color: white;
  display: flex;
  padding: 0 10px;
}
</style>
