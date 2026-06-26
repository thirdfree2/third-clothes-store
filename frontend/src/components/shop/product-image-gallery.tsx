"use client";

import { useEffect, useState } from "react";
import type { ClothesImage } from "@/lib/api/types";

type ProductImageGalleryProps = {
  images: ClothesImage[];
  productID: number;
  productName: string;
};

export function ProductImageGallery({
  images,
  productID,
  productName,
}: ProductImageGalleryProps) {
  const [viewerImageID, setViewerImageID] = useState<number | null>(null);
  const [isViewerOpen, setIsViewerOpen] = useState(false);
  const primaryImage = images[0] ?? null;
  const viewerImage = images.find((image) => image.id === viewerImageID) ?? null;
  const secondaryImages = images.slice(1, 5);

  useEffect(() => {
    if (images.length === 0) {
      setViewerImageID(null);
      setIsViewerOpen(false);
      return;
    }

    if (viewerImageID && !images.some((image) => image.id === viewerImageID)) {
      setViewerImageID(null);
      setIsViewerOpen(false);
    }
  }, [images, viewerImageID]);

  useEffect(() => {
    if (!isViewerOpen) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setIsViewerOpen(false);
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
  }, [images, isViewerOpen, viewerImageID]);

  function openViewer(index: number) {
    const image = images[index];

    if (!image?.image_url) {
      return;
    }

    setViewerImageID(image.id);
    setIsViewerOpen(true);
  }

  function showAdjacentViewerImage(direction: "previous" | "next") {
    setViewerImageID((current) => {
      if (current === null || images.length === 0) {
        return current;
      }

      const currentIndex = images.findIndex((image) => image.id === current);
      if (currentIndex === -1) {
        return images[0]?.id ?? null;
      }

      const nextIndex =
        direction === "previous"
          ? currentIndex === 0
            ? images.length - 1
            : currentIndex - 1
          : currentIndex === images.length - 1
            ? 0
            : currentIndex + 1;
      const nextImage = images[nextIndex] ?? null;

      return nextImage?.id ?? null;
    });
  }

  return (
    <div className="space-y-4">
      <div className="relative overflow-hidden rounded-md border border-black/10 bg-white shadow-soft">
        <button
          type="button"
          onClick={() => openViewer(0)}
          disabled={!primaryImage?.image_url}
          className="block w-full text-left disabled:cursor-default"
          aria-label={primaryImage?.image_url ? `View ${productName} image larger` : undefined}
        >
          <div className="aspect-[4/5] bg-[#ded4c9] sm:aspect-[5/4] lg:aspect-[4/5]">
            {primaryImage?.image_url ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={primaryImage.image_url}
                alt={productName}
                className="h-full w-full object-cover transition duration-300 hover:scale-[1.02]"
              />
            ) : (
              <div className="flex h-full items-center justify-center p-8 text-center text-sm font-semibold uppercase tracking-[0.2em] text-black/35">
                Clothes Store
              </div>
            )}
          </div>
        </button>
        <div className="absolute left-4 top-4 rounded-md bg-white/92 px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.16em] text-moss shadow-soft backdrop-blur">
          Product #{productID}
        </div>
      </div>

      {secondaryImages.length > 0 ? (
        <div className="grid grid-cols-4 gap-3">
          {secondaryImages.map((image, index) => {
            const imageIndex = index + 1;

            return (
              <button
                key={image.id}
                type="button"
                onClick={() => openViewer(imageIndex)}
                className="overflow-hidden rounded-md border border-black/10 bg-white shadow-soft transition hover:border-moss"
                aria-label={`View ${productName} image ${imageIndex + 1} larger`}
              >
                <div className="aspect-square bg-[#ded4c9]">
                  {image.image_url ? (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img
                      src={image.image_url}
                      alt={`${productName} ${image.id}`}
                      className="h-full w-full object-cover"
                    />
                  ) : null}
                </div>
              </button>
            );
          })}
        </div>
      ) : null}

      {isViewerOpen && viewerImage?.image_url ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/82 p-4 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-label={`${productName} large image viewer`}
          onClick={() => setIsViewerOpen(false)}
        >
          <button
            type="button"
            onClick={(event) => {
              event.stopPropagation();
              setIsViewerOpen(false);
            }}
            className="absolute right-4 top-4 grid h-11 w-11 place-items-center rounded-md bg-white text-2xl font-semibold text-ink shadow-soft transition hover:bg-[#f7f3ef]"
            aria-label="Close image viewer"
          >
            &times;
          </button>

          {images.length > 1 ? (
            <button
              type="button"
              onClick={(event) => {
                event.stopPropagation();
                showAdjacentViewerImage("previous");
              }}
              className="absolute left-4 top-1/2 grid h-11 w-11 -translate-y-1/2 place-items-center rounded-md bg-white text-2xl font-semibold text-ink shadow-soft transition hover:bg-[#f7f3ef]"
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
              alt={`${productName} large view`}
              className="mx-auto max-h-[88vh] w-auto max-w-full rounded-md object-contain shadow-soft"
            />
          </div>

          {images.length > 1 ? (
            <button
              type="button"
              onClick={(event) => {
                event.stopPropagation();
                showAdjacentViewerImage("next");
              }}
              className="absolute right-4 top-1/2 grid h-11 w-11 -translate-y-1/2 place-items-center rounded-md bg-white text-2xl font-semibold text-ink shadow-soft transition hover:bg-[#f7f3ef]"
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
