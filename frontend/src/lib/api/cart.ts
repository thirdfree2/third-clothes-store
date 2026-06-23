import { apiFetch } from "./client";
import type { Cart } from "./types";

export const cartApi = {
  get(token: string) {
    return apiFetch<Cart>("/api/v1/cart", {
      token,
      cache: "no-store",
    });
  },

  addItem(
    token: string,
    input: { clothes_id: number; quantity: number; size?: string },
  ) {
    return apiFetch<Cart>("/api/v1/cart/items", {
      method: "POST",
      token,
      body: {
        clothes_id: input.clothes_id,
        quantity: input.quantity,
        size: input.size ?? "",
      },
    });
  },

  updateItem(token: string, itemId: number, input: { quantity: number }) {
    return apiFetch<Cart>(`/api/v1/cart/items/${itemId}`, {
      method: "PATCH",
      token,
      body: input,
    });
  },

  deleteItem(token: string, itemId: number) {
    return apiFetch<Cart>(`/api/v1/cart/items/${itemId}`, {
      method: "DELETE",
      token,
    });
  },

  clear(token: string) {
    return apiFetch<Cart>("/api/v1/cart", {
      method: "DELETE",
      token,
    });
  },
};
