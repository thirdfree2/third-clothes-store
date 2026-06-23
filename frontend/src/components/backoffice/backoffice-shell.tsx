"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { type ReactNode } from "react";
import { RouteGuard } from "@/components/route-guard";
import { tokenStore } from "@/lib/auth/token-store";

type BackofficeShellProps = {
  children: ReactNode;
};

export function BackofficeShell({ children }: BackofficeShellProps) {
  const router = useRouter();
  const pathname = usePathname();

  function handleLogout() {
    tokenStore.clear("admin");
    router.replace("/backoffice/login");
  }

  const navItems = [
    { href: "/backoffice", label: "Clothes", exact: true },
    { href: "/backoffice/orders", label: "Orders" },
  ];

  return (
    <RouteGuard area="admin" permissions={["clothes:read"]}>
      <div className="min-h-screen bg-linen text-ink">
        <aside className="fixed inset-y-0 left-0 hidden w-64 border-r border-black/10 bg-white px-5 py-6 lg:block">
          <Link href="/backoffice" className="block">
            <p className="text-sm font-medium uppercase tracking-wide text-moss">
              Backoffice
            </p>
            <h1 className="mt-2 text-xl font-semibold">Clothes Store</h1>
          </Link>

          <nav className="mt-8 space-y-1">
            {navItems.map((item) => {
              const isActive = item.exact
                ? pathname === item.href
                : pathname === item.href || pathname.startsWith(`${item.href}/`);

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`block rounded-md px-3 py-2 text-sm font-medium ${
                    isActive
                      ? "bg-ink text-white"
                      : "text-black/65 hover:bg-stone-100 hover:text-ink"
                  }`}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </aside>

        <div className="lg:pl-64">
          <header className="sticky top-0 z-10 border-b border-black/10 bg-linen/95 px-5 py-4 backdrop-blur">
            <div className="mx-auto flex max-w-7xl items-center justify-between gap-4">
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-moss lg:hidden">
                  Backoffice
                </p>
                <h2 className="text-lg font-semibold">Product operations</h2>
              </div>
              <button
                type="button"
                onClick={handleLogout}
                className="h-10 rounded-md border border-black/15 bg-white px-4 text-sm font-medium shadow-soft hover:border-black/30"
              >
                Logout
              </button>
            </div>
          </header>

          {children}
        </div>
      </div>
    </RouteGuard>
  );
}
