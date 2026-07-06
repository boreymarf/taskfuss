<script setup lang="ts">
import { ref } from "vue";
import BaseButton from "./BaseButton.vue";
import BaseModal from "./BaseModal.vue";

const isOpen = ref(false);

const open = () => {
  isOpen.value = true;
};
</script>

<template>
  <Story :layout="{ type: 'grid', width: '100%' }">
    <div>
      <BaseButton @click="open"> Open modal window </BaseButton>

      <BaseModal v-model="isOpen">
        <h2>This is a modal example</h2>
        <p>You can put anything in here.</p>
      </BaseModal>
    </div>
  </Story>
</template>
