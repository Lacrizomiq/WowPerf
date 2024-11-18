// src/libs/userService.ts
import api, { APIError } from "./api";
import axios, { AxiosError } from "axios";

// API response types
interface UserProfile {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

interface ApiResponse {
  message: string;
  error?: string;
  code?: string;
}

interface UpdateEmailResponse extends ApiResponse {
  email: string;
}

interface UpdateUsernameResponse extends ApiResponse {
  username: string;
}

interface UpdatePasswordResponse extends ApiResponse {}

interface DeleteAccountResponse extends ApiResponse {}

// Error codes
export enum UserErrorCode {
  PROFILE_NOT_FOUND = "profile_not_found",
  INVALID_EMAIL = "invalid_email",
  EMAIL_EXISTS = "email_exists",
  INVALID_PASSWORD = "invalid_password",
  INVALID_USERNAME = "invalid_username",
  USERNAME_EXISTS = "username_exists",
  USERNAME_CHANGE_LIMIT = "username_change_limit",
  UNAUTHORIZED = "unauthorized",
  NETWORK_ERROR = "network_error",
  UNKNOWN_ERROR = "unknown_error",
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
  async getProfile(): Promise<UserProfile> {
    try {
      const response = await api.get<UserProfile>("/user/profile");
      return response.data;
    } catch (error) {
      console.error("Error fetching user profile:", error);

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        // Check HTTP status
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

        throw new UserServiceError(
          UserErrorCode.UNKNOWN_ERROR,
          err.response?.data?.error || "Failed to fetch profile",
          err
        );
      }

      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while fetching profile",
        error
      );
    }
  },

  async updateEmail(newEmail: string): Promise<UpdateEmailResponse> {
    try {
      const response = await api.put<UpdateEmailResponse>("/user/email", {
        new_email: newEmail,
      });
      return response.data;
    } catch (error) {
      console.error("Error updating email:", error);

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        // Check HTTP status
        if (err.response?.status === 401) {
          throw new UserServiceError(
            UserErrorCode.UNAUTHORIZED,
            "Not authorized to update email"
          );
        }

        // Check custom error codes
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
              UserErrorCode.UNKNOWN_ERROR,
              err.response?.data?.error || "Failed to update email",
              err
            );
        }
      }

      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while updating email",
        error
      );
    }
  },

  async changePassword(
    currentPassword: string,
    newPassword: string
  ): Promise<UpdatePasswordResponse> {
    try {
      const response = await api.put<UpdatePasswordResponse>("/user/password", {
        current_password: currentPassword,
        new_password: newPassword,
      });
      return response.data;
    } catch (error) {
      console.error("Error changing password:", error);

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        // Check HTTP status
        if (err.response?.status === 401) {
          throw new UserServiceError(
            UserErrorCode.UNAUTHORIZED,
            "Not authorized to change password"
          );
        }

        // Check custom error codes
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
              UserErrorCode.UNKNOWN_ERROR,
              err.response?.data?.error || "Failed to change password",
              err
            );
        }
      }

      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while changing password",
        error
      );
    }
  },

  async changeUsername(newUsername: string): Promise<UpdateUsernameResponse> {
    try {
      const response = await api.put<UpdateUsernameResponse>("/user/username", {
        new_username: newUsername,
      });
      return response.data;
    } catch (error) {
      console.error("Error changing username:", error);

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        // Check HTTP status
        if (err.response?.status === 401) {
          throw new UserServiceError(
            UserErrorCode.UNAUTHORIZED,
            "Not authorized to change username"
          );
        }

        if (err.response?.status === 429) {
          throw new UserServiceError(
            UserErrorCode.USERNAME_CHANGE_LIMIT,
            "Username can only be changed once every 30 days"
          );
        }

        // Check custom error codes
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
              UserErrorCode.UNKNOWN_ERROR,
              err.response?.data?.error || "Failed to change username",
              err
            );
        }
      }

      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while changing username",
        error
      );
    }
  },

  async deleteAccount(): Promise<DeleteAccountResponse> {
    try {
      const response = await api.delete<DeleteAccountResponse>("/user/account");
      return response.data;
    } catch (error) {
      console.error("Error deleting account:", error);

      if (axios.isAxiosError(error)) {
        const err = error as AxiosError<APIError>;

        // Check HTTP status
        if (err.response?.status === 401) {
          throw new UserServiceError(
            UserErrorCode.UNAUTHORIZED,
            "Not authorized to delete account"
          );
        }

        throw new UserServiceError(
          UserErrorCode.UNKNOWN_ERROR,
          err.response?.data?.error || "Failed to delete account",
          err
        );
      }

      throw new UserServiceError(
        UserErrorCode.NETWORK_ERROR,
        "Network error while deleting account",
        error
      );
    }
  },
};

// Export types
export type {
  UserProfile,
  ApiResponse,
  UpdateEmailResponse,
  UpdateUsernameResponse,
  UpdatePasswordResponse,
  DeleteAccountResponse,
};
