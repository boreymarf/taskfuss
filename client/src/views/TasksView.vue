<script setup lang="ts">
import { onMounted } from 'vue';
import { useTasksStore } from '../stores/tasks';
// import { useAuthStore } from '../stores/auth';

// const authStore = useAuthStore()
const tasksStore = useTasksStore()

onMounted(() => {
  tasksStore.fetchTasks();
});
</script>

<template>
  <div>

  </div>
</template>

<style lang="scss" scoped></style>
