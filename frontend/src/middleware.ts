import { NextRequest, NextResponse } from "next/server";
import { findGuardRule, isPublicRoute } from "@/lib/auth/route-rules";
import { accessTokenCookieNames } from "@/lib/auth/constants";
import type { AuthClaims } from "@/lib/auth/types";

export function middleware(request: NextRequest) {
  const { pathname, search } = request.nextUrl;
  if (isPublicRoute(pathname)) {
    return NextResponse.next();
  }

  const rule = findGuardRule(pathname);

  if (!rule) {
    return NextResponse.next();
  }

  const area = rule.area ?? "customer";
  const token = request.cookies.get(accessTokenCookieNames[area])?.value;
  const claims = token ? decodeTokenPayload(token) : null;

  if (!token || isExpired(claims)) {
    const loginURL = new URL(area === "admin" ? "/backoffice/login" : "/login", request.url);
    loginURL.searchParams.set("returnTo", `${pathname}${search}`);
    return NextResponse.redirect(loginURL);
  }

  if (!hasPermissions(claims, rule.permissions)) {
    return NextResponse.redirect(new URL("/unauthorized", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/cart",
    "/checkout",
    "/orders/:path*",
    "/profile/:path*",
    "/admin/:path*",
    "/backoffice/:path*",
  ],
};

function decodeTokenPayload(token: string): AuthClaims | null {
  const payload = token.split(".")[1];
  if (!payload) {
    return null;
  }

  try {
    return JSON.parse(base64URLDecode(payload)) as AuthClaims;
  } catch {
    return null;
  }
}

function base64URLDecode(value: string): string {
  const base64 = value.replace(/-/g, "+").replace(/_/g, "/");
  const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), "=");
  const binary = atob(padded);
  const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));

  return new TextDecoder().decode(bytes);
}

function isExpired(claims: AuthClaims | null): boolean {
  if (!claims?.exp) {
    return true;
  }

  return claims.exp * 1000 <= Date.now();
}

function hasPermissions(claims: AuthClaims | null, permissions: string[]): boolean {
  const userPermissions = new Set(claims?.permissions ?? []);
  return permissions.every((permission) => userPermissions.has(permission));
}
