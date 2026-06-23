"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { ApiError } from "@/lib/api/client";
import { customerApi } from "@/lib/api/customer";
import { purchasesApi } from "@/lib/api/purchases";
import type { CustomerProfile, Purchase } from "@/lib/api/types";
import { decodeAccessToken } from "@/lib/auth/jwt";
import { tokenStore } from "@/lib/auth/token-store";

type ProfileState =
  | { status: "loading"; profile: null; error: null }
  | { status: "ready"; profile: CustomerProfile; error: null }
  | { status: "not-found"; profile: null; error: string }
  | { status: "error"; profile: null; error: string };

type PurchaseHistoryState =
  | { status: "loading"; purchases: Purchase[]; error: null }
  | { status: "ready"; purchases: Purchase[]; error: null }
  | { status: "error"; purchases: Purchase[]; error: string };

export function CustomerProfileDetail() {
  const router = useRouter();
  const [state, setState] = useState<ProfileState>({
    status: "loading",
    profile: null,
    error: null,
  });
  const [email, setEmail] = useState<string | null>(null);
  const [isUpdating, setIsUpdating] = useState(false);
  const [historyState, setHistoryState] = useState<PurchaseHistoryState>({
    status: "loading",
    purchases: [],
    error: null,
  });

  useEffect(() => {
    void loadProfile();
    void loadPurchaseHistory();
  }, []);

  async function loadProfile() {
    const token = tokenStore.get("customer");

    if (!token) {
      setState({
        status: "error",
        profile: null,
        error: "Customer token is missing",
      });
      return;
    }

    setEmail(decodeAccessToken(token)?.email ?? null);

    try {
      const profile = await customerApi.getMe(token);
      setState({ status: "ready", profile, error: null });
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Cannot load profile";

      setState({
        status: isProfileNotFound(err) ? "not-found" : "error",
        profile: null,
        error: message,
      });
    }
  }

  async function updateMe() {
    const token = tokenStore.get("customer");

    if (!token) {
      setState({
        status: "error",
        profile: null,
        error: "Customer token is missing",
      });
      return;
    }

    setIsUpdating(true);

    try {
      const profile = await customerApi.updateMe(token, {
        first_name: null,
        last_name: null,
        phone: null,
        date_of_birth: null,
        marketing_opt_in: false,
      });

      setState({ status: "ready", profile, error: null });
    } catch (err) {
      setState({
        status: "error",
        profile: null,
        error: err instanceof ApiError ? err.message : "Cannot update profile",
      });
    } finally {
      setIsUpdating(false);
    }
  }

  async function loadPurchaseHistory() {
    const token = tokenStore.get("customer");

    if (!token) {
      setHistoryState({
        status: "error",
        purchases: [],
        error: "Customer token is missing",
      });
      return;
    }

    try {
      const purchases = await purchasesApi.list(token);
      setHistoryState({ status: "ready", purchases, error: null });
    } catch (err) {
      setHistoryState({
        status: "error",
        purchases: [],
        error: err instanceof ApiError ? err.message : "Cannot load purchases",
      });
    }
  }

  function logout() {
    tokenStore.clear("customer");
    router.replace("/login");
    router.refresh();
  }

  if (state.status === "loading") {
    return <ProfileMessage message="Loading profile..." />;
  }

  if (state.status === "error") {
    return <ProfileMessage message={state.error} tone="error" />;
  }

  if (state.status === "not-found") {
    return (
      <ProfileShell onLogout={logout}>
        <section className="mt-8 rounded-md border border-black/10 bg-white p-6 text-center shadow-soft">
          <h2 className="text-xl font-semibold text-ink">Profile not found</h2>
          <p className="mt-2 text-sm text-black/55">{state.error}</p>
          <button
            type="button"
            onClick={updateMe}
            disabled={isUpdating}
            className="mt-5 inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isUpdating ? "Updating..." : "UpdateMe"}
          </button>
        </section>
      </ProfileShell>
    );
  }

  const { profile } = state;
  const fullName = [profile.first_name, profile.last_name].filter(Boolean).join(" ");

  return (
    <ProfileShell onLogout={logout}>
      <section className="mt-8 rounded-md border border-black/10 bg-white p-6 shadow-soft">
        <dl className="grid gap-5 sm:grid-cols-2">
          <Info label="User ID" value={String(profile.user_id)} />
          <Info label="Email" value={email ?? "None"} />
          <Info label="Name" value={fullName || "None"} />
          <Info label="Phone" value={profile.phone ?? "None"} />
          <Info label="Date of birth" value={profile.date_of_birth ?? "None"} />
          <Info
            label="Marketing"
            value={profile.marketing_opt_in ? "Opted in" : "Opted out"}
          />
        </dl>
      </section>
      <PurchaseHistory state={historyState} />
    </ProfileShell>
  );
}

function ProfileShell({
  children,
  onLogout,
}: {
  children: React.ReactNode;
  onLogout: () => void;
}) {
  return (
    <main className="mx-auto min-h-screen w-full max-w-4xl px-5 py-8">
      <header className="flex flex-wrap items-center justify-between gap-4 border-b border-black/10 pb-5">
        <div>
          <Link href="/" className="text-sm font-medium text-moss">
            Back to shop
          </Link>
          <h1 className="mt-2 text-3xl font-semibold text-ink">User profile</h1>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Link
            href="/cart"
            className="inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black"
          >
            Cart
          </Link>
          <button
            type="button"
            onClick={onLogout}
            className="inline-flex h-10 items-center rounded-md border border-clay/30 bg-white px-4 text-sm font-semibold text-clay hover:border-clay"
          >
            Logout
          </button>
        </div>
      </header>
      {children}
    </main>
  );
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wide text-black/45">
        {label}
      </dt>
      <dd className="mt-1 font-medium text-ink">{value}</dd>
    </div>
  );
}

function ProfileMessage({
  message,
  tone = "muted",
}: {
  message: string;
  tone?: "muted" | "error";
}) {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-4xl items-center justify-center px-5">
      <p className={tone === "error" ? "text-clay" : "text-black/55"}>{message}</p>
    </main>
  );
}

function PurchaseHistory({ state }: { state: PurchaseHistoryState }) {
  if (state.status === "loading") {
    return (
      <section className="mt-6 rounded-md border border-black/10 bg-white p-6 shadow-soft">
        <h2 className="text-xl font-semibold text-ink">Items to receive</h2>
        <p className="mt-3 text-sm text-black/55">Loading orders...</p>
      </section>
    );
  }

  if (state.status === "error") {
    return (
      <section className="mt-6 rounded-md border border-black/10 bg-white p-6 shadow-soft">
        <h2 className="text-xl font-semibold text-ink">Items to receive</h2>
        <p className="mt-3 text-sm text-clay">{state.error}</p>
      </section>
    );
  }

  if (state.purchases.length === 0) {
    return (
      <section className="mt-6 rounded-md border border-black/10 bg-white p-6 text-center shadow-soft">
        <h2 className="text-xl font-semibold text-ink">Items to receive</h2>
        <p className="mt-2 text-sm text-black/55">No orders yet.</p>
        <Link
          href="/"
          className="mt-5 inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black"
        >
          Shop now
        </Link>
      </section>
    );
  }

  return (
    <section className="mt-6 rounded-md border border-black/10 bg-white p-6 shadow-soft">
      <h2 className="text-xl font-semibold text-ink">Items to receive</h2>
      <div className="mt-5 space-y-5">
        {state.purchases.map((purchase) => (
          <article key={purchase.id} className="border-t border-black/10 pt-5 first:border-t-0 first:pt-0">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <h3 className="font-semibold text-ink">Order #{purchase.id}</h3>
                <p className="mt-1 text-sm text-black/55">
                  {purchase.created_at ? formatDate(purchase.created_at) : "Order date not available"}
                </p>
              </div>
              <p className="font-semibold text-clay">{formatMoney(purchase.total_amount)}</p>
            </div>

            <div className="mt-4 space-y-3">
              {purchase.items.map((item) => {
                const productHref = `/products/${item.clothes_id}`;

                return (
                  <div
                    key={item.id}
                    className="grid gap-3 rounded-md border border-black/10 p-3 sm:grid-cols-[72px_1fr_auto]"
                  >
                    <Link
                      href={productHref}
                      className="block h-[72px] w-[72px] overflow-hidden rounded-md bg-stone-200"
                    >
                      {item.product_image_url ? (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          src={item.product_image_url}
                          alt={item.product_name || `Product #${item.clothes_id}`}
                          className="h-full w-full object-cover"
                        />
                      ) : null}
                    </Link>
                    <div className="min-w-0">
                      <Link href={productHref} className="font-semibold text-ink">
                        {item.product_name || `Product #${item.clothes_id}`}
                      </Link>
                      <p className="mt-1 text-sm text-black/55">
                        Size: {item.size || "No size"} / Qty: {item.quantity}
                      </p>
                      <p className="mt-2 inline-flex rounded-md bg-stone-100 px-2 py-1 text-xs font-medium text-black/60">
                        Waiting to receive
                      </p>
                    </div>
                    <p className="font-semibold text-ink sm:text-right">
                      {formatMoney(item.unit_price * item.quantity)}
                    </p>
                  </div>
                );
              })}
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}

function isProfileNotFound(err: unknown) {
  return (
    err instanceof ApiError &&
    (err.code === 404101 || err.message === "customer profile not found")
  );
}

function formatMoney(value: number) {
  return `${value.toLocaleString("th-TH", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })} THB`;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}
