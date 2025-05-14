import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Initialize from localStorage if available, otherwise default to false
const storedDarkMode = browser ? localStorage.getItem('darkMode') === 'true' : false;

export const darkMode = writable<boolean>(storedDarkMode);

// Subscribe to changes and update localStorage
if (browser) {
  darkMode.subscribe((value) => {
    localStorage.setItem('darkMode', String(value));
    if (value) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  });
} 