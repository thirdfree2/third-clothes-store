import { apiFetch } from "./client";
import type { Category } from "./types";

export const categoriesApi = {
    list(token?: string | null) {
        return apiFetch<Category[]>("/api/v1/categories", {
            method: "GET",
            ...(token ? { token } : {}),
            cache: "no-store",
        });
    },
};
