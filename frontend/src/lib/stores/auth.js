import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const initialToken = browser ? localStorage.getItem('token') : null;
const initialUser = browser ? JSON.parse(localStorage.getItem('user')) : null;
const initialAdmin = browser ? JSON.parse(localStorage.getItem('admin')) : null;

export const token = writable(initialToken);
export const user = writable(initialUser);
export const admin = writable(initialAdmin);

if (browser) {
  token.subscribe((v) => {
    if (v) localStorage.setItem('token', v);
    else localStorage.removeItem('token');
  });

  user.subscribe((v) => {
    if (v) localStorage.setItem('user', JSON.stringify(v));
    else localStorage.removeItem('user');
  });

  admin.subscribe((v) => {
    if (v) localStorage.setItem('admin', JSON.stringify(v));
    else localStorage.removeItem('admin');
  });
}

export function logout() {
  token.set(null);
  user.set(null);
  admin.set(null);
}
