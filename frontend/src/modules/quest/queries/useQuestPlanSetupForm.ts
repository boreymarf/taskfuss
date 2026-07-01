import { getPlanSetupFormApiQuestPlanPlanIdSetupFormGet } from '@/api/generated'
import { useQuery } from '@tanstack/vue-query'
import type { QuestPlanSetupForm } from '../types/setupForm'

export function useQuestPlanSetupForm(planId: string) {
  return useQuery({
    queryKey: ['questPlanSetupForm', planId],
    queryFn: async () => {
      const result = await getPlanSetupFormApiQuestPlanPlanIdSetupFormGet({
        path: {
          plan_id: planId
        }
      })
      if (result.error) throw result.error
      return result.data as QuestPlanSetupForm
    },
    enabled: !!planId, // запрос не выполнится, пока нет id
    staleTime: 10 * 60 * 1000,
    gcTime: 20 * 60 * 1000,
    retry: 2,
    refetchOnWindowFocus: false,
  })
}
