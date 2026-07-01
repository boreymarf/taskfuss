<script setup lang="ts">
import { useQuestPlans } from './useQuestPlans';


const { data, error } = useQuestPlans()
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
