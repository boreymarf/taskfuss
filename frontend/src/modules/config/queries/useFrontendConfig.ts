import { getFrontendConfigurationApiConfigGet, type FrontendConfig } from '@/api/generated'
import { useQuery } from '@tanstack/vue-query'

export function useFrontendConfig() {
  return useQuery({
    queryKey: ['frontendConfig'],
    queryFn: async () => {
      const result = await getFrontendConfigurationApiConfigGet()
      if (result.error) throw result.error
      return result.data as FrontendConfig
    },
    staleTime: 20 * 60 * 1000, // 20 min
    gcTime: 30 * 60 * 1000, // 30 min
    retry: 3,
    refetchOnWindowFocus: false,
  })
}
