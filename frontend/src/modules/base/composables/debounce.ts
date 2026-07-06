import { ref, type Ref } from "vue";

export function useDebounce(fn: Function, delay: number) {
  let timer: ReturnType<typeof setTimeout> | null = null;
  const remaining: Ref<number> = ref(0);
  let startTime = 0;

  const debounced = (...args: any[]) => {
    if (timer) clearTimeout(timer);

    startTime = Date.now();
    remaining.value = delay;

    // Обновляем счётчик каждые ~30ms для плавности (опционально)
    const interval = setInterval(() => {
      const elapsed = Date.now() - startTime;
      remaining.value = Math.max(0, delay - elapsed);
      if (remaining.value === 0) clearInterval(interval);
    }, 30);

    timer = setTimeout(() => {
      fn(...args);
      timer = null;
      remaining.value = 0;
      clearInterval(interval);
    }, delay);
  };

  debounced.cancel = () => {
    if (timer) {
      clearTimeout(timer);
      timer = null;
      remaining.value = 0;
    }
  };

  return { debounced, remaining };
}
