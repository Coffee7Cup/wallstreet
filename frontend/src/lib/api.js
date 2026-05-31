import axios from 'axios';
import { get } from 'svelte/store';
import { token, logout } from '$lib/stores/auth';
import { PUBLIC_HOST } from '$env/static/public';

const api = axios.create({
  baseURL: `http://${PUBLIC_HOST}:3000/api/v1`,
  headers: {
    'Content-Type': 'application/json'
  }
});

api.interceptors.request.use(
  (config) => {
    const t = get(token);
    if (t) {
      config.headers.Authorization = `Bearer ${t}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      logout();
    }
    return Promise.reject(error);
  }
);

export default api;
