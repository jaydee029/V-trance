// import axios from 'axios';
// import { store } from '../store/store';
// import { updateAccessToken, logout } from '../store/authSlice';

// const api = axios.create({
//   baseURL: 'http://localhost:8080',
//   withCredentials: true
// });

// api.interceptors.request.use((config) => {
//   const state = store.getState();
//   const accessToken = state.authController.accessToken;
//   const refreshToken = state.authController.refreshToken;

//   if (config.url.includes('/signup') || config.url.includes('/login')) {
//     return config;
//   }

//   if (accessToken) {
//     config.headers.Authorization = `Bearer ${accessToken}`;
//   }
//   if (refreshToken) {
//     config.headers['X-Refresh-Token'] = refreshToken;
//   }
//   return config;
// });

// api.interceptors.response.use(
//   (response) => response,
//   async (error) => {
//     const originalRequest = error.config;

//     if (
//       error.response?.status === 401 &&
//       !originalRequest._retry ) {
//       originalRequest._retry = true;
//       try {
//         const res = await axios.get('http://localhost:8080/user/refresh', {
//           headers: {
//             'X-Refresh-Token': store.getState().authController.refreshToken,
//           },
//         });

//         if (res.data.accessToken) {
//           store.dispatch(updateAccessToken(res.data.accessToken));

//           originalRequest.headers.Authorization = `Bearer ${res.data.accessToken}`;
//           return axios(originalRequest);
//         }
//       } catch (refreshErr) {
//         store.dispatch(logout());
//         return Promise.reject(refreshErr);
//       }
//     }

//     return Promise.reject(error);
//   }
// );

// export default api;

import { store } from '../store/store';
import { updateAccessToken, logout } from '../store/authSlice';

const BASE_URL = 'http://localhost:8080';

const apiFetch = async (endpoint, options = {}, skipAuth = false) => {
  const state = store.getState();
  const accessToken = state.authController.accessToken;
  const refreshToken = state.authController.refreshToken;

  const url = `${BASE_URL}${endpoint}`;
  const headers = {
    ...options.headers,
    'Content-Type': 'application/json',
  };

  if (!skipAuth) {
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }
    if (refreshToken) {
      headers['X-Refresh-Token'] = refreshToken;
    }
  }

  let response = await fetch(url, {
    ...options,
    credentials: 'include',
    headers,
    mode: 'no-cors'
  });

  // Handle token refresh logic
  if (response.status === 401 && !options._retry) {
    options._retry = true;

    try {
      const refreshRes = await fetch(`${BASE_URL}/user/refresh`, {
        method: 'GET',
        headers: {
          'X-Refresh-Token': refreshToken,
        },
        credentials: 'include',
        mode: 'no-cors'
      });

      const refreshData = await refreshRes.json();

      if (refreshRes.ok && refreshData.accessToken) {
        store.dispatch(updateAccessToken(refreshData.accessToken));

        // Retry original request with new access token
        return apiFetch(endpoint, {
          ...options,
          headers: {
            ...headers,
            Authorization: `Bearer ${refreshData.accessToken}`,
          },
          _retry: true,
        });
      } else {
        store.dispatch(logout());
        throw new Error('Refresh failed');
      }
    } catch (err) {
      store.dispatch(logout());
      throw err;
    }
  }

  if (!response.ok) {
    const errData = await response.json().catch(() => ({}));
    const error = new Error('API Error');
    error.data = errData;
    throw error;
  }

  return response.json();
};

export default apiFetch;