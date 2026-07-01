<script setup lang="ts">
import { useQuestPlanSetupForm } from './useQuestPlanSetupForm';



const { data, error } = useQuestPlanSetupForm("simple_checklist")
</script>

<template>
  <Story>
    <Variant>
      <h1>data</h1>
      <pre>
        {{ JSON.stringify(data, null, 2) }}
      </pre>
      <h1>error</h1>
      <pre>
        {{ JSON.stringify(error, null, 2) }}
      </pre>
    </Variant>
  </Story>
</template>
