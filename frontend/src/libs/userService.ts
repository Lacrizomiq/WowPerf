// src/libs/userService.ts
import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";

// API response types aligned with the backend
export interface UserProfile {
  id: number;
  username: string;
  email: string;
  battle_tag?: string;
  battle_net_id?: string;
  favorite_character_id?: number;
}

export interface ApiResponse {
  message: string;
  code: string;
  error?: string;
}

// Standardized error codes
export enum UserErrorCode {
  // Validation errors
  INVALID_EMAIL = "invalid_email",
  INVALID_PASSWORD = "invalid_password",
  INVALID_USERNAME = "invalid_username",

  // Conflict errors
  EMAIL_EXISTS = "email_exists",
  USERNAME_EXISTS = "username_exists",

  // Security errors
  UNAUTHORIZED = "unauthorized",
  INVALID_CSRF_TOKEN = "INVALID_CSRF_TOKEN",

  // Limitation errors
  USERNAME_CHANGE_LIMIT = "username_change_limit",

  // Technical errors
  PROFILE_NOT_FOUND = "profile_not_found",
  NETWORK_ERROR = "network_error",
  SERVER_ERROR = "server_error",
}

export class UserServiceError extends Error {
  constructor(
    public code: UserErrorCode,
    message: string,
    public originalError?: unknown
  ) {
    super(message);
    this.name = "UserServiceError";
  }
}

export const userService = {
  // Route protected by JWT only
  async getProfile(): Promise<UserProfile> {
    try {
      const response = await api.get<UserProfile>("/user/profile");
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        if (err.response?.status === 401) {
          throw new UserServiceError(
            UserErrorCode.UNAUTHORIZED,
            "Not authorized to access profile"
          );
        }

        if (err.response?.status === 404) {
          throw new UserServiceError(
            UserErrorCode.PROFILE_NOT_FOUND,
            "Profile not found"
          );
        }
      }
      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Failed to fetch profile"
      );
    }
  },

  // Route protected by JWT + CSRF
  async updateEmail(newEmail: string): Promise<ApiResponse> {
    try {
      const response = await api.put<ApiResponse>("/user/email", {
        new_email: newEmail,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        if (err.response?.data?.code === "INVALID_CSRF_TOKEN") {
          throw new UserServiceError(
            UserErrorCode.INVALID_CSRF_TOKEN,
            "Security verification failed"
          );
        }

        switch (err.response?.data?.code) {
          case "invalid_email":
            throw new UserServiceError(
              UserErrorCode.INVALID_EMAIL,
              "Invalid email format"
            );
          case "email_exists":
            throw new UserServiceError(
              UserErrorCode.EMAIL_EXISTS,
              "Email already in use"
            );
          default:
            throw new UserServiceError(
              UserErrorCode.SERVER_ERROR,
              err.response?.data?.error || "Failed to update email"
            );
        }
      }
      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while updating email"
      );
    }
  },

  // Route protected by JWT + CSRF
  async changePassword(
    currentPassword: string,
    newPassword: string
  ): Promise<ApiResponse> {
    try {
      const response = await api.put<ApiResponse>("/user/password", {
        current_password: currentPassword,
        new_password: newPassword,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        if (err.response?.data?.code === "INVALID_CSRF_TOKEN") {
          throw new UserServiceError(
            UserErrorCode.INVALID_CSRF_TOKEN,
            "Security verification failed"
          );
        }

        switch (err.response?.data?.code) {
          case "invalid_current_password":
            throw new UserServiceError(
              UserErrorCode.INVALID_PASSWORD,
              "Current password is incorrect"
            );
          case "invalid_new_password":
            throw new UserServiceError(
              UserErrorCode.INVALID_PASSWORD,
              "New password does not meet requirements"
            );
          default:
            throw new UserServiceError(
              UserErrorCode.SERVER_ERROR,
              "Failed to change password"
            );
        }
      }
      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while changing password"
      );
    }
  },

  // Route protected by JWT + CSRF
  async changeUsername(newUsername: string): Promise<ApiResponse> {
    try {
      const response = await api.put<ApiResponse>("/user/username", {
        new_username: newUsername,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        if (err.response?.data?.code === "INVALID_CSRF_TOKEN") {
          throw new UserServiceError(
            UserErrorCode.INVALID_CSRF_TOKEN,
            "Security verification failed"
          );
        }

        if (err.response?.status === 429) {
          throw new UserServiceError(
            UserErrorCode.USERNAME_CHANGE_LIMIT,
            "Username can only be changed once every 30 days"
          );
        }

        switch (err.response?.data?.code) {
          case "username_exists":
            throw new UserServiceError(
              UserErrorCode.USERNAME_EXISTS,
              "Username already taken"
            );
          case "invalid_username":
            throw new UserServiceError(
              UserErrorCode.INVALID_USERNAME,
              "Invalid username format"
            );
          default:
            throw new UserServiceError(
              UserErrorCode.SERVER_ERROR,
              "Failed to change username"
            );
        }
      }
      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while changing username"
      );
    }
  },

  // Route protected by JWT + CSRF
  async deleteAccount(): Promise<ApiResponse> {
    try {
      const response = await api.delete<ApiResponse>("/user/account");
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        if (err.response?.data?.code === "INVALID_CSRF_TOKEN") {
          throw new UserServiceError(
            UserErrorCode.INVALID_CSRF_TOKEN,
            "Security verification failed"
          );
        }

        if (err.response?.status === 401) {
          throw new UserServiceError(
            UserErrorCode.UNAUTHORIZED,
            "Not authorized to delete account"
          );
        }
      }
      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while deleting account"
      );
    }
  },
};
