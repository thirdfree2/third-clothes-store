import Link from "next/link";
import { notFound } from "next/navigation";
import { AddToCartPanel } from "@/components/shop/add-to-cart-panel";
import { clothesApi } from "@/lib/api/clothes";

type ProductDetailPageProps = {
  params: Promise<{
    id: string;
  }>;
};

export default async function ProductDetailPage({ params }: ProductDetailPageProps) {
  const { id: rawID } = await params;
  const id = Number(rawID);

  if (!Number.isInteger(id) || id <= 0) {
    notFound();
  }

  const product = await clothesApi.getById(id).catch(() => null);

  if (!product) {
    notFound();
  }

  const primaryImage = product.images[0];
  const categoryNames = product.categories.map((category) => category.name);

  return (
    <main className="mx-auto min-h-screen w-full max-w-6xl px-5 py-8">
      <header className="flex flex-wrap items-center justify-between gap-4 border-b border-black/10 pb-5">
        <div>
          <Link href="/" className="text-sm font-medium text-moss">
            Back to shop
          </Link>
          <h1 className="mt-2 text-3xl font-semibold text-ink">{product.name}</h1>
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
          <p className="rounded-md bg-white px-4 py-2 text-lg font-semibold text-clay shadow-soft">
            {formatMoney(product.price)}
          </p>
        </nav>
      </header>

      <section className="grid gap-8 py-8 lg:grid-cols-[minmax(0,1.05fr)_minmax(320px,0.95fr)]">
        <div className="space-y-4">
          <div className="overflow-hidden rounded-md border border-black/10 bg-white shadow-soft">
            <div className="aspect-[4/3] bg-stone-200">
              {primaryImage?.image_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={primaryImage.image_url}
                  alt={product.name}
                  className="h-full w-full object-cover"
                />
              ) : null}
            </div>
          </div>

          {product.images.length > 1 ? (
            <div className="grid grid-cols-4 gap-3">
              {product.images.slice(1).map((image) => (
                <div
                  key={image.id}
                  className="overflow-hidden rounded-md border border-black/10 bg-white"
                >
                  <div className="aspect-square bg-stone-200">
                    {image.image_url ? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img
                        src={image.image_url}
                        alt={`${product.name} ${image.id}`}
                        className="h-full w-full object-cover"
                      />
                    ) : null}
                  </div>
                </div>
              ))}
            </div>
          ) : null}
        </div>

        <div className="rounded-md border border-black/10 bg-white p-5 shadow-soft">
          <dl className="grid gap-5 sm:grid-cols-2">
            <Info label="Product ID" value={String(product.id)} />
            <Info label="Price" value={formatMoney(product.price)} />
            <Info label="Color" value={product.color?.name ?? "None"} />
            <Info
              label="Categories"
              value={categoryNames.length > 0 ? categoryNames.join(", ") : "None"}
            />
            <Info label="Images" value={String(product.images.length)} />
            <Info label="Updated" value={formatDate(product.updated_at)} />
          </dl>

          {categoryNames.length > 0 ? (
            <div className="mt-6 border-t border-black/10 pt-5">
              <div className="flex flex-wrap gap-2">
                {product.categories.map((category) => (
                  <span
                    key={category.id}
                    className="rounded-md bg-stone-100 px-3 py-1.5 text-sm font-medium text-black/65"
                  >
                    {category.name}
                  </span>
                ))}
              </div>
            </div>
          ) : null}

          <AddToCartPanel clothesID={product.id} />
        </div>
      </section>
    </main>
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
