import toast from "react-hot-toast";
import type { ToastOptions } from "react-hot-toast";
type ToastType = "success" | "error" | "info" | "warning";

interface ShowToastOptions {
  type?: ToastType;
  message: string;
  dismissable?: boolean;
  durationMs?: number;
}

/**
 * Unified toast service
 * Example: toastService.show({ message: "Saved!", type: "success" })
 */

export const toastService = {
  show({
    message,
    type = "info",
    dismissable = true,
    durationMs = 3000,
  }: ShowToastOptions) {
    const options: ToastOptions = {
      duration: durationMs,
      id: dismissable ? undefined : "persistent-toast",
    };

    switch (type) {
      case "success":
        toast.success(message, options);
        break;
      case "error":
        toast.error(message, options);
        break;
      case "warning":
        toast(message, { ...options, icon: "⚠️" });
        break;
      default:
        toast(message, options);
        break;
    }
  },

  dismiss(id?: string) {
    toast.dismiss(id);
  },
};
