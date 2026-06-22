import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const initialTheme = browser ? localStorage.getItem('theme') || 'light' : 'light';

export const theme = writable(initialTheme);

if (browser) {
  theme.subscribe((v) => {
    localStorage.setItem('theme', v);
    if (v === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  });
}
