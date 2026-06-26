"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { ChangeEvent, FormEvent, useEffect, useMemo, useState } from "react";
import { colorApi } from "@/lib/api/color";
import { clothesApi } from "@/lib/api/clothes";
import { ApiError } from "@/lib/api/client";
import type { Category, Clothes, ClothesImage, Color } from "@/lib/api/types";
import { decodeAccessToken, hasAllPermissions } from "@/lib/auth/jwt";
import { tokenStore } from "@/lib/auth/token-store";
import { categoriesApi } from "@/lib/api/category";

type DetailState =
  | { status: "loading"; item: null; error: null }
  | { status: "ready"; item: Clothes; error: null }
  | { status: "error"; item: null; error: string };

type ClothesDetailProps = {
  id: number;
};

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

export function ClothesDetail({ id }: ClothesDetailProps) {
  const router = useRouter();
  const [state, setState] = useState<DetailState>({
    status: "loading",
    item: null,
    error: null,
  });
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
  const [isEditing, setIsEditing] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const canEdit = useMemo(() => {
    const token = tokenStore.get("admin");
    return hasAllPermissions(token ? decodeAccessToken(token) : null, ["clothes:write"]);
  }, []);
  const canManageImages = useMemo(() => {
    const token = tokenStore.get("admin");
    return hasAllPermissions(token ? decodeAccessToken(token) : null, ["images:write"]);
  }, []);

  async function handleDelete(item: Clothes) {
    setDeleteError(null);

    if (!window.confirm(`Delete ${item.name}?`)) {
      return;
    }

    const token = tokenStore.get("admin");

    if (!token) {
      setDeleteError("Admin token is missing");
      return;
    }

    setIsDeleting(true);

    try {
      await clothesApi.delete(item.id, token);
      router.push("/backoffice");
      router.refresh();
    } catch (err) {
      if (err instanceof ApiError) {
        setDeleteError(err.message);
      } else {
        setDeleteError("Cannot delete clothes");
      }
    } finally {
      setIsDeleting(false);
    }
  }

  useEffect(() => {
    let isMounted = true;

    clothesApi
      .getById(id)
      .then((item) => {
        if (!isMounted) {
          return;
        }

        setState({ status: "ready", item, error: null });
      })
      .catch((error: unknown) => {
        if (!isMounted) {
          return;
        }

        setState({
          status: "error",
          item: null,
          error: error instanceof Error ? error.message : "Cannot load clothes",
        });
      });

    return () => {
      isMounted = false;
    };
  }, [id]);

  useEffect(() => {
    if (!canEdit) {
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
      return;
    }

    setColorState({ status: "loading", items: [], error: null });

    colorApi
      .list(token)
      .then((items) => {
        if (!isMounted) {
          return;
        }

        setColorState({ status: "ready", items, error: null });
      })
      .catch((error: unknown) => {
        if (!isMounted) {
          return;
        }

        setColorState({
          status: "error",
          items: [],
          error: error instanceof Error ? error.message : "Cannot load colors",
        });
      });

    return () => {
      isMounted = false;
    };
  }, [canEdit]);

  useEffect(() => {
    if (!canEdit) {
      return;
    }

    let isMounted = true;
    const token = tokenStore.get("admin");

    if (!token) {
      setCategoryState({
        status: "error",
        items: [],
        error: "Admin token is missing",
      });
      return;
    }

    setCategoryState({ status: "loading", items: [], error: null });

    categoriesApi
      .list(token)
      .then((items) => {
        if (!isMounted) {
          return;
        }

        setCategoryState({ status: "ready", items, error: null });
      })
      .catch((error: unknown) => {
        if (!isMounted) {
          return;
        }

        setCategoryState({
          status: "error",
          items: [],
          error: error instanceof Error ? error.message : "Cannot load categories",
        });
      });

    return () => {
      isMounted = false;
    };
  }, [canEdit]);

  if (state.status === "loading") {
    return <DetailMessage message="Loading clothes..." />;
  }

  if (state.status === "error") {
    return <DetailMessage message={state.error} tone="error" />;
  }

  return (
    <section className="mx-auto w-full max-w-5xl px-5 py-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <Link href="/backoffice" className="text-sm font-medium text-moss">
            Back to clothes
          </Link>
          <h1 className="mt-2 text-3xl font-semibold">{state.item.name}</h1>
        </div>

        {canEdit ? (
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => setIsEditing((current) => !current)}
              className="h-10 rounded-md border border-black/15 bg-white px-4 text-sm font-medium shadow-soft hover:border-black/30"
            >
              {isEditing ? "Cancel edit" : "Edit"}
            </button>
            <button
              type="button"
              onClick={() => handleDelete(state.item)}
              disabled={isDeleting}
              className="h-10 rounded-md border border-clay/30 bg-white px-4 text-sm font-medium text-clay shadow-soft hover:border-clay disabled:cursor-not-allowed disabled:opacity-60"
            >
              {isDeleting ? "Deleting..." : "Delete"}
            </button>
          </div>
        ) : null}
      </div>

      {deleteError ? <p className="mt-4 text-sm font-medium text-clay">{deleteError}</p> : null}

      {isEditing ? (
        <EditForm
          item={state.item}
          colors={colorState.items}
          colorStatus={colorState.status}
          colorError={colorState.error}
          categories={categoryState.items}
          categoryStatus={categoryState.status}
          categoryError={categoryState.error}
          onCancel={() => setIsEditing(false)}
          onSaved={(item) => {
            setState({ status: "ready", item, error: null });
            setIsEditing(false);
          }}
        />
      ) : (
        <ReadView
          item={state.item}
          canManageImages={canManageImages}
          onImagesChanged={(images) => {
            setState((current) => {
              if (current.status !== "ready") {
                return current;
              }

              return {
                status: "ready",
                item: { ...current.item, images },
                error: null,
              };
            });
          }}
        />
      )}
    </section>
  );
}

function ReadView({
  item,
  canManageImages,
  onImagesChanged,
}: {
  item: Clothes;
  canManageImages: boolean;
  onImagesChanged: (images: ClothesImage[]) => void;
}) {
  const [viewerImageID, setViewerImageID] = useState<number | null>(null);
  const primaryImage = item.images[0] ?? null;
  const viewerImage = item.images.find((image) => image.id === viewerImageID) ?? null;

  useEffect(() => {
    if (item.images.length === 0) {
      setViewerImageID(null);
      return;
    }

    if (viewerImageID && !item.images.some((image) => image.id === viewerImageID)) {
      setViewerImageID(null);
    }
  }, [item.images, viewerImageID]);

  useEffect(() => {
    if (viewerImageID === null) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setViewerImageID(null);
      }

      if (event.key === "ArrowLeft") {
        showAdjacentViewerImage("previous");
      }

      if (event.key === "ArrowRight") {
        showAdjacentViewerImage("next");
      }
    }

    document.body.style.overflow = "hidden";
    window.addEventListener("keydown", handleKeyDown);

    return () => {
      document.body.style.overflow = "";
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [item.images, viewerImageID]);

  function openViewer(index: number) {
    const image = item.images[index];

    if (!image?.image_url) {
      return;
    }

    setViewerImageID(image.id);
  }

  function showAdjacentViewerImage(direction: "previous" | "next") {
    setViewerImageID((current) => {
      if (current === null || item.images.length === 0) {
        return current;
      }

      const currentIndex = item.images.findIndex((image) => image.id === current);
      if (currentIndex === -1) {
        return item.images[0]?.id ?? null;
      }

      const nextIndex =
        direction === "previous"
          ? currentIndex === 0
            ? item.images.length - 1
            : currentIndex - 1
          : currentIndex === item.images.length - 1
            ? 0
            : currentIndex + 1;

      return item.images[nextIndex]?.id ?? null;
    });
  }

  return (
    <div className="mt-6 grid gap-6 lg:grid-cols-[320px_1fr]">
      <div className="space-y-4">
        <div className="overflow-hidden rounded-md border border-black/10 bg-white shadow-soft">
          <button
            type="button"
            onClick={() => openViewer(0)}
            disabled={!primaryImage?.image_url}
            className="block w-full disabled:cursor-default"
            aria-label={primaryImage?.image_url ? `View ${item.name} image larger` : undefined}
          >
            <div className="aspect-square bg-stone-200">
              {primaryImage?.image_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={primaryImage.image_url}
                  alt={item.name}
                  className="h-full w-full object-cover transition duration-300 hover:scale-[1.02]"
                />
              ) : null}
            </div>
          </button>
        </div>

        {canManageImages ? (
          <ImageManager item={item} onImagesChanged={onImagesChanged} />
        ) : null}
      </div>

      <div className="rounded-md border border-black/10 bg-white p-5 shadow-soft">
        <dl className="grid gap-5 sm:grid-cols-2">
          <Info label="ID" value={String(item.id)} />
          <Info label="Price" value={formatMoney(item.price)} />
          <Info label="Color" value={item.color?.name ?? "None"} />
          <Info
            label="Categories"
            value={
              item.categories.length > 0
                ? item.categories.map((category) => category.name).join(", ")
                : "None"
            }
          />
          <Info label="Images" value={String(item.images.length)} />
          <Info label="Updated" value={formatDate(item.updated_at)} />
        </dl>

        <div className="mt-6 border-t border-black/10 pt-5">
          <h2 className="text-sm font-semibold text-ink">Images</h2>
          <div className="mt-3 grid grid-cols-2 gap-3 sm:grid-cols-3">
            {item.images.length > 0 ? (
              item.images.map((image, index) => (
                <button
                  key={image.id}
                  type="button"
                  onClick={() => openViewer(index)}
                  className="overflow-hidden rounded-md border border-black/10 bg-stone-100 transition hover:border-moss"
                  aria-label={`View ${item.name} image ${index + 1} larger`}
                >
                  <div className="aspect-square">
                    {image.image_url ? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img
                        src={image.image_url}
                        alt={`${item.name} image ${image.id}`}
                        className="h-full w-full object-cover"
                      />
                    ) : null}
                  </div>
                </button>
              ))
            ) : (
              <p className="col-span-full text-sm text-black/45">No images</p>
            )}
          </div>
        </div>
      </div>

      {viewerImage?.image_url ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/82 p-4 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-label={`${item.name} large image viewer`}
          onClick={() => setViewerImageID(null)}
        >
          <button
            type="button"
            onClick={(event) => {
              event.stopPropagation();
              setViewerImageID(null);
            }}
            className="absolute right-4 top-4 grid h-11 w-11 place-items-center rounded-md bg-white text-2xl font-semibold text-ink shadow-soft transition hover:bg-stone-100"
            aria-label="Close image viewer"
          >
            &times;
          </button>

          {item.images.length > 1 ? (
            <button
              type="button"
              onClick={(event) => {
                event.stopPropagation();
                showAdjacentViewerImage("previous");
              }}
              className="absolute left-4 top-1/2 grid h-11 w-11 -translate-y-1/2 place-items-center rounded-md bg-white text-2xl font-semibold text-ink shadow-soft transition hover:bg-stone-100"
              aria-label="Previous image"
            >
              &larr;
            </button>
          ) : null}

          <div
            className="max-h-[88vh] w-full max-w-5xl"
            onClick={(event) => event.stopPropagation()}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={viewerImage.image_url}
              alt={`${item.name} large view`}
              className="mx-auto max-h-[88vh] w-auto max-w-full rounded-md object-contain shadow-soft"
            />
          </div>

          {item.images.length > 1 ? (
            <button
              type="button"
              onClick={(event) => {
                event.stopPropagation();
                showAdjacentViewerImage("next");
              }}
              className="absolute right-4 top-1/2 grid h-11 w-11 -translate-y-1/2 place-items-center rounded-md bg-white text-2xl font-semibold text-ink shadow-soft transition hover:bg-stone-100"
              aria-label="Next image"
            >
              &rarr;
            </button>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}

function ImageManager({
  item,
  onImagesChanged,
}: {
  item: Clothes;
  onImagesChanged: (images: ClothesImage[]) => void;
}) {
  const [files, setFiles] = useState<File[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [deletingImageID, setDeletingImageID] = useState<number | null>(null);

  function handleFilesChange(event: ChangeEvent<HTMLInputElement>) {
    setFiles(Array.from(event.target.files ?? []));
    setError(null);
  }

  async function handleUpload(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);

    const token = tokenStore.get("admin");

    if (!token) {
      setError("Admin token is missing");
      return;
    }

    if (files.length === 0) {
      setError("Choose at least one image");
      return;
    }

    setIsUploading(true);

    try {
      const uploadedImages = await clothesApi.uploadImages(item.id, token, files);
      onImagesChanged([...item.images, ...uploadedImages]);
      setFiles([]);
      event.currentTarget.reset();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Cannot upload images");
      }
    } finally {
      setIsUploading(false);
    }
  }

  async function handleDelete(imageID: number) {
    setError(null);

    const token = tokenStore.get("admin");

    if (!token) {
      setError("Admin token is missing");
      return;
    }

    setDeletingImageID(imageID);

    try {
      await clothesApi.deleteImage(item.id, imageID, token);
      onImagesChanged(item.images.filter((image) => image.id !== imageID));
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Cannot delete image");
      }
    } finally {
      setDeletingImageID(null);
    }
  }

  return (
    <div className="rounded-md border border-black/10 bg-white p-4 shadow-soft">
      <form onSubmit={handleUpload}>
        <label className="text-sm font-medium text-ink">
          Upload images
          <input
            type="file"
            accept="image/*"
            multiple
            onChange={handleFilesChange}
            className="mt-2 block w-full text-sm text-black/65 file:mr-3 file:h-9 file:rounded-md file:border-0 file:bg-stone-100 file:px-3 file:text-sm file:font-medium file:text-ink"
          />
        </label>

        <button
          type="submit"
          disabled={isUploading}
          className="mt-3 h-10 w-full rounded-md bg-ink px-4 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isUploading ? "Uploading..." : `Upload${files.length > 0 ? ` ${files.length}` : ""}`}
        </button>
      </form>

      {item.images.length > 0 ? (
        <div className="mt-4 grid grid-cols-2 gap-2">
          {item.images.map((image) => (
            <div
              key={image.id}
              className="overflow-hidden rounded-md border border-black/10 bg-stone-100"
            >
              <div className="aspect-square">
                {image.image_url ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={image.image_url}
                    alt={`${item.name} image ${image.id}`}
                    className="h-full w-full object-cover"
                  />
                ) : null}
              </div>
              <button
                type="button"
                onClick={() => handleDelete(image.id)}
                disabled={deletingImageID === image.id}
                className="h-9 w-full border-t border-black/10 bg-white text-xs font-medium text-clay disabled:cursor-not-allowed disabled:opacity-60"
              >
                {deletingImageID === image.id ? "Deleting..." : "Delete"}
              </button>
            </div>
          ))}
        </div>
      ) : null}

      {error ? <p className="mt-3 text-sm text-clay">{error}</p> : null}
    </div>
  );
}

function EditForm({
  item,
  colors,
  colorStatus,
  colorError,
  categories,
  categoryStatus,
  categoryError,
  onCancel,
  onSaved,
}: {
  item: Clothes;
  colors: Color[];
  colorStatus: ColorState["status"];
  colorError: string | null;
  categories: Category[];
  categoryStatus: CategoryState["status"];
  categoryError: string | null;
  onCancel: () => void;
  onSaved: (item: Clothes) => void;
}) {
  const [name, setName] = useState(item.name);
  const [price, setPrice] = useState(String(item.price));
  const [colorID, setColorID] = useState(item.color_id ? String(item.color_id) : "");
  const [categoryIDs, setCategoryIDs] = useState<number[]>(
    item.categories.map((category) => category.id),
  );
  const [error, setError] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

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
      const updatedItem = await clothesApi.update(item.id, token, {
        name: name.trim(),
        price: parsedPrice,
        color_id: parsedColorID,
        category_ids: categoryIDs,
      });
      onSaved(updatedItem);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Cannot save clothes");
      }
    } finally {
      setIsSaving(false);
    }
  }

  return (
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
          disabled={colorStatus === "loading" || colorStatus === "error"}
          options={colors.map((color) => ({
            value: String(color.id),
            label: color.name,
            swatch: color.hex_code,
          }))}
          placeholder={
            colorStatus === "loading"
              ? "Loading colors..."
              : colorStatus === "error"
                ? "Cannot load colors"
                : "No color"
          }
          helpText={colorError}
        />

        <MultiSelectInput
          label="Categories"
          values={categoryIDs}
          options={categories.map((category) => ({
            value: category.id,
            label: category.name,
          }))}
          onToggle={toggleCategory}
          disabled={categoryStatus === "loading" || categoryStatus === "error"}
          emptyText={
            categoryStatus === "loading"
              ? "Loading categories..."
              : categoryStatus === "error"
                ? "Cannot load categories"
                : "No categories available"
          }
          helpText={categoryError}
        />
      </div>

      {error ? <p className="mt-4 text-sm text-clay">{error}</p> : null}

      <div className="mt-6 flex justify-end gap-2">
        <button
          type="button"
          onClick={onCancel}
          className="h-10 rounded-md border border-black/15 px-4 text-sm font-medium"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isSaving}
          className="h-10 rounded-md bg-ink px-4 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isSaving ? "Saving..." : "Save changes"}
        </button>
      </div>
    </form>
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

function DetailMessage({
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
