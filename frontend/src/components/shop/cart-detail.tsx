"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { ApiError } from "@/lib/api/client";
import { cartApi } from "@/lib/api/cart";
import { purchasesApi } from "@/lib/api/purchases";
import type { Cart, CartItem } from "@/lib/api/types";
import { tokenStore } from "@/lib/auth/token-store";

type CartState =
  | { status: "loading"; cart: Cart | null; error: null }
  | { status: "ready"; cart: Cart; error: null }
  | { status: "error"; cart: Cart | null; error: string };

export function CartDetail() {
  const [state, setState] = useState<CartState>({
    status: "loading",
    cart: null,
    error: null,
  });
  const [actionError, setActionError] = useState<string | null>(null);
  const [busyItemID, setBusyItemID] = useState<number | null>(null);
  const [isClearing, setIsClearing] = useState(false);
  const [isCheckingOut, setIsCheckingOut] = useState(false);
  const [checkoutMessage, setCheckoutMessage] = useState<string | null>(null);

  async function loadCart() {
    const token = tokenStore.get("customer");

    if (!token) {
      setState({
        status: "error",
        cart: null,
        error: "Customer token is missing",
      });
      return;
    }

    try {
      const cart = await cartApi.get(token);
      setState({ status: "ready", cart, error: null });
    } catch (err) {
      setState({
        status: "error",
        cart: null,
        error: err instanceof ApiError ? err.message : "Cannot load cart",
      });
    }
  }

  useEffect(() => {
    void loadCart();
  }, []);

  const totals = useMemo(() => {
    const items = state.cart?.items ?? [];

    return items.reduce(
      (summary, item) => {
        const lineTotal = item.quantity * (item.clothes?.price ?? 0);
        return {
          quantity: summary.quantity + item.quantity,
          subtotal: summary.subtotal + lineTotal,
        };
      },
      { quantity: 0, subtotal: 0 },
    );
  }, [state.cart]);

  async function updateQuantity(item: CartItem, quantity: number) {
    const nextQuantity = Math.max(quantity, 1);
    const token = tokenStore.get("customer");

    if (!token) {
      setActionError("Customer token is missing");
      return;
    }

    setActionError(null);
    setBusyItemID(item.id);

    try {
      await cartApi.updateItem(token, item.id, { quantity: nextQuantity });
      await loadCart();
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : "Cannot update item");
    } finally {
      setBusyItemID(null);
    }
  }

  async function deleteItem(item: CartItem) {
    const token = tokenStore.get("customer");

    if (!token) {
      setActionError("Customer token is missing");
      return;
    }

    setActionError(null);
    setBusyItemID(item.id);

    try {
      await cartApi.deleteItem(token, item.id);
      await loadCart();
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : "Cannot remove item");
    } finally {
      setBusyItemID(null);
    }
  }

  async function clearCart() {
    const token = tokenStore.get("customer");

    if (!token) {
      setActionError("Customer token is missing");
      return;
    }

    if (!window.confirm("Clear all items from cart?")) {
      return;
    }

    setActionError(null);
    setIsClearing(true);

    try {
      await cartApi.clear(token);
      await loadCart();
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : "Cannot clear cart");
    } finally {
      setIsClearing(false);
    }
  }

  async function checkout() {
    const token = tokenStore.get("customer");

    if (!token || state.status !== "ready" || state.cart.items.length === 0) {
      setActionError("Cannot checkout right now");
      return;
    }

    setActionError(null);
    setCheckoutMessage(null);
    setIsCheckingOut(true);

    try {
      const purchase = await purchasesApi.createFromCart(token, state.cart);
      await cartApi.clear(token);
      await loadCart();
      setCheckoutMessage(`Order #${purchase.id} created successfully`);
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : "Cannot checkout");
    } finally {
      setIsCheckingOut(false);
    }
  }

  if (state.status === "loading") {
    return <CartMessage message="Loading cart..." />;
  }

  if (state.status === "error") {
    return <CartMessage message={state.error} tone="error" />;
  }

  const items = state.cart.items;

  return (
    <main className="mx-auto min-h-screen w-full max-w-6xl px-5 py-8">
      <header className="flex flex-wrap items-center justify-between gap-4 border-b border-black/10 pb-5">
        <div>
          <Link href="/" className="text-sm font-medium text-moss">
            Back to shop
          </Link>
          <h1 className="mt-2 text-3xl font-semibold text-ink">Cart</h1>
        </div>
        <p className="rounded-md bg-white px-4 py-2 text-sm font-medium text-black/60 shadow-soft">
          {totals.quantity} {totals.quantity === 1 ? "item" : "items"}
        </p>
      </header>

      {actionError ? <p className="mt-4 text-sm font-medium text-clay">{actionError}</p> : null}
      {checkoutMessage ? (
        <p className="mt-4 text-sm font-medium text-moss">{checkoutMessage}</p>
      ) : null}

      {items.length === 0 ? (
        <section className="mt-8 rounded-md border border-black/10 bg-white p-8 text-center shadow-soft">
          <h2 className="text-xl font-semibold text-ink">
            {checkoutMessage ? "Checkout complete" : "Your cart is empty"}
          </h2>
          <p className="mt-2 text-sm text-black/55">
            {checkoutMessage
              ? "You can review items to receive in your profile."
              : "Choose a product and add it to cart."}
          </p>
          <div className="mt-5 flex flex-wrap items-center justify-center gap-2">
            <Link
              href="/"
              className="inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black"
            >
              Shop now
            </Link>
            {checkoutMessage ? (
              <Link
                href="/profile"
                className="inline-flex h-10 items-center rounded-md border border-black/10 bg-white px-4 text-sm font-semibold text-ink hover:border-moss hover:text-moss"
              >
                View profile
              </Link>
            ) : null}
          </div>
        </section>
      ) : (
        <section className="mt-8 grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
          <div className="space-y-3">
            {items.map((item) => (
              <CartLineItem
                key={item.id}
                item={item}
                isBusy={busyItemID === item.id}
                onDecrease={() => updateQuantity(item, item.quantity - 1)}
                onIncrease={() => updateQuantity(item, item.quantity + 1)}
                onDelete={() => deleteItem(item)}
              />
            ))}
          </div>

          <aside className="h-fit rounded-md border border-black/10 bg-white p-5 shadow-soft">
            <h2 className="text-lg font-semibold text-ink">Order summary</h2>
            <dl className="mt-5 space-y-3 text-sm">
              <SummaryRow label="Items" value={String(totals.quantity)} />
              <SummaryRow label="Subtotal" value={formatMoney(totals.subtotal)} />
              <SummaryRow label="Shipping" value="Calculated later" />
            </dl>
            <div className="mt-5 border-t border-black/10 pt-4">
              <SummaryRow label="Total" value={formatMoney(totals.subtotal)} strong />
            </div>

            <button
              type="button"
              onClick={checkout}
              disabled={isCheckingOut || items.length === 0}
              className="mt-5 h-11 w-full rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black"
            >
              {isCheckingOut ? "Checking out..." : "Checkout"}
            </button>
            <button
              type="button"
              onClick={clearCart}
              disabled={isClearing}
              className="mt-2 h-10 w-full rounded-md border border-clay/30 px-4 text-sm font-medium text-clay hover:border-clay disabled:cursor-not-allowed disabled:opacity-60"
            >
              {isClearing ? "Clearing..." : "Clear cart"}
            </button>
          </aside>
        </section>
      )}
    </main>
  );
}

function CartLineItem({
  item,
  isBusy,
  onDecrease,
  onIncrease,
  onDelete,
}: {
  item: CartItem;
  isBusy: boolean;
  onDecrease: () => void;
  onIncrease: () => void;
  onDelete: () => void;
}) {
  const product = item.clothes;
  const lineTotal = item.quantity * (product?.price ?? 0);

  return (
    <article className="grid gap-4 rounded-md border border-black/10 bg-white p-4 shadow-soft sm:grid-cols-[96px_1fr_auto]">
      <Link
        href={`/products/${item.clothes_id}`}
        className="block h-24 w-24 overflow-hidden rounded-md bg-stone-200"
      >
        {product?.images[0]?.image_url ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={product.images[0].image_url}
            alt={product.name}
            className="h-full w-full object-cover"
          />
        ) : null}
      </Link>

      <div className="min-w-0">
        <Link href={`/products/${item.clothes_id}`} className="font-semibold text-ink">
          {product?.name ?? `Product #${item.clothes_id}`}
        </Link>
        <p className="mt-1 text-sm text-black/55">
          Size: {item.size || "No size"}
        </p>
        <p className="mt-2 text-sm font-medium text-clay">
          {formatMoney(product?.price ?? 0)}
        </p>

        <div className="mt-4 inline-flex h-9 items-center overflow-hidden rounded-md border border-black/15 bg-white">
          <button
            type="button"
            onClick={onDecrease}
            disabled={isBusy || item.quantity <= 1}
            className="h-full w-9 text-lg font-medium disabled:cursor-not-allowed disabled:opacity-40"
          >
            -
          </button>
          <span className="min-w-10 px-3 text-center text-sm font-semibold">
            {item.quantity}
          </span>
          <button
            type="button"
            onClick={onIncrease}
            disabled={isBusy}
            className="h-full w-9 text-lg font-medium disabled:cursor-not-allowed disabled:opacity-40"
          >
            +
          </button>
        </div>
      </div>

      <div className="flex flex-row items-center justify-between gap-4 sm:flex-col sm:items-end">
        <p className="font-semibold text-ink">{formatMoney(lineTotal)}</p>
        <button
          type="button"
          onClick={onDelete}
          disabled={isBusy}
          className="h-9 rounded-md border border-clay/30 px-3 text-sm font-medium text-clay hover:border-clay disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isBusy ? "Updating..." : "Remove"}
        </button>
      </div>
    </article>
  );
}

function SummaryRow({
  label,
  value,
  strong = false,
}: {
  label: string;
  value: string;
  strong?: boolean;
}) {
  return (
    <div className={`flex items-center justify-between gap-4 ${strong ? "text-base font-semibold" : ""}`}>
      <dt className="text-black/55">{label}</dt>
      <dd className="text-right text-ink">{value}</dd>
    </div>
  );
}

function CartMessage({
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

function formatMoney(value: number) {
  return `${value.toLocaleString("th-TH", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })} THB`;
}
