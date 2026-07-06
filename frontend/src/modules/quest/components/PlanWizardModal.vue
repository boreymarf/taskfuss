<script setup lang="ts">
import BaseButton from '@/modules/base/components/BaseButton.vue'
import BaseModal from '@/modules/base/components/BaseModal.vue'
import { ref, computed } from 'vue'
import { useQuestPlans } from '../queries/useQuestPlans'
import { useQuestPlanSetupForm } from '../queries/useQuestPlanSetupForm'
import Form from './Form.vue'

const modelValue = defineModel<boolean>({ default: false })

const {
  data: plans,
  isLoading: plansLoading,
  isError: plansError,
  error: plansErrorObj,
  refetch: refetchPlans,
} = useQuestPlans()
const selectedPlanId = ref<string | null>(null)

const {
  data: setupForm,
  isLoading: formLoading,
  isError: formError,
  error: formErrorObj,
} = useQuestPlanSetupForm(selectedPlanId)
const selectPlan = (id: string) => {
  selectedPlanId.value = id
}

const hasSelectedPlan = computed(() => selectedPlanId.value !== null)
</script>

<template>
  <BaseModal v-model="modelValue" :showConfirm="false">
    <div v-if="!hasSelectedPlan">
      <div v-if="plansLoading">Loading plans...</div>
      <div v-else-if="plansError">Error: {{ plansErrorObj?.message }}</div>
      <div v-else-if="plans && plans.length === 0">No plans found!</div>
      <div v-else-if="plans && plans.length > 0">
        <BaseButton @click="() => refetchPlans()"> Refresh </BaseButton>
        <p>Select a plan</p>
        <BaseButton v-for="plan in plans" :key="plan.id" @click="selectPlan(plan.id)">
          {{ plan.name }}
        </BaseButton>
      </div>
    </div>

    <div v-else>
      <div v-if="formLoading">Loading setup form...</div>
      <div v-else-if="formError">Error: {{ formErrorObj?.message }}</div>
      <div v-else-if="setupForm">
        <Form :fields="setupForm" />
        <BaseButton> Cancel </BaseButton>
        <BaseButton> Submit </BaseButton>
      </div>
    </div>
  </BaseModal>
</template>
