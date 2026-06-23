import { apiFetch } from "./client";
import type { Color } from "./types";

export const colorApi = {
  list(token: string) {
    return apiFetch<Color[]>("/api/v1/colors", {
      method: "GET",
      token,
      cache: "no-store",
    });
  },
};
