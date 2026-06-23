"use client";

import Link from "next/link";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { categoriesApi } from "@/lib/api/category";
import { colorApi } from "@/lib/api/color";
import { ApiError } from "@/lib/api/client";
import { clothesApi } from "@/lib/api/clothes";
import type { Category, Color } from "@/lib/api/types";
import { decodeAccessToken, hasAllPermissions } from "@/lib/auth/jwt";
import { tokenStore } from "@/lib/auth/token-store";

type ColorState =
  | { status: "idle"; items: Color[]; error: null }
  | { status: "loading"; items: Color[]; error: null }
  | { status: "ready"; items: Color[]; error: null }
  | { status: "error"; items: Color[]; error: string };

type CategoryState =
  | { status: "idle"; items: Category[]; error: null }
  | { status: "loading"; items: Category[]; error: null }
  | { status: "ready"; items: Category[]; error: null }
  | { status: "error"; items: Category[]; error: string };

export function ClothesCreateForm() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [colorID, setColorID] = useState("");
  const [categoryIDs, setCategoryIDs] = useState<number[]>([]);
  const [colorState, setColorState] = useState<ColorState>({
    status: "idle",
    items: [],
    error: null,
  });
  const [categoryState, setCategoryState] = useState<CategoryState>({
    status: "idle",
    items: [],
    error: null,
  });
  const [error, setError] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  const canCreate = useMemo(() => {
    const token = tokenStore.get("admin");
    return hasAllPermissions(token ? decodeAccessToken(token) : null, ["clothes:write"]);
  }, []);

  useEffect(() => {
    if (!canCreate) {
      return;
    }

    let isMounted = true;
    const token = tokenStore.get("admin");

    if (!token) {
      setColorState({
        status: "error",
        items: [],
        error: "Admin token is missing",
      });
      setCategoryState({
        status: "error",
        items: [],
        error: "Admin token is missing",
      });
      return;
    }

    setColorState({ status: "loading", items: [], error: null });
    setCategoryState({ status: "loading", items: [], error: null });

    colorApi
      .list(token)
      .then((items) => {
        if (isMounted) {
          setColorState({ status: "ready", items, error: null });
        }
      })
      .catch((err: unknown) => {
        if (isMounted) {
          setColorState({
            status: "error",
            items: [],
            error: err instanceof Error ? err.message : "Cannot load colors",
          });
        }
      });

    categoriesApi
      .list(token)
      .then((items) => {
        if (isMounted) {
          setCategoryState({ status: "ready", items, error: null });
        }
      })
      .catch((err: unknown) => {
        if (isMounted) {
          setCategoryState({
            status: "error",
            items: [],
            error: err instanceof Error ? err.message : "Cannot load categories",
          });
        }
      });

    return () => {
      isMounted = false;
    };
  }, [canCreate]);

  function toggleCategory(categoryID: number) {
    setCategoryIDs((current) =>
      current.includes(categoryID)
        ? current.filter((id) => id !== categoryID)
        : [...current, categoryID],
    );
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);

    const token = tokenStore.get("admin");
    const parsedPrice = Number(price);
    const parsedColorID = colorID ? Number(colorID) : null;

    if (!token) {
      setError("Admin token is missing");
      return;
    }

    if (!name.trim() || Number.isNaN(parsedPrice) || parsedPrice < 0) {
      setError("Name and valid price are required");
      return;
    }

    setIsSaving(true);

    try {
      const item = await clothesApi.create(token, {
        name: name.trim(),
        price: parsedPrice,
        color_id: parsedColorID,
        category_ids: categoryIDs,
      });
      router.push(`/backoffice/clothes/${item.id}`);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Cannot create clothes");
      }
    } finally {
      setIsSaving(false);
    }
  }

  if (!canCreate) {
    return <CreateMessage message="You do not have permission to create clothes" tone="error" />;
  }

  return (
    <section className="mx-auto w-full max-w-5xl px-5 py-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <Link href="/backoffice" className="text-sm font-medium text-moss">
            Back to clothes
          </Link>
          <h1 className="mt-2 text-3xl font-semibold">Create clothes</h1>
        </div>
      </div>

      <form
        onSubmit={handleSubmit}
        className="mt-6 rounded-md border border-black/10 bg-white p-5 shadow-soft"
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <label className="text-sm font-medium text-ink">
            Name
            <input
              value={name}
              onChange={(event) => setName(event.target.value)}
              className="mt-2 h-11 w-full rounded-md border border-black/15 px-3 outline-none focus:border-moss"
              required
            />
          </label>

          <label className="text-sm font-medium text-ink">
            Price
            <input
              value={price}
              onChange={(event) => setPrice(event.target.value)}
              type="number"
              min="0"
              step="0.01"
              className="mt-2 h-11 w-full rounded-md border border-black/15 px-3 outline-none focus:border-moss"
              required
            />
          </label>

          <SelectInput
            label="Color"
            value={colorID}
            onChange={setColorID}
            disabled={colorState.status === "loading" || colorState.status === "error"}
            options={colorState.items.map((color) => ({
              value: String(color.id),
              label: color.name,
              swatch: color.hex_code,
            }))}
            placeholder={
              colorState.status === "loading"
                ? "Loading colors..."
                : colorState.status === "error"
                  ? "Cannot load colors"
                  : "No color"
            }
            helpText={colorState.error}
          />

          <MultiSelectInput
            label="Categories"
            values={categoryIDs}
            options={categoryState.items.map((category) => ({
              value: category.id,
              label: category.name,
            }))}
            onToggle={toggleCategory}
            disabled={categoryState.status === "loading" || categoryState.status === "error"}
            emptyText={
              categoryState.status === "loading"
                ? "Loading categories..."
                : categoryState.status === "error"
                  ? "Cannot load categories"
                  : "No categories available"
            }
            helpText={categoryState.error}
          />
        </div>

        {error ? <p className="mt-4 text-sm text-clay">{error}</p> : null}

        <div className="mt-6 flex justify-end gap-2">
          <Link
            href="/backoffice"
            className="inline-flex h-10 items-center rounded-md border border-black/15 px-4 text-sm font-medium"
          >
            Cancel
          </Link>
          <button
            type="submit"
            disabled={isSaving}
            className="h-10 rounded-md bg-ink px-4 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isSaving ? "Creating..." : "Create clothes"}
          </button>
        </div>
      </form>
    </section>
  );
}

function MultiSelectInput({
  label,
  values,
  options,
  onToggle,
  disabled = false,
  emptyText = "No options available",
  helpText,
}: {
  label: string;
  values: number[];
  options: { value: number; label: string }[];
  onToggle: (value: number) => void;
  disabled?: boolean;
  emptyText?: string;
  helpText?: string | null;
}) {
  return (
    <fieldset className="text-sm font-medium text-ink">
      <legend>{label}</legend>
      <div className="mt-2 min-h-11 rounded-md border border-black/15 bg-white p-2">
        {options.length > 0 ? (
          <div className="flex flex-wrap gap-2">
            {options.map((option) => {
              const isSelected = values.includes(option.value);

              return (
                <label
                  key={option.value}
                  className={`inline-flex h-8 cursor-pointer items-center rounded-md border px-3 text-xs font-medium ${
                    isSelected
                      ? "border-moss bg-moss text-white"
                      : "border-black/10 bg-stone-100 text-black/65"
                  } ${disabled ? "cursor-not-allowed opacity-50" : ""}`}
                >
                  <input
                    type="checkbox"
                    checked={isSelected}
                    onChange={() => onToggle(option.value)}
                    disabled={disabled}
                    className="sr-only"
                  />
                  {option.label}
                </label>
              );
            })}
          </div>
        ) : (
          <p className="px-1 py-1.5 text-sm font-normal text-black/45">{emptyText}</p>
        )}
      </div>
      {helpText ? <p className="mt-2 text-xs text-clay">{helpText}</p> : null}
    </fieldset>
  );
}

function SelectInput({
  label,
  value,
  options,
  onChange,
  disabled = false,
  placeholder = "Select",
  helpText,
}: {
  label: string;
  value: string;
  options: { value: string; label: string; swatch?: string }[];
  onChange: (value: string) => void;
  disabled?: boolean;
  placeholder?: string;
  helpText?: string | null;
}) {
  const selectedOption = options.find((option) => option.value === value);

  return (
    <label className="text-sm font-medium text-ink">
      {label}
      <div className="relative mt-2">
        {selectedOption?.swatch ? (
          <span
            className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 rounded-full border border-black/15"
            style={{ backgroundColor: selectedOption.swatch }}
          />
        ) : null}
        <select
          value={value}
          onChange={(event) => onChange(event.target.value)}
          disabled={disabled}
          className={`h-11 w-full rounded-md border border-black/15 bg-white px-3 outline-none focus:border-moss disabled:cursor-not-allowed disabled:bg-stone-100 disabled:text-black/45 ${
            selectedOption?.swatch ? "pl-9" : ""
          }`}
        >
          <option value="">{placeholder}</option>
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>
      {helpText ? <p className="mt-2 text-xs text-clay">{helpText}</p> : null}
    </label>
  );
}

function CreateMessage({
  message,
  tone = "muted",
}: {
  message: string;
  tone?: "muted" | "error";
}) {
  return (
    <main className="mx-auto flex min-h-[50vh] w-full max-w-4xl items-center justify-center px-5">
      <p className={tone === "error" ? "text-clay" : "text-black/55"}>{message}</p>
    </main>
  );
}
