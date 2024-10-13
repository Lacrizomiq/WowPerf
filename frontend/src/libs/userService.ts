import api from "./api";

export const userService = {
  async getProfile() {
    try {
      const response = await api.get("/user/profile");
      return response.data;
    } catch (error) {
      console.error("Error fetching user profile:", error);
      throw error;
    }
  },

  async updateEmail(newEmail: string) {
    try {
      const response = await api.put("/user/email", { new_email: newEmail });
      return response.data;
    } catch (error) {
      console.error("Error updating email:", error);
      throw error;
    }
  },

  async changePassword(currentPassword: string, newPassword: string) {
    try {
      const response = await api.put("/user/password", {
        current_password: currentPassword,
        new_password: newPassword,
      });
      return response.data;
    } catch (error) {
      console.error("Error changing password:", error);
      throw error;
    }
  },

  async changeUsername(newUsername: string) {
    try {
      const response = await api.put("/user/username", {
        new_username: newUsername,
      });
      return response.data;
    } catch (error) {
      console.error("Error changing username:", error);
      throw error;
    }
  },

  async deleteAccount() {
    try {
      const response = await api.delete("/user/account");
      return response.data;
    } catch (error) {
      console.error("Error deleting account:", error);
      throw error;
    }
  },
};
