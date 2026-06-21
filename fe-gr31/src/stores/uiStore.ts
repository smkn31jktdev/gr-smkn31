import { writable } from 'svelte/store';

export interface ToastItem {
  id: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
}

export const toasts = writable<ToastItem[]>([]);
export const sidebarCollapsed = writable<boolean>(false);

export function addToast(message: string, type: 'info' | 'success' | 'warning' | 'error' = 'info') {
  const id = Math.random().toString(36).substring(2, 9);
  toasts.update(($toasts) => [...$toasts, { id, message, type }]);
  
  setTimeout(() => {
    toasts.update(($toasts) => $toasts.filter((t) => t.id !== id));
  }, 4000);
}
