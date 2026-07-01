<script setup lang="ts">
interface Props {
  type?: 'text' | 'password' | 'email' | 'number' | 'tel' | 'url'
  modelValue?: string | number
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  name?: string
  id?: string
  autofocus?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  modelValue: '',
  placeholder: '',
  disabled: false,
  readonly: false,
  name: undefined,
  id: undefined,
  autofocus: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
  (e: 'input', event: Event): void
  (e: 'focus', event: FocusEvent): void
  (e: 'blur', event: FocusEvent): void
}>()

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  const value = props.type === 'number' ? Number(target.value) : target.value
  emit('update:modelValue', value)
  emit('input', event)
}
</script>

<template>
  <input
    class="bg-gray-300 w-full"
    :type="type"
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :id="id"
    :autofocus="autofocus"
    @input="handleInput"
    @focus="emit('focus', $event)"
    @blur="emit('blur', $event)"
  />
</template>
