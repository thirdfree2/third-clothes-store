import type { GuardRule } from "./types";

export const publicRoutes = ["/login", "/backoffice/login", "/unauthorized"];

export const guardedRoutes: GuardRule[] = [
  {
    path: "/cart",
    permissions: ["cart:read"],
  },
  {
    path: "/checkout",
    permissions: ["orders:create"],
  },
  {
    path: "/orders",
    permissions: ["orders:read_own"],
    match: "prefix",
  },
  {
    path: "/profile",
    permissions: ["profile:read"],
    match: "prefix",
  },
  {
    path: "/admin",
    permissions: ["clothes:read"],
    match: "prefix",
    area: "admin",
  },
  {
    path: "/admin/clothes/new",
    permissions: ["clothes:write"],
    area: "admin",
  },
  {
    path: "/backoffice",
    permissions: ["clothes:read"],
    match: "prefix",
    area: "admin",
  },
  {
    path: "/backoffice/clothes/new",
    permissions: ["clothes:write"],
    area: "admin",
  },
];

export function findGuardRule(pathname: string): GuardRule | null {
  if (isPublicRoute(pathname)) {
    return null;
  }

  const sortedRules = [...guardedRoutes].sort((a, b) => b.path.length - a.path.length);

  return (
    sortedRules.find((rule) => {
      if (rule.match === "prefix") {
        return pathname === rule.path || pathname.startsWith(`${rule.path}/`);
      }

      return pathname === rule.path;
    }) ?? null
  );
}

export function isPublicRoute(pathname: string): boolean {
  return publicRoutes.some((route) => pathname === route);
}
