import { onScopeDispose } from 'vue'

export function useDebounce<T extends (...args: any[]) => any>(
  fn: T,
  delay = 300
) {
  let timer: ReturnType<typeof setTimeout> | null = null

  const debounced = (...args: Parameters<T>) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => fn(...args), delay)
  }

  const cancel = () => {
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  // This functions is automatically called when the composable is removed 
  onScopeDispose(cancel) 

  return { debounced, cancel }
}