import { apiFetch } from "./client";
import type {
  CustomerAddress,
  CustomerProfile,
  UpdateCustomerAddressInput,
  UpdateCustomerProfileInput,
} from "./types";

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

  getDefaultAddress(token: string) {
    return apiFetch<CustomerAddress>("/api/v1/me/address", {
      token,
      cache: "no-store",
    });
  },

  updateDefaultAddress(token: string, input: UpdateCustomerAddressInput) {
    return apiFetch<CustomerAddress>("/api/v1/me/address", {
      method: "PATCH",
      token,
      body: input,
    });
  },
};
