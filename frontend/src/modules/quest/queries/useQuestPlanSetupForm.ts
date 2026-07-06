import { getPlanSetupFormApiQuestPlanPlanIdSetupFormGet } from '@/api/generated'
import { useQuery } from '@tanstack/vue-query'
import type { QuestPlanSetupForm } from '../types/setupForm'
import { toValue, type MaybeRef } from 'vue'

export function useQuestPlanSetupForm(planId: MaybeRef<string | null>) {
  return useQuery({
    queryKey: ['questPlanSetupForm', toValue(planId)],
    queryFn: async () => {
      const id = toValue(planId)
      if (!id) throw new Error('planId is required')
      const result = await getPlanSetupFormApiQuestPlanPlanIdSetupFormGet({
        path: { plan_id: id },
      })
      if (result.error) throw result.error
      return result.data as QuestPlanSetupForm
    },
    enabled: () => !!toValue(planId), // enabled может быть функцией
    staleTime: 10 * 60 * 1000,
    gcTime: 20 * 60 * 1000,
    retry: 2,
    refetchOnWindowFocus: false,
  })
}
