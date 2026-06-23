import { apiFetch } from "./client";
import type { CustomerProfile, UpdateCustomerProfileInput } from "./types";

export const customerApi = {
  getMe(token: string) {
    return apiFetch<CustomerProfile>("/api/v1/me", {
      token,
      cache: "no-store",
    });
  },

  updateMe(token: string, input: UpdateCustomerProfileInput) {
    return apiFetch<CustomerProfile>("/api/v1/me", {
      method: "PATCH",
      token,
      body: input,
    });
  },
};
