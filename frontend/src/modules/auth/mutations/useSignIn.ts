import {
  signInApiAuthSignInPost,
  type Token,
  type UserLogin,
} from "@/api/generated";
import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { useAuthStore } from "../stores/useAuthStore";
import router from "@/router";

export function useSignIn() {
  const authStore = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userData: UserLogin) =>
      signInApiAuthSignInPost({ body: userData }).then((result) => {
        if (result.error) throw result.error;
        return result.data as Token;
      }),
    onSuccess: (token) => {
      authStore.setAuth(token.access_token, token.token_type);
      queryClient.invalidateQueries({ queryKey: ["user"] });
      router.push("/");
    },
  });
}
