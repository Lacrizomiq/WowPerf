import axios from "axios";
import { authService } from "@/libs/authService";
import {Simulate} from "react-dom/test-utils";
import error = Simulate.error;

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
});

api.interceptors.response.use(
    (response) => response,
    async (error) => {
      if (error.response?.status === 401) {
        try {
          // Try to refresh the token via the refresh token in the cookie
          await authService.refreshToken();
          // Retry the original request
          return axios(error.config);
        } catch (refreshError) {
          // if the refresh failed, redirect to the login page
          window.location.href = "/login";
          return Promise.reject(refreshError);
        }
      }
      return Promise.reject(error);
    }
);


export default api;
