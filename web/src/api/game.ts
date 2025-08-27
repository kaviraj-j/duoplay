import api from "@/lib/axios";
import type { GameListPayload } from "@/types";

export const gameApi = {
  getGamesList: async (): Promise<{
    type: string;
    data: GameListPayload[];
  }> => {
    const response = await api.get("/game/list");
    return response.data as { type: string; data: GameListPayload[] };
  },
};
