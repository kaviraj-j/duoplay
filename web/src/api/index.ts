import api from "@/lib/axios";

export const healthCheck = async () => {
  const response = await api.get("/health");
  return response.data;
};


export * from "./user";
export * from "./game";
export * from "./room";
