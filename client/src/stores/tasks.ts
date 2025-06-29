import { defineStore } from "pinia";
import type { Task } from "../types/stores/tasks";
import axios from "axios";
import { useAuthStore } from "./auth";

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
      const authStore = useAuthStore()

      if (!authStore.auth_token) {
        this.error = 'No authentication token found';
        return;
      }

      this.isLoading = true;
      try {
        const response = await axios.get('/api/tasks', {
          headers: {
            Authorization: `Bearer ${authStore.auth_token}`
          }
        })
        this.tasks = response.data;
      } catch (err) {
        this.error = err;
      } finally {
        this.isLoading = false;
      }
    },
  }
})
