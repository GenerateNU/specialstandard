// src/lib/api/apiClient.ts
import axios, { AxiosRequestConfig, AxiosResponse } from "axios";

const apiClient = axios.create({
  // eslint-disable-next-line node/prefer-global/process
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

apiClient.interceptors.request.use(
  (config) => {
    // Add any auth headers or other request modifications here
    // const token = getAuthToken();
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;

    if (status === 401) {
      // Handle unauthorized access
      console.error("Unauthorized access");

      // Clear cookies and redirect to homepage after a delay
      setTimeout(() => {
        document.cookie.split(";").forEach((cookie) => {
          const eqPos = cookie.indexOf("=");
          const name = eqPos > -1 ? cookie.substring(0, eqPos) : cookie;
          document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
        });
        window.location.href = "/";
      }, 2000);
    } else if (status === 403) {
      console.error("Forbidden access");
    } else if (status === 404) {
      console.error("Resource not found");
    } else if (status >= 500) {
      console.error("Server error occurred");
    } else {
      console.error("An error occurred:", error.message);
    }

    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    return Promise.reject(error);
  }
);

export const customAxios = <T>(config: AxiosRequestConfig): Promise<T> => {
  return apiClient(config).then((response: AxiosResponse<T>) => response.data);
};
