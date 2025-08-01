import { defineStore } from "pinia";
import axios from "axios";
import { useAuthStore } from "./auth";
import { useRouter } from "vue-router";
import { api } from "../api/client/http";
import type { DtoGetAllTasksResponse, DtoTask } from "../api/generated";

const router = useRouter()
export const useTasksStore = defineStore("tasks", {
  state: () => ({
    tasks: [] as DtoTask[],
    isLoading: false as Boolean,
    error: null as null | Error | unknown
  }),

  getters: {

  },

  actions: {
    async fetchAllTasks() {
      const authStore = useAuthStore()

      if (!authStore.auth_token) {
        this.error = 'No authentication token found';
        router.push('login/')
        return;
      }

      this.isLoading = true;
      try {
        const response = await api.tasks.tasksGet({ authorization: `Bearer ${authStore.auth_token}` })

        const dtoGetAllTasksResponse = response.data as DtoGetAllTasksResponse;

        if (!dtoGetAllTasksResponse) {
          throw new Error('Invalid response data: missing response payload');
        }

        this.tasks = dtoGetAllTasksResponse.tasks || [];
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
