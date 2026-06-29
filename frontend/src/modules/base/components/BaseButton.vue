<script setup lang="ts">

interface Props {
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  disabled: false,
  loading: false,
})

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()
</script>

<template>
  <button
    class="bg-gray-300"
    :type="type"
    :disabled="props.disabled || props.loading"
    @click="emit('click', $event)"
  >
    <slot />
  </button>
</template>
