import { apiFetch } from "./client";
import type { LoginResponse, RegisterResponse } from "./types";

export const authApi = {
  login(input: {
    email: string;
    password: string;
    userType?: "CUSTOMER" | "ADMIN" | string;
  }) {
    return apiFetch<LoginResponse>("/api/v1/auth/login", {
      method: "POST",
      body: {
        email: input.email,
        password: input.password,
        userType: input.userType ?? "CUSTOMER",
      },
    });
  },

  register(input: { email: string; password: string }) {
    return apiFetch<RegisterResponse>("/api/v1/auth/register", {
      method: "POST",
      body: input,
    });
  },
};
