import type { AuthClaims } from "./types";

export function decodeAccessToken(token: string): AuthClaims | null {
  const parts = token.split(".");
  if (parts.length < 2) {
    return null;
  }

  try {
    return JSON.parse(base64URLDecode(parts[1])) as AuthClaims;
  } catch {
    return null;
  }
}

export function isTokenExpired(claims: AuthClaims | null): boolean {
  if (!claims?.exp) {
    return true;
  }

  return claims.exp * 1000 <= Date.now();
}

export function hasAllPermissions(
  claims: AuthClaims | null,
  permissions: string[],
): boolean {
  if (permissions.length === 0) {
    return true;
  }

  const userPermissions = new Set(claims?.permissions ?? []);
  return permissions.every((permission) => userPermissions.has(permission));
}

function base64URLDecode(value: string): string {
  const base64 = value.replace(/-/g, "+").replace(/_/g, "/");
  const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), "=");

  if (typeof window === "undefined") {
    return Buffer.from(padded, "base64").toString("utf8");
  }

  return decodeURIComponent(
    window
      .atob(padded)
      .split("")
      .map((char) => `%${char.charCodeAt(0).toString(16).padStart(2, "0")}`)
      .join(""),
  );
}
