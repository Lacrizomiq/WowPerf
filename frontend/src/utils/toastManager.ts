// utils/toastManager.ts
import toast, { Toast, ToastOptions } from "react-hot-toast";

// IDs constants to avoid duplicate toasts
export const TOAST_IDS = {
  // Battle.net connection
  BATTLENET_LINKING: "battlenet-link-required",
  BATTLENET_LINK_SUCCESS: "battlenet-link-success",
  BATTLENET_LINK_ERROR: "battlenet-link-error",
  BATTLENET_UNLINK_SUCCESS: "battlenet-unlink-success",

  // Characters
  CHARACTERS_SYNC: "characters-sync",
  CHARACTERS_SYNC_SUCCESS: "characters-sync-success",
  CHARACTERS_SYNC_ERROR: "characters-sync-error",
  CHARACTERS_REFRESH: "characters-refresh",
  CHARACTERS_REFRESH_SUCCESS: "characters-refresh-success",
  CHARACTERS_REFRESH_ERROR: "characters-refresh-error",
  CHARACTER_FAVORITE: "character-favorite",
  CHARACTER_TOGGLE: "character-toggle",
};

// Standard options for all toasts
const defaultOptions: ToastOptions = {
  duration: 4000, // 4 seconds by default
  position: "top-center",
  style: {
    maxWidth: "500px",
  },
};

// Utility function to display error messages
export function showError(message: string, id?: string) {
  // If a toast with this ID already exists, delete it
  if (id) toast.dismiss(id);

  return toast.error(message, {
    ...defaultOptions,
    id: id,
  });
}

// Utility function to display success messages
export function showSuccess(message: string, id?: string) {
  // If a toast with this ID already exists, delete it
  if (id) toast.dismiss(id);

  return toast.success(message, {
    ...defaultOptions,
    id: id,
  });
}

// Utility function to display standard notifications
export function showInfo(message: string, id?: string) {
  // If a toast with this ID already exists, delete it
  if (id) toast.dismiss(id);

  return toast(message, {
    ...defaultOptions,
    id: id,
  });
}

// Function to delete toasts
export function dismissToast(id: string) {
  toast.dismiss(id);
}
