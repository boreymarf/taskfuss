<script setup lang="ts">
import { ref } from 'vue'
import PlanWizardModal from './PlanWizardModal.vue'
import BaseButton from '@/modules/base/components/BaseButton.vue'

const isOpen = ref<boolean>(false)

const open = () => {
  isOpen.value = true
}
</script>

<template>
  <Story>
    <Variant>
      <BaseButton @click="open"> Create new quest </BaseButton>
      <PlanWizardModal v-model="isOpen" />
    </Variant>
  </Story>
</template>
