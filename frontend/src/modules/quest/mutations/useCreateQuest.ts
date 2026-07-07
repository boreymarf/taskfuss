// composables/useCreateQuest.ts
import { useMutation } from '@tanstack/vue-query'
import { createQuestApiQuestPost, type QuestCreate } from '@/api/generated'
import { useAuthStore } from '@/modules/auth'

export function useCreateQuest() {
  const authStore = useAuthStore()
  // const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (questData: QuestCreate) => {
      const authHeader = authStore.authHeader // уже строка или null

      if (!authHeader) {
        throw new Error('Not authorized')
      }

      const result = await createQuestApiQuestPost({
        body: questData,
        headers: {
          Authorization: authHeader,
        },
      })
      if (result.error) throw result.error
      return result.data
    },

    onSuccess: (_data) => {
      // queryClient.invalidateQueries({ queryKey: ['quests'] })
      // Можно ещё редиректить, как в signIn:
      // router.push("/quests");
    },
  })
}
