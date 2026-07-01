import { getAllQuestPlansApiQuestPlanGet, type PlanRegistryPublic } from '@/api/generated'
import { useQuery } from '@tanstack/vue-query'

export function useQuestPlans() {
  return useQuery({
    queryKey: ['questPlans'],
    queryFn: async () => {
      const result = await getAllQuestPlansApiQuestPlanGet()
      if (result.error) throw result.error
      return result.data as PlanRegistryPublic[]
    },
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
    retry: 2,
    refetchOnWindowFocus: false,
  })
}
