import { apiFetch } from "./client";
import type { Cart, Purchase } from "./types";

export const purchasesApi = {
  list(token: string) {
    return apiFetch<Purchase[]>("/api/v1/purchases", {
      token,
      cache: "no-store",
    });
  },

  create(
    token: string,
    input: {
      items: { clothes_id: number; quantity: number; size?: string }[];
    },
  ) {
    return apiFetch<Purchase>("/api/v1/purchases", {
      method: "POST",
      token,
      body: {
        items: input.items.map((item) => ({
          clothes_id: item.clothes_id,
          quantity: item.quantity,
          size: item.size ?? "",
        })),
      },
    });
  },

  createFromCart(token: string, cart: Cart) {
    return purchasesApi.create(token, {
      items: cart.items.map((item) => ({
        clothes_id: item.clothes_id,
        quantity: item.quantity,
        size: item.size,
      })),
    });
  },
};
