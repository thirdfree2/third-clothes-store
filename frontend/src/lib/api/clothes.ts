import { apiFetch, apiFormDataFetch } from "./client";
import type { Clothes, ClothesImage, PaginatedPayload } from "./types";

export const clothesApi = {
  list(params: { page?: number; perPage?: number; category_id?: number } = {}) {
    return apiFetch<PaginatedPayload<Clothes>>("/api/v1/clothes", {
      query: params,
      cache: "no-store",
    });
  },

  getById(id: number) {
    return apiFetch<Clothes>(`/api/v1/clothes/${id}`, {
      cache: "no-store",
    });
  },

  create(
    token: string,
    input: {
      name: string;
      price: number;
      color_id: number | null;
      category_ids: number[];
    },
  ) {
    return apiFetch<Clothes>("/api/v1/clothes", {
      method: "POST",
      token,
      body: input,
    });
  },

  update(
    id: number,
    token: string,
    input: {
      name: string;
      price: number;
      color_id: number | null;
      category_ids: number[];
    },
  ) {
    return apiFetch<Clothes>(`/api/v1/clothes/${id}`, {
      method: "PUT",
      token,
      body: input,
    });
  },

  delete(id: number, token: string) {
    return apiFetch<{ deleted: boolean }>(`/api/v1/clothes/${id}`, {
      method: "DELETE",
      token,
    });
  },

  uploadImages(id: number, token: string, files: File[]) {
    const body = new FormData();

    for (const file of files) {
      body.append("files", file);
    }

    return apiFormDataFetch<ClothesImage[]>(`/api/v1/clothes/${id}/images`, {
      method: "POST",
      token,
      body,
    });
  },

  deleteImage(id: number, imageID: number, token: string) {
    return apiFetch<{ deleted: boolean }>(`/api/v1/clothes/${id}/images/${imageID}`, {
      method: "DELETE",
      token,
    });
  },
};
