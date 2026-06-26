import Link from "next/link";
import { notFound } from "next/navigation";
import { AddToCartPanel } from "@/components/shop/add-to-cart-panel";
import { CartCountLink } from "@/components/shop/cart-count-link";
import { ProductImageGallery } from "@/components/shop/product-image-gallery";
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

  const categoryNames = product.categories.map((category) => category.name);

  return (
    <main className="min-h-screen bg-[linear-gradient(180deg,#fbf8f4_0%,#f7f3ef_48%,#efe7df_100%)] text-ink">
      <header className="mx-auto flex w-full max-w-6xl flex-wrap items-center justify-between gap-4 px-5 py-5">
        <Link href="/" className="flex items-center gap-3" aria-label="Back to shop">
          <span className="grid h-10 w-10 place-items-center rounded-md border border-black/10 bg-white text-lg font-semibold text-ink shadow-soft">
            &larr;
          </span>
          <span>
            <span className="block text-sm font-semibold text-moss">Back to shop</span>
            <span className="block text-xs font-medium text-black/50">
              Clothes Store
            </span>
          </span>
        </Link>
        <nav className="flex flex-wrap items-center gap-2" aria-label="Product navigation">
          <CartCountLink />
          <Link
            href="/profile"
            className="inline-flex h-10 items-center rounded-md border border-black/10 bg-white/90 px-4 text-sm font-semibold text-ink shadow-soft transition hover:border-moss hover:text-moss"
          >
            User profile
          </Link>
        </nav>
      </header>

      <section className="mx-auto grid w-full max-w-6xl gap-8 px-5 pb-12 pt-4 lg:grid-cols-[minmax(0,1.02fr)_minmax(360px,0.82fr)] lg:items-start">
        <ProductImageGallery
          images={product.images}
          productID={product.id}
          productName={product.name}
        />

        <aside className="rounded-md border border-black/10 bg-white/90 p-5 shadow-soft backdrop-blur lg:sticky lg:top-5">
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-clay">
            Ready-to-wear
          </p>
          <h1 className="mt-3 text-4xl font-semibold leading-tight text-ink">
            {product.name}
          </h1>
          <p className="mt-3 text-3xl font-semibold text-clay">
            {formatMoney(product.price)}
          </p>

          <div className="mt-5 flex flex-wrap gap-2">
            {product.color ? (
              <span className="inline-flex items-center gap-2 rounded-md border border-black/10 bg-white px-3 py-2 text-sm font-semibold text-ink">
                <span
                  aria-hidden="true"
                  className="h-3 w-3 rounded-full border border-black/20"
                  style={{ backgroundColor: product.color.hex_code }}
                />
                {product.color.name}
              </span>
            ) : null}
            {categoryNames.length > 0
              ? product.categories.map((category) => (
                  <span
                    key={category.id}
                    className="rounded-md bg-[#f1ebe4] px-3 py-2 text-sm font-semibold text-black/65"
                  >
                    {category.name}
                  </span>
                ))
              : null}
          </div>

          <div className="mt-6 grid gap-3 border-y border-black/10 py-5 sm:grid-cols-3">
            <Info label="Images" value={String(product.images.length)} />
            <Info label="Updated" value={formatDate(product.updated_at)} />
            <Info label="Currency" value="THB" />
          </div>

          <div className="mt-5 rounded-md bg-[#f7f3ef] p-4">
            <p className="text-sm font-semibold text-ink">Store note</p>
            <p className="mt-1 text-sm leading-6 text-black/58">
              Choose your size before adding this piece to your cart. You can
              update quantity later from the cart page.
            </p>
          </div>

          <AddToCartPanel clothesID={product.id} />
        </aside>
      </section>

      <section className="mx-auto w-full max-w-6xl px-5 pb-12">
        <div className="grid gap-4 border-t border-black/10 pt-6 sm:grid-cols-3">
          <Feature title="Easy styling" text="Pairs cleanly with everyday layers." />
          <Feature title="Size first" text="Required size selection keeps cart items clear." />
          <Feature title="Curated rack" text="Part of the current Clothes Store collection." />
        </div>
      </section>
    </main>
  );
}

function Feature({ title, text }: { title: string; text: string }) {
  return (
    <div className="rounded-md border border-black/10 bg-white/76 p-4 shadow-soft">
      <p className="font-semibold text-ink">{title}</p>
      <p className="mt-1 text-sm leading-6 text-black/55">{text}</p>
    </div>
  );
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-semibold uppercase tracking-[0.16em] text-black/42">
        {label}
      </dt>
      <dd className="mt-1 text-sm font-semibold text-ink">{value}</dd>
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
  }).format(new Date(value));
}
