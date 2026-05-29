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

const cache = new Map();
const CACHE_TTL = 30000; // 30 seconds

api.interceptors.request.use(
  (config) => {
    const t = get(token);
    if (t) {
      config.headers.Authorization = `Bearer ${t}`;
    }

    if (config.method === 'get') {
      const cacheKey = config.url + JSON.stringify(config.params || {});
      const cachedData = cache.get(cacheKey);
      if (cachedData && Date.now() - cachedData.timestamp < CACHE_TTL) {
        config.adapter = (cfg) => {
          return Promise.resolve({
            data: cachedData.data,
            status: 200,
            statusText: 'OK',
            headers: cfg.headers,
            config: cfg,
            request: {}
          });
        };
      }
    }

    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => {
    if (response.config.method === 'get') {
      const cacheKey = response.config.url + JSON.stringify(response.config.params || {});
      cache.set(cacheKey, {
        data: response.data,
        timestamp: Date.now()
      });
    }
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      logout();
    }
    return Promise.reject(error);
  }
);

export default api;
