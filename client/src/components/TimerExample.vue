<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue';

const seconds = ref<number>(0);
const timerId = ref<number | null>(null);

const formattedTime = computed<string>(() => {
  const mins: number = Math.floor(seconds.value / 60);
  const secs: number = seconds.value % 60;
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
});

function startTimer(): void {
  if (!timerId.value) {
    timerId.value = window.setInterval(() => {
      seconds.value++;
    }, 1000);
  }
}

function stopTimer(): void {
  if (timerId.value) {
    clearInterval(timerId.value);
    timerId.value = null;
  }
}

function resetTimer(): void {
  stopTimer();
  seconds.value = 0;
}

onUnmounted((): void => {
  stopTimer();
});
</script>

<template>
  <div>
    <div class="timer-display">{{ formattedTime }}</div>
    <div class="controls">
      <button @click="startTimer" :disabled="timerId !== null">
        Start
      </button>
      <button @click="stopTimer" :disabled="timerId === null">
        Stop
      </button>
      <button @click="resetTimer">
        Reset
      </button>
    </div>
  </div>
</template>
