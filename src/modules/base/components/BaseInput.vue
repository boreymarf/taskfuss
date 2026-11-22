<script setup lang="ts">
interface Props {
  type?: string;
}

withDefaults(defineProps<Props>(), {
  type: 'text',
});

const model = defineModel<string | number>();
</script>

<template>
  <input :type="type" v-bind="$attrs" v-model="model" />
</template>

<style scoped lang="scss"></style>