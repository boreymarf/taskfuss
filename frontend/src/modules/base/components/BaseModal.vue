<script setup lang="ts">
import BaseButton from "./BaseButton.vue";

interface Props {
  showClose?: boolean;
  showConfirm?: boolean;
  confirmText?: string;
  closeOnOverlay?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showClose: true,
  showConfirm: true,
  confirmText: "Confirm",
  closeOnOverlay: true,
});

const model = defineModel<boolean>({ default: false });

const close = () => {
  model.value = false;
};
</script>

<template>
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="fixed inset-0 bg-black/50 flex justify-center items-center z-9999"
      @click.self="closeOnOverlay ? close : undefined"
    >
      <div class="bg-white p-5 max-w-md w-[90%] relative">
        <BaseButton
          v-if="props.showClose"
          class="absolute top-2 right-2"
          @click="close"
        >
          ×
        </BaseButton>

        <div>
          <slot />
        </div>

        <BaseButton v-if="props.showConfirm" @click="close">
          {{ props.confirmText }}
        </BaseButton>
      </div>
    </div>
  </Teleport>
</template>
