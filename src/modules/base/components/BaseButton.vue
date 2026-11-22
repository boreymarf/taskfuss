<script setup lang="ts">
interface Props {
  type?: 'button' | 'submit' | 'reset';
}

withDefaults(defineProps<Props>(), {
  type: 'button',
});

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void;
}>();

const handleClick = (event: MouseEvent) => {
  emit('click', event);
};
</script>

<template>
  <button :type="type" v-bind="$attrs" @click="handleClick">
    <slot />
  </button>
</template>

<style scoped lang="scss"></style>