import { signUpApiAuthSignUpPost, type Token, type UserCreate } from '@/api/generated'
import { useMutation } from '@tanstack/vue-query'
import { useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/useAuthStore'

export function useSignUp() {
  const authStore = useAuthStore()
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation({
    mutationFn: (userData: UserCreate) =>
      signUpApiAuthSignUpPost({
        body: userData,
      }).then((result) => {
        if (result.error) throw result.error
        return result.data as Token
      }),
    onSuccess: (token) => {
      authStore.setAuth(token.access_token, token.token_type)
      queryClient.invalidateQueries({ queryKey: ['user'] })
      router.push('/')
    },
  })
}
