import { defineStore } from "pinia";
import type { Task } from "../types/stores/tasks";
import axios from "axios";

export const useTasksStore = defineStore("tasks", {
  state: () => ({
    tasks: [] as Task[],
    isLoading: false as Boolean,
    error: null as null | Error | unknown
  }),

  getters: {

  },

  actions: {
    async fetchTasks() {
      this.isLoading = true;
      try {
        const response = await axios.get('/api/tasks'); // Запрос к вашему Go API
        this.tasks = response.data; // Предположим, что Go возвращает JSON-массив задач
      } catch (err) {
        this.error = err;
      } finally {
        this.isLoading = false;
      }
    },
  }
})
