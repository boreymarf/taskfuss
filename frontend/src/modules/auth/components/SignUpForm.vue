<script setup lang="ts">
import BaseButton from "@/modules/base/components/BaseButton.vue";
import BaseInput from "@/modules/base/components/BaseInput.vue";
import { ref } from "vue";

export type UserCreate = {
  login: string;
  password: string;
};

const formData = ref<UserCreate>({
  login: "",
  password: "",
});

const emit = defineEmits<{
  (e: "submit", data: UserCreate): void;
}>();

const handleSubmit = () => {
  emit("submit", formData.value);
};
</script>

<template>
  <form @submit.prevent="handleSubmit" class="signup-form">
    <BaseInput v-model="formData.login" placeholder="Login" autofocus />
    <BaseInput
      v-model="formData.password"
      type="password"
      placeholder="Password"
    />

    <BaseButton type="submit"> Sign Up </BaseButton>
  </form>
</template>

<style scoped>
.signup-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
