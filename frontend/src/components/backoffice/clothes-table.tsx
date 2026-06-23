"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { ApiError } from "@/lib/api/client";
import { clothesApi } from "@/lib/api/clothes";
import type { Clothes, PaginationMeta } from "@/lib/api/types";
import { decodeAccessToken, hasAllPermissions } from "@/lib/auth/jwt";
import { tokenStore } from "@/lib/auth/token-store";

type TableState =
  | { status: "loading"; items: Clothes[]; pagination: PaginationMeta | null; error: null }
  | { status: "ready"; items: Clothes[]; pagination: PaginationMeta; error: null }
  | { status: "error"; items: Clothes[]; pagination: PaginationMeta | null; error: string };

const defaultPerPage = 5;

export function ClothesTable() {
  const [state, setState] = useState<TableState>({
    status: "loading",
    items: [],
    pagination: null,
    error: null,
  });
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const [deletingID, setDeletingID] = useState<number | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  const canDelete = useMemo(() => {
    const token = tokenStore.get("admin");
    return hasAllPermissions(token ? decodeAccessToken(token) : null, ["clothes:write"]);
  }, []);

  useEffect(() => {
    let isMounted = true;

    setState((current) => ({
      status: "loading",
      items: current.items,
      pagination: current.pagination,
      error: null,
    }));

    clothesApi
      .list({ page, perPage: defaultPerPage })
      .then((payload) => {
        if (!isMounted) {
          return;
        }

        setState({
          status: "ready",
          items: payload.items,
          pagination: payload.pagination,
          error: null,
        });
      })
      .catch((error: unknown) => {
        if (!isMounted) {
          return;
        }

        setState({
          status: "error",
          items: [],
          pagination: null,
          error: error instanceof Error ? error.message : "Cannot load clothes",
        });
      });

    return () => {
      isMounted = false;
    };
  }, [page]);

  const filteredItems = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase();
    if (!normalizedQuery) {
      return state.items;
    }

    return state.items.filter((item) => {
      const categories = item.categories.map((category) => category.name).join(" ");
      return `${item.name} ${categories}`.toLowerCase().includes(normalizedQuery);
    });
  }, [query, state.items]);

  async function handleDelete(item: Clothes) {
    setActionError(null);

    if (!window.confirm(`Delete ${item.name}?`)) {
      return;
    }

    const token = tokenStore.get("admin");

    if (!token) {
      setActionError("Admin token is missing");
      return;
    }

    setDeletingID(item.id);

    try {
      await clothesApi.delete(item.id, token);
      setState((current) => ({
        ...current,
        items: current.items.filter((currentItem) => currentItem.id !== item.id),
        ...(current.pagination
          ? {
              pagination: {
                ...current.pagination,
                total: Math.max(current.pagination.total - 1, 0),
              },
            }
          : {}),
      }));
    } catch (err) {
      if (err instanceof ApiError) {
        setActionError(err.message);
      } else {
        setActionError("Cannot delete clothes");
      }
    } finally {
      setDeletingID(null);
    }
  }

  return (
    <section className="mx-auto w-full max-w-7xl px-5 py-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm font-medium uppercase tracking-wide text-moss">
            Inventory
          </p>
          <h1 className="mt-2 text-3xl font-semibold">Clothes</h1>
        </div>

        <div className="flex w-full flex-wrap items-end gap-3 sm:w-auto">
          <label className="w-full max-w-xs text-sm font-medium text-black/70">
            Search
            <input
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Name or category"
              className="mt-2 h-10 w-full rounded-md border border-black/15 bg-white px-3 outline-none focus:border-moss"
            />
          </label>
          <Link
            href="/backoffice/clothes/new"
            className="inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-medium text-white hover:bg-black"
          >
            Create clothes
          </Link>
        </div>
      </div>

      {actionError ? <p className="mt-4 text-sm font-medium text-clay">{actionError}</p> : null}

      <div className="mt-6 overflow-hidden rounded-md border border-black/10 bg-white shadow-soft">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[840px] border-collapse text-left text-sm">
            <thead className="bg-stone-100 text-xs uppercase tracking-wide text-black/55">
              <tr>
                <th className="px-4 py-3 font-semibold">Product</th>
                <th className="px-4 py-3 font-semibold">Price</th>
                <th className="px-4 py-3 font-semibold">Color</th>
                <th className="px-4 py-3 font-semibold">Categories</th>
                <th className="px-4 py-3 text-right font-semibold">Images</th>
                <th className="px-4 py-3 text-right font-semibold">ID</th>
                <th className="px-4 py-3 text-right font-semibold">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-black/10">
              {state.status === "loading" ? (
                <TableMessage message="Loading clothes..." />
              ) : null}

              {state.status === "error" ? (
                <TableMessage message={state.error} tone="error" />
              ) : null}

              {state.status === "ready" && filteredItems.length === 0 ? (
                <TableMessage message="No clothes found" />
              ) : null}

              {filteredItems.map((item) => (
                <tr key={item.id} className="hover:bg-stone-50">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <div className="h-12 w-12 overflow-hidden rounded-md bg-stone-200">
                        {item.images[0]?.image_url ? (
                          // eslint-disable-next-line @next/next/no-img-element
                          <img
                            src={item.images[0].image_url}
                            alt={item.name}
                            className="h-full w-full object-cover"
                          />
                        ) : null}
                      </div>
                      <div>
                        <p className="font-semibold text-ink">{item.name}</p>
                        <p className="mt-1 text-xs text-black/50">
                          Updated {formatDate(item.updated_at)}
                        </p>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3 font-medium text-clay">
                    {formatMoney(item.price)}
                  </td>
                  <td className="px-4 py-3">
                    {item.color ? (
                      <span className="inline-flex items-center gap-2">
                        <span
                          className="h-4 w-4 rounded-full border border-black/15"
                          style={{ backgroundColor: item.color.hex_code }}
                        />
                        {item.color.name}
                      </span>
                    ) : (
                      <span className="text-black/40">None</span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex flex-wrap gap-1.5">
                      {item.categories.length > 0 ? (
                        item.categories.map((category) => (
                          <span
                            key={category.id}
                            className="rounded-md bg-stone-100 px-2 py-1 text-xs font-medium text-black/65"
                          >
                            {category.name}
                          </span>
                        ))
                      ) : (
                        <span className="text-black/40">None</span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3 text-right">{item.images.length}</td>
                  <td className="px-4 py-3 text-right text-black/45">{item.id}</td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex justify-end gap-2">
                      <Link
                        href={`/backoffice/clothes/${item.id}`}
                        className="inline-flex h-9 items-center rounded-md border border-black/15 px-3 text-sm font-medium hover:border-black/30"
                      >
                        View
                      </Link>
                      {canDelete ? (
                        <button
                          type="button"
                          onClick={() => handleDelete(item)}
                          disabled={deletingID === item.id}
                          className="inline-flex h-9 items-center rounded-md border border-clay/30 px-3 text-sm font-medium text-clay hover:border-clay disabled:cursor-not-allowed disabled:opacity-60"
                        >
                          {deletingID === item.id ? "Deleting..." : "Delete"}
                        </button>
                      ) : null}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-black/10 px-4 py-3">
          <p className="text-sm text-black/55">
            {state.pagination
              ? `Page ${state.pagination.page} of ${Math.max(state.pagination.total_pages, 1)} - ${state.pagination.total} items`
              : "Page 1"}
          </p>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setPage((current) => Math.max(current - 1, 1))}
              disabled={state.status === "loading" || !state.pagination || state.pagination.page <= 1}
              className="h-9 rounded-md border border-black/15 bg-white px-3 text-sm font-medium disabled:cursor-not-allowed disabled:opacity-45"
            >
              Previous
            </button>
            <button
              type="button"
              onClick={() => setPage((current) => current + 1)}
              disabled={
                state.status === "loading" ||
                !state.pagination ||
                state.pagination.page >= state.pagination.total_pages
              }
              className="h-9 rounded-md border border-black/15 bg-white px-3 text-sm font-medium disabled:cursor-not-allowed disabled:opacity-45"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}

function TableMessage({
  message,
  tone = "muted",
}: {
  message: string;
  tone?: "muted" | "error";
}) {
  return (
    <tr>
      <td
        colSpan={7}
        className={`px-4 py-10 text-center ${
          tone === "error" ? "text-clay" : "text-black/55"
        }`}
      >
        {message}
      </td>
    </tr>
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
