import { clothesApi } from "@/lib/api/clothes";
import { CartCountLink } from "@/components/shop/cart-count-link";
import Link from "next/link";
import type { ReactNode } from "react";

const perPage = 9;

type HomePageProps = {
  searchParams?: Promise<{
    page?: string | string[];
  }>;
};

export default async function HomePage({ searchParams }: HomePageProps) {
  const params = await searchParams;
  const currentPage = parsePage(params?.page);
  const clothes = await clothesApi
    .list({ page: currentPage, perPage })
    .catch(() => null);
  const pagination = clothes?.pagination;
  const totalPages = pagination?.total_pages ?? 0;
  const pageNumbers = getPageNumbers(pagination?.page ?? currentPage, totalPages);
  const featuredItem = clothes?.items[0];
  const totalItems = pagination?.total ?? clothes?.items.length ?? 0;

  return (
    <main className="min-h-screen overflow-hidden bg-[linear-gradient(180deg,#fbf8f4_0%,#f7f3ef_42%,#efe7df_100%)] text-ink">
      <header className="mx-auto flex w-full max-w-6xl flex-wrap items-center justify-between gap-4 px-5 py-5">
        <Link href="/" className="flex items-center gap-3" aria-label="Clothes Store">
          <span className="grid h-11 w-11 place-items-center rounded-md bg-ink text-sm font-black text-white shadow-soft">
            CS
          </span>
          <span>
            <span className="block text-base font-semibold uppercase tracking-wide">
              Clothes Store
            </span>
            <span className="block text-xs font-medium text-black/55">
              Everyday pieces, neatly curated
            </span>
          </span>
        </Link>
        <nav className="flex flex-wrap items-center gap-2" aria-label="Shop navigation">
          <CartCountLink />
          <Link
            href="/profile"
            className="inline-flex h-10 items-center rounded-md border border-black/10 bg-white/90 px-4 text-sm font-semibold text-ink shadow-soft transition hover:border-moss hover:text-moss"
          >
            User profile
          </Link>
        </nav>
      </header>

      <section className="mx-auto grid w-full max-w-6xl gap-8 px-5 pb-8 pt-3 lg:min-h-[520px] lg:grid-cols-[minmax(0,0.95fr)_minmax(360px,0.8fr)] lg:items-center">
        <div className="max-w-2xl">
          <p className="text-sm font-semibold uppercase tracking-[0.22em] text-moss">
            New wardrobe arrivals
          </p>
          <h1 className="mt-4 max-w-2xl text-5xl font-semibold leading-[0.98] text-ink sm:text-6xl">
            Style that feels ready for today.
          </h1>
          <p className="mt-5 max-w-xl text-base leading-7 text-black/62">
            Browse relaxed essentials, polished layers, and color-forward pieces
            selected for easy daily dressing.
          </p>
          <div className="mt-8 flex flex-wrap items-center gap-3">
            <a
              href="#collection"
              className="inline-flex h-12 items-center rounded-md bg-ink px-5 text-sm font-semibold text-white shadow-soft transition hover:bg-black"
            >
              Shop collection
            </a>
            <CartCountLink
              className="inline-flex h-12 items-center rounded-md border border-black/12 bg-white px-5 text-sm font-semibold text-ink shadow-soft transition hover:border-moss hover:text-moss"
            >
              View cart
            </CartCountLink>
          </div>
          <div className="mt-9 grid max-w-lg grid-cols-3 gap-3">
            <Stat value={String(totalItems)} label="Pieces" />
            <Stat value="S-XL" label="Sizes" />
            <Stat value="THB" label="Currency" />
          </div>
        </div>

        <div className="relative min-h-[360px]">
          <div className="absolute left-3 top-5 h-48 w-32 rounded-md bg-clay/16" />
          <div className="absolute bottom-4 right-0 h-44 w-44 rounded-md bg-moss/14" />
          <Link
            href={featuredItem ? `/products/${featuredItem.id}` : "#collection"}
            className="group relative mx-auto block max-w-[420px] overflow-hidden rounded-md border border-black/10 bg-white shadow-soft"
          >
            <div className="aspect-[3/4] bg-[#ded4c9]">
              {featuredItem?.images[0]?.image_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={featuredItem.images[0].image_url}
                  alt={featuredItem.name}
                  className="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                />
              ) : (
                <div className="flex h-full items-center justify-center p-8 text-center text-sm font-medium uppercase tracking-[0.2em] text-black/38">
                  Clothes Store
                </div>
              )}
            </div>
            <div className="absolute inset-x-4 bottom-4 rounded-md border border-white/50 bg-white/92 p-4 shadow-soft backdrop-blur">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-moss">
                Featured
              </p>
              <div className="mt-2 flex items-end justify-between gap-4">
                <div>
                  <h2 className="text-lg font-semibold text-ink">
                    {featuredItem?.name ?? "Curated clothing"}
                  </h2>
                  <p className="mt-1 text-sm text-black/55">
                    {featuredItem
                      ? categoryText(featuredItem.categories)
                      : "Fresh looks arriving soon"}
                  </p>
                </div>
                {featuredItem ? (
                  <p className="shrink-0 text-base font-semibold text-clay">
                    {formatMoney(featuredItem.price)}
                  </p>
                ) : null}
              </div>
            </div>
          </Link>
        </div>
      </section>

      <section
        id="collection"
        className="mx-auto w-full max-w-6xl px-5 pb-12 pt-4"
      >
        <div className="mb-5 flex flex-wrap items-end justify-between gap-4 border-b border-black/10 pb-4">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.2em] text-clay">
              Collection
            </p>
            <h2 className="mt-2 text-3xl font-semibold text-ink">
              Ready-to-wear picks
            </h2>
          </div>
          <p className="text-sm font-medium text-black/55">
            Page {currentPage}
            {totalPages > 0 ? ` of ${totalPages}` : ""}
          </p>
        </div>

        {clothes ? (
          clothes.items.length > 0 ? (
            <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
              {clothes.items.map((item) => (
                <ProductCard key={item.id} item={item} />
              ))}
            </div>
          ) : (
            <EmptyState
              title="No clothes in this collection yet"
              description="Add products from the backoffice and they will appear here."
            />
          )
        ) : (
          <EmptyState
            title="The rack is taking a moment"
            description="Products could not be loaded right now. Please check the API service and refresh."
          />
        )}
      </section>

      {totalPages > 1 ? (
        <nav
          aria-label="Product pagination"
          className="mx-auto flex w-full max-w-6xl flex-wrap items-center justify-center gap-2 px-5 pb-12"
        >
          <PaginationLink
            href={pageHref(Math.max(1, currentPage - 1))}
            isDisabled={currentPage <= 1}
          >
            Previous
          </PaginationLink>

          {pageNumbers.map((pageNumber) => (
            <PaginationLink
              key={pageNumber}
              href={pageHref(pageNumber)}
              isActive={pageNumber === currentPage}
            >
              {pageNumber}
            </PaginationLink>
          ))}

          <PaginationLink
            href={pageHref(Math.min(totalPages, currentPage + 1))}
            isDisabled={currentPage >= totalPages}
          >
            Next
          </PaginationLink>
        </nav>
      ) : null}
    </main>
  );
}

type ProductCardProps = {
  item: NonNullable<Awaited<ReturnType<typeof clothesApi.list>>>["items"][number];
};

function ProductCard({ item }: ProductCardProps) {
  const primaryImage = item.images[0];

  return (
    <Link
      href={`/products/${item.id}`}
      className="group overflow-hidden rounded-md border border-black/10 bg-white shadow-soft transition duration-200 hover:-translate-y-1 hover:border-moss/50"
    >
      <div className="relative aspect-[4/5] bg-[#ded4c9]">
        {primaryImage?.image_url ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={primaryImage.image_url}
            alt={item.name}
            className="h-full w-full object-cover transition duration-500 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full items-center justify-center p-8 text-center text-xs font-semibold uppercase tracking-[0.22em] text-black/35">
            No image
          </div>
        )}
        {item.color ? (
          <span className="absolute left-3 top-3 rounded-md bg-white/92 px-3 py-1 text-xs font-semibold text-ink shadow-soft backdrop-blur">
            {item.color.name}
          </span>
        ) : null}
      </div>
      <div className="p-4">
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0">
            <h3 className="truncate text-base font-semibold text-ink">{item.name}</h3>
            <p className="mt-1 line-clamp-1 text-sm text-black/55">
              {categoryText(item.categories)}
            </p>
          </div>
          <p className="shrink-0 text-sm font-bold text-clay">
            {formatMoney(item.price)}
          </p>
        </div>
        <div className="mt-4 flex items-center justify-between border-t border-black/10 pt-3">
          <span className="text-xs font-semibold uppercase tracking-[0.16em] text-black/42">
            View details
          </span>
          <span className="grid h-8 w-8 place-items-center rounded-md bg-ink text-white transition group-hover:bg-moss">
            &rarr;
          </span>
        </div>
      </div>
    </Link>
  );
}

function Stat({ value, label }: { value: string; label: string }) {
  return (
    <div className="rounded-md border border-black/10 bg-white/72 p-4 shadow-soft backdrop-blur">
      <p className="text-xl font-semibold text-ink">{value}</p>
      <p className="mt-1 text-xs font-semibold uppercase tracking-[0.16em] text-black/45">
        {label}
      </p>
    </div>
  );
}

function EmptyState({
  title,
  description,
}: {
  title: string;
  description: string;
}) {
  return (
    <div className="rounded-md border border-dashed border-black/18 bg-white/76 px-5 py-12 text-center shadow-soft">
      <p className="text-lg font-semibold text-ink">{title}</p>
      <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-black/55">
        {description}
      </p>
    </div>
  );
}

function PaginationLink({
  children,
  href,
  isActive = false,
  isDisabled = false,
}: {
  children: ReactNode;
  href: string;
  isActive?: boolean;
  isDisabled?: boolean;
}) {
  const className = [
    "inline-flex h-10 min-w-10 items-center justify-center rounded-md border px-3 text-sm font-semibold shadow-soft transition",
    isActive
      ? "border-ink bg-ink text-white"
      : "border-black/10 bg-white text-ink hover:border-moss hover:text-moss",
    isDisabled ? "pointer-events-none opacity-45" : "",
  ].join(" ");

  return (
    <Link
      href={href}
      aria-disabled={isDisabled}
      aria-current={isActive ? "page" : undefined}
      className={className}
    >
      {children}
    </Link>
  );
}

function parsePage(value: string | string[] | undefined): number {
  const pageValue = Array.isArray(value) ? value[0] : value;
  const page = Number(pageValue);

  if (!Number.isInteger(page) || page < 1) {
    return 1;
  }

  return page;
}

function pageHref(page: number): string {
  return page === 1 ? "/" : `/?page=${page}`;
}

function categoryText(categories: { name: string }[]): string {
  return categories.length > 0
    ? categories.map((category) => category.name).join(", ")
    : "Everyday wear";
}

function formatMoney(value: number): string {
  return `${value.toLocaleString("th-TH")} THB`;
}

function getPageNumbers(currentPage: number, totalPages: number): number[] {
  const start = Math.max(1, currentPage - 2);
  const end = Math.min(totalPages, currentPage + 2);

  return Array.from({ length: end - start + 1 }, (_, index) => start + index);
}
