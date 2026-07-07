<script setup lang="ts">
import BaseButton from '@/modules/base/components/BaseButton.vue'
import PlanWizardModal from '@/modules/quest/components/PlanWizardModal.vue'
import { ref } from 'vue'

const isOpen = ref<boolean>(false)

const open = () => {
  isOpen.value = true
}
</script>

<template>
  <div>
    <h1>This is a front page!</h1>
    <p>It's currently empty</p>
    <BaseButton @click="open"> Create new quest </BaseButton>
    <PlanWizardModal v-model="isOpen" />
  </div>
</template>
