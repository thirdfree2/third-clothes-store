export type AuthClaims = {
  sub: number | string;
  email?: string;
  user_type?: "CUSTOMER" | "ADMIN" | string;
  role?: string;
  permissions?: string[];
  iss?: string;
  iat?: number;
  exp?: number;
  type?: string;
};

export type AuthArea = "customer" | "admin";

export type GuardRule = {
  path: string;
  permissions: string[];
  match?: "exact" | "prefix";
  area?: AuthArea;
};
