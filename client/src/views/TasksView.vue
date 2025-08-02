<script setup lang="ts">
import { onMounted } from 'vue';
import { useTasksStore } from '../stores/tasks';
import TaskEntryButton from '../components/TaskEntryButton.vue';

const tasksStore = useTasksStore()

onMounted(() => {
  tasksStore.fetchAllTasks();
});
</script>

<template>
  <div class="main">
    <h1>Tasks</h1>

    <div v-if="tasksStore.isLoading">Loading tasks...</div>
    <div v-if="tasksStore.error" class="error">{{ tasksStore.error }}</div>

    <div class="grid-header container">
      <div>Task Title</div>
      <div>1 август</div>
      <div>31 июня</div>
      <div>30 июня</div>
      <div>29 июня</div>
    </div>
    <ul v-if="tasksStore.tasks.length > 0" class="grid-container">
      <li v-for="task in tasksStore.tasks" :key="task.id" class="grid-item container">
        <div>{{ task.title }}</div>
        <TaskEntryButton />
        <TaskEntryButton />
        <TaskEntryButton />
        <TaskEntryButton />
      </li>
    </ul>

    <!-- Empty state -->
    <div v-else-if="!tasksStore.isLoading">
      No tasks found.
    </div>
  </div>
</template>

<style lang="scss" scoped>
.grid-header {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
  font-weight: bold;
  padding: 10px;
}

.grid-container {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  list-style: none;
  padding: 0;
}

.grid-item {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
  padding: 10px;
  border: 1px solid #ddd;
}

/* Responsive adjustments */
@media (max-width: 768px) {

  .grid-header,
  .grid-item {
    grid-template-columns: 1fr;
  }

  .grid-header div,
  .grid-item div {
    padding: 5px 0;
    border-bottom: 1px solid #eee;
  }

  .grid-header div:last-child,
  .grid-item div:last-child {
    border-bottom: none;
  }
}
</style>
