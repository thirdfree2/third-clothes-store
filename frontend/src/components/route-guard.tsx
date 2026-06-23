"use client";

import { usePathname, useRouter } from "next/navigation";
import { type ReactNode, useEffect, useMemo, useState } from "react";
import {
  decodeAccessToken,
  hasAllPermissions,
  isTokenExpired,
} from "@/lib/auth/jwt";
import { findGuardRule } from "@/lib/auth/route-rules";
import { tokenStore } from "@/lib/auth/token-store";
import type { AuthArea } from "@/lib/auth/types";

type RouteGuardProps = {
  children: ReactNode;
  permissions?: string[];
  area?: AuthArea;
};

export function RouteGuard({ children, permissions, area }: RouteGuardProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [isAllowed, setIsAllowed] = useState(false);
  const [isChecking, setIsChecking] = useState(true);

  const requiredPermissions = useMemo(() => {
    if (permissions) {
      return permissions;
    }

    return findGuardRule(pathname)?.permissions ?? [];
  }, [pathname, permissions]);

  const authArea = useMemo(() => {
    return area ?? findGuardRule(pathname)?.area ?? "customer";
  }, [area, pathname]);

  useEffect(() => {
    const token = tokenStore.get(authArea);
    const claims = token ? decodeAccessToken(token) : null;
    const returnTo = `${pathname}${window.location.search}`;

    if (!token || isTokenExpired(claims)) {
      tokenStore.clear(authArea);
      const loginPath = authArea === "admin" ? "/backoffice/login" : "/login";
      router.replace(`${loginPath}?returnTo=${encodeURIComponent(returnTo)}`);
      return;
    }

    if (!hasAllPermissions(claims, requiredPermissions)) {
      router.replace("/unauthorized");
      return;
    }

    setIsAllowed(true);
    setIsChecking(false);
  }, [authArea, pathname, requiredPermissions, router]);

  if (isChecking || !isAllowed) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-linen px-5">
        <p className="text-sm font-medium text-black/60">Checking access...</p>
      </main>
    );
  }

  return children;
}
