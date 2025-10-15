// src/lib/api/apiClient.ts
import type { AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import axios from 'axios'

// Extend the config type to include your custom property
interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean
}

const apiClient = axios.create({
  // eslint-disable-next-line node/prefer-global/process
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
})

apiClient.interceptors.request.use(
  (config) => {
    // Add any auth headers or other request modifications here
    // const token = getAuthToken();
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    console.warn('Request:', {
      url: config.url,
      fullURL: config.baseURL + config.url,
      withCredentials: config.withCredentials,
      headers: config.headers,
      cookies: document.cookie,
    })
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

let isRetrying = false

apiClient.interceptors.response.use(
  response => response,
  async (error) => {
    const status = error.response?.status
    const config = error.config as CustomAxiosRequestConfig

    if (status === 401) {
      // Retry once if this is the first 401 and we haven't retried yet
      if (!isRetrying && !config._retry) {
        isRetrying = true
        config._retry = true

        // Wait a moment for cookies to settle
        await new Promise(resolve => setTimeout(resolve, 100))

        try {
          const result = await apiClient.request(config)
          isRetrying = false
          return result
        }
        catch (retryError) {
          isRetrying = false
          // If retry fails, then redirect
          console.error('Unauthorized access - redirecting to login')
          document.cookie.split(';').forEach((cookie) => {
            const eqPos = cookie.indexOf('=')
            const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
            document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
          })
          window.location.href = '/login'
          return Promise.reject(retryError)
        }
      }
      else {
        // Second 401 or already retrying - redirect immediately
        console.error('Unauthorized access - redirecting to login')
        document.cookie.split(';').forEach((cookie) => {
          const eqPos = cookie.indexOf('=')
          const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim()
          document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
        })
        window.location.href = '/login'
      }
    }
    else if (status === 403) {
      console.error('Forbidden access')
    }
    else if (status === 404) {
      console.error('Resource not found')
    }
    else if (status >= 500) {
      console.error('Server error occurred')
    }
    else {
      console.error('An error occurred:', error.message)
    }

    return Promise.reject(error)
  },
)

export function customAxios<T>(config: AxiosRequestConfig): Promise<T> {
  return Promise.resolve(
    apiClient({
      ...config,
      withCredentials: true, // Force this to always be true
    }),
  ).then((response: AxiosResponse<T>) => response.data)
}

export default apiClient
