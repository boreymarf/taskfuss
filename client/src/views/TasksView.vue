<script setup lang="ts">
import { onMounted } from 'vue';
import { useTasksStore } from '../stores/tasks';
import TaskRow from "../components/TaskRow.vue"

const tasksStore = useTasksStore()

onMounted(() => {
  tasksStore.fetchAllTasks();
});
</script>

<template>
  <h1>Tasks</h1>

  <div v-if="tasksStore.isLoading">Loading tasks...</div>
  <div v-if="tasksStore.error" class="error">{{ tasksStore.error }}</div>

  <ul v-if="tasksStore.tasks.length > 0">
    <TaskRow v-for="task in tasksStore.tasks" :key="task.id" :task="task" />
  </ul>

  <!-- Empty state -->
  <div v-else-if="!tasksStore.isLoading">
    No tasks found.
  </div>

</template>

<style lang="scss" scoped></style>
