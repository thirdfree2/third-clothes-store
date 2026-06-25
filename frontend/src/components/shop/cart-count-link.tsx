"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { cartApi } from "@/lib/api/cart";
import type { Cart } from "@/lib/api/types";
import { tokenStore } from "@/lib/auth/token-store";

type CartCountLinkProps = {
  children?: string;
  className?: string;
};

const cartUpdatedEvent = "cart:updated";

export function CartCountLink({
  children = "Cart",
  className = "inline-flex h-10 items-center rounded-md border border-black/10 bg-white/90 px-4 text-sm font-semibold text-ink shadow-soft transition hover:border-moss hover:text-moss",
}: CartCountLinkProps) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    const token = tokenStore.get("customer");

    if (token) {
      cartApi
        .get(token)
        .then((cart) => setCount(getCartCount(cart)))
        .catch(() => setCount(0));
    }

    function handleCartUpdated(event: Event) {
      const detail = (event as CustomEvent<{ cart?: Cart; count?: number }>).detail;

      if (typeof detail?.count === "number") {
        setCount(detail.count);
        return;
      }

      if (detail?.cart) {
        setCount(getCartCount(detail.cart));
      }
    }

    window.addEventListener(cartUpdatedEvent, handleCartUpdated);
    return () => window.removeEventListener(cartUpdatedEvent, handleCartUpdated);
  }, []);

  return (
    <Link href="/cart" className={`${className} relative gap-2`}>
      <span>{children}</span>
      {count > 0 ? (
        <span
          aria-label={`${count} items in cart`}
          className="grid h-5 min-w-5 place-items-center rounded-full bg-clay px-1.5 text-[11px] font-bold leading-none text-white"
        >
          {count > 99 ? "99+" : count}
        </span>
      ) : null}
    </Link>
  );
}

export function notifyCartUpdated(cart: Cart) {
  window.dispatchEvent(
    new CustomEvent(cartUpdatedEvent, {
      detail: {
        cart,
        count: getCartCount(cart),
      },
    }),
  );
}

function getCartCount(cart: Cart): number {
  return cart.items.reduce((total, item) => total + item.quantity, 0);
}
