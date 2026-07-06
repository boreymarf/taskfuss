import FrontView from "@/views/FrontView.vue";
import SignInView from "@/views/SignInView.vue";
import SignUpView from "@/views/SignUpView.vue";
import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: FrontView,
    },
    {
      path: "/sign-in/",
      name: "sign-in",
      component: SignInView,
    },
    {
      path: "/sign-up/",
      name: "sign-up",
      component: SignUpView,
    },
  ],
});

export default router;
