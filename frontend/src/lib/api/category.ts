import { apiFetch } from "./client";
import type { Category } from "./types";

export const categoriesApi = {
    list(token: string) {
        return apiFetch<Category[]>("/api/v1/categories", {
            method: "GET",
            token,
            cache: "no-store",
        });
    },
};
