"use client";

import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import { FormEvent, useState } from "react";
import { ApiError } from "@/lib/api/client";
import { cartApi } from "@/lib/api/cart";
import { tokenStore } from "@/lib/auth/token-store";
import { notifyCartUpdated } from "./cart-count-link";

type AddToCartPanelProps = {
  clothesID: number;
};

const sizeOptions = ["S", "M", "L", "XL"];

export function AddToCartPanel({ clothesID }: AddToCartPanelProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [quantity, setQuantity] = useState(1);
  const [size, setSize] = useState("");
  const [status, setStatus] = useState<"idle" | "adding" | "added">("idle");
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);

    if (!size) {
      setError("Please choose a size before adding this item to your cart.");
      return;
    }

    const token = tokenStore.get("customer");

    if (!token) {
      router.push(`/login?returnTo=${encodeURIComponent(pathname)}`);
      return;
    }

    setStatus("adding");

    try {
      const cart = await cartApi.addItem(token, {
        clothes_id: clothesID,
        quantity,
        size,
      });
      setStatus("added");
      notifyCartUpdated(cart);
    } catch (err) {
      setStatus("idle");
      if (err instanceof ApiError) {
        setError(err.message);
        if (err?.message === "Request failed with status 401") {
          router.push(`/login?returnTo=${encodeURIComponent(pathname)}`);
        }
      } else {
        setError("Cannot add item to cart");
      }
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mt-6 border-t border-black/10 pt-5">
      <div>
        <div className="flex items-center justify-between gap-3">
          <p className="text-sm font-semibold text-ink">Size</p>
          <p className="text-xs font-medium text-black/45">Required</p>
        </div>
        <div className="mt-3 grid grid-cols-4 gap-2" role="radiogroup" aria-label="Size">
          {sizeOptions.map((option) => {
            const isSelected = size === option;

            return (
              <button
                key={option}
                type="button"
                role="radio"
                aria-checked={isSelected}
                onClick={() => {
                  setSize(option);
                  setStatus("idle");
                  setError(null);
                }}
                className={[
                  "h-11 rounded-md border text-sm font-semibold transition",
                  isSelected
                    ? "border-ink bg-ink text-white"
                    : "border-black/15 bg-white text-ink hover:border-moss hover:text-moss",
                ].join(" ")}
              >
                {option}
              </button>
            );
          })}
        </div>
      </div>

      <div className="mt-4">
        <label className="text-sm font-medium text-ink">
          Quantity
          <input
            type="number"
            min="1"
            value={quantity}
            onChange={(event) => {
              setQuantity(Math.max(Number(event.target.value), 1));
              setStatus("idle");
            }}
            className="mt-2 h-11 w-full rounded-md border border-black/15 bg-white px-3 outline-none focus:border-moss"
          />
        </label>
      </div>

      <button
        type="submit"
        disabled={status === "adding"}
        className="mt-4 h-12 w-full rounded-md bg-ink px-4 text-sm font-semibold text-white shadow-soft transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
      >
        {status === "adding" ? "Adding..." : status === "added" ? "Added to cart" : "Add to cart"}
      </button>

      {error ? (
        <p className="mt-3 rounded-md bg-clay/10 px-3 py-2 text-sm font-semibold text-clay">
          {error}
        </p>
      ) : null}

      {status === "added" ? (
        <Link
          href="/cart"
          className="mt-3 inline-flex h-11 w-full items-center justify-center rounded-md border border-black/15 bg-white px-4 text-sm font-semibold text-ink transition hover:border-moss hover:text-moss"
        >
          View cart
        </Link>
      ) : null}
    </form>
  );
}
