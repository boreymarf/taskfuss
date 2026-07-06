<script setup lang="ts">
import BaseCheckbox from "./BaseCheckbox.vue";
import { logEvent } from "histoire/client";
import { ref } from "vue";

const checked = ref(false);
</script>

<template>
  <Story :layout="{ type: 'grid', width: 400 }">
    <Variant title="normal">
      <BaseCheckbox v-model="checked" @change="logEvent('Change', $event)" />
    </Variant>
    <Variant title="disabled">
      <BaseCheckbox
        v-model="checked"
        :disabled="true"
        @change="logEvent('Change', $event)"
      />
    </Variant>
    <Variant title="indeterminate">
      <BaseCheckbox
        v-model="checked"
        :indeterminate="true"
        @change="logEvent('Change', $event)"
      />
    </Variant>
  </Story>
</template>
