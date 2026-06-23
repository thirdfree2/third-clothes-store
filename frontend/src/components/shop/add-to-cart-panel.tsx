"use client";

import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import { FormEvent, useState } from "react";
import { ApiError } from "@/lib/api/client";
import { cartApi } from "@/lib/api/cart";
import { tokenStore } from "@/lib/auth/token-store";

type AddToCartPanelProps = {
  clothesID: number;
};

const sizeOptions = ["", "S", "M", "L", "XL"];

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

    const token = tokenStore.get("customer");

    if (!token) {
      router.push(`/login?returnTo=${encodeURIComponent(pathname)}`);
      return;
    }

    setStatus("adding");

    try {
      await cartApi.addItem(token, {
        clothes_id: clothesID,
        quantity,
        size,
      });
      setStatus("added");
    } catch (err) {
      setStatus("idle");
      if (err instanceof ApiError) {
        setError(err.message);
        if (err?.message == 'Request failed with status 401') {
          router.push(`/login?returnTo=${encodeURIComponent(pathname)}`);
        }
      } else {
        setError("Cannot add item to cart");
      }
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mt-6 border-t border-black/10 pt-5">
      <div className="grid gap-3 sm:grid-cols-2">
        <label className="text-sm font-medium text-ink">
          Size
          <select
            value={size}
            onChange={(event) => {
              setSize(event.target.value);
              setStatus("idle");
            }}
            className="mt-2 h-11 w-full rounded-md border border-black/15 bg-white px-3 outline-none focus:border-moss"
          >
            {sizeOptions.map((option) => (
              <option key={option || "none"} value={option}>
                {option || "No size"}
              </option>
            ))}
          </select>
        </label>

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
        className="mt-4 h-11 w-full rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
      >
        {status === "adding" ? "Adding..." : status === "added" ? "Added to cart" : "Add to cart"}
      </button>

      {error ? <p className="mt-3 text-sm font-medium text-clay">{error}</p> : null}

      {status === "added" ? (
        <Link
          href="/cart"
          className="mt-3 inline-flex h-10 w-full items-center justify-center rounded-md border border-black/15 px-4 text-sm font-semibold text-ink hover:border-black/30"
        >
          View cart
        </Link>
      ) : null}
    </form>
  );
}
