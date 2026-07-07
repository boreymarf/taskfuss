<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import BaseButton from '@/modules/base/components/BaseButton.vue'
import BaseModal from '@/modules/base/components/BaseModal.vue'
import Form from './Form.vue'
import { useQuestPlans } from '../queries/useQuestPlans'
import { useQuestPlanSetupForm } from '../queries/useQuestPlanSetupForm'
import { useCreateQuest } from '../mutations/useCreateQuest'

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

const setupFormData = ref(undefined)

const { mutate: createQuest, error: createQuestError } = useCreateQuest()

const hasSelectedPlan = computed(() => selectedPlanId.value !== null)

const selectPlan = (id: string) => {
  selectedPlanId.value = id
}

const close = () => {
  modelValue.value = false
}

const submitSetupForm = () => {
  if (selectedPlanId.value != null && setupFormData.value != null) {
    createQuest({ plan_id: selectedPlanId.value, setup_form: setupFormData.value })
    close()
  }
}

watch(modelValue, (newValue) => {
  if (!newValue) {
    selectedPlanId.value = null
    setupFormData.value = undefined
  }
})
</script>

<template>
  <BaseModal v-model="modelValue" :showConfirm="false">
    <div v-if="!hasSelectedPlan">
      <div v-if="plansLoading">Loading plans...</div>
      <div v-else-if="plansError">Error: {{ plansErrorObj?.message }}</div>
      <div v-else-if="plans && plans.length === 0">No plans found!</div>
      <div v-else-if="plans && plans.length > 0">
        <BaseButton @click="refetchPlans()">Refresh</BaseButton>
        <p>Select a plan</p>
        <BaseButton
          v-for="plan in plans"
          :key="plan.id"
          @click="selectPlan(plan.id)"
        >
          {{ plan.name }}
        </BaseButton>
      </div>
    </div>

    <div v-else>
      <div v-if="formLoading">Loading setup form...</div>
      <div v-else-if="formError">Error: {{ formErrorObj?.message }}</div>
      <div v-else-if="setupForm">
        <Form :fields="setupForm" v-model="setupFormData" />
        <div v-if="createQuestError" class="error">
          Failed to create quest: {{ createQuestError.message }}
        </div>
        <BaseButton @click="close">Cancel</BaseButton>
        <BaseButton @click="submitSetupForm">Submit</BaseButton>
      </div>
    </div>
  </BaseModal>
</template>
