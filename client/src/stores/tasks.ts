import { defineStore } from "pinia";
import type { Task } from "../types/stores/tasks";
import axios from "axios";
import { useAuthStore } from "./auth";
import { useRouter } from "vue-router";

const router = useRouter()
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
        router.push('login/')
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
        if (axios.isAxiosError(err)) {
          if (err.response?.status === 401) {
            authStore.logout()
            router.push('login/')
          } else {
            this.error = err;
          }
        }
      } finally {
        this.isLoading = false;
      }
    },
  }
})
