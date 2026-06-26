import { apiFetch } from "./client";
import type { Color } from "./types";

export const colorApi = {
  list(token?: string | null) {
    return apiFetch<Color[]>("/api/v1/colors", {
      method: "GET",
      ...(token ? { token } : {}),
      cache: "no-store",
    });
  },
};
