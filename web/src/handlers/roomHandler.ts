import { toastService } from "@/services/toastService";

export const messageHandler = (event: MessageEvent) => {
  try {
    const data = JSON.parse(event.data);
    console.log("WebSocket message received:", data);
    // for now just display the message via toast

    toastService.show({
      message: data.message,
      type: data.type || "info",
      dismissable: true,
      durationMs: 5000,
    });
  } catch (error: unknown) {
    console.error("Failed to parse WebSocket message:", error);
  }
};
