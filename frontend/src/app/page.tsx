import { clothesApi } from "@/lib/api/clothes";
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

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-8 px-5 py-8">
      <header className="flex flex-wrap items-center justify-between gap-4 border-b border-black/10 pb-5">
        <div>
          <p className="text-sm font-medium uppercase tracking-wide text-moss">
            Clothes Store
          </p>
          <h1 className="mt-2 text-3xl font-semibold text-ink">Shop</h1>
        </div>
        <nav className="flex flex-wrap items-center gap-2">
          <Link
            href="/cart"
            className="inline-flex h-10 items-center rounded-md border border-black/10 bg-white px-4 text-sm font-semibold text-ink shadow-soft hover:border-moss hover:text-moss"
          >
            Cart
          </Link>
          <Link
            href="/profile"
            className="inline-flex h-10 items-center rounded-md border border-black/10 bg-white px-4 text-sm font-semibold text-ink shadow-soft hover:border-moss hover:text-moss"
          >
            User profile
          </Link>
          <div className="rounded-md bg-white px-4 py-2 text-sm shadow-soft">
            API: {process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:9080"}
          </div>
        </nav>
      </header>

      <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {clothes?.items.map((item) => (
          <Link
            key={item.id}
            href={`/products/${item.id}`}
            className="overflow-hidden rounded-md border border-black/10 bg-white shadow-soft"
          >
            <div className="aspect-[4/3] bg-stone-200">
              {item.images[0]?.image_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={item.images[0].image_url}
                  alt={item.name}
                  className="h-full w-full object-cover"
                />
              ) : null}
            </div>
            <div className="flex items-start justify-between gap-3 p-4">
              <div>
                <h2 className="font-semibold text-ink">{item.name}</h2>
                <p className="mt-1 text-sm text-black/55">
                  {item.categories.map((category) => category.name).join(", ")}
                </p>
              </div>
              <p className="shrink-0 font-semibold text-clay">
                {item.price.toLocaleString()} THB
              </p>
            </div>
          </Link>
        ))}
      </section>

      {totalPages > 1 ? (
        <nav
          aria-label="Product pagination"
          className="flex flex-wrap items-center justify-center gap-2"
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
    "inline-flex h-10 min-w-10 items-center justify-center rounded-md border px-3 text-sm font-semibold shadow-soft",
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

function getPageNumbers(currentPage: number, totalPages: number): number[] {
  const start = Math.max(1, currentPage - 2);
  const end = Math.min(totalPages, currentPage + 2);

  return Array.from({ length: end - start + 1 }, (_, index) => start + index);
}
