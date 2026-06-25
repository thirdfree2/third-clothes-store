"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { CartCountLink } from "@/components/shop/cart-count-link";
import { ApiError } from "@/lib/api/client";
import { customerApi } from "@/lib/api/customer";
import { purchasesApi } from "@/lib/api/purchases";
import type { CustomerProfile, Purchase } from "@/lib/api/types";
import { decodeAccessToken } from "@/lib/auth/jwt";
import { tokenStore } from "@/lib/auth/token-store";

type ProfileState =
  | { status: "loading"; profile: null; error: null }
  | { status: "ready"; profile: CustomerProfile; error: null }
  | { status: "not-found"; profile: null; error: string }
  | { status: "error"; profile: null; error: string };

type PurchaseHistoryState =
  | { status: "loading"; purchases: Purchase[]; error: null }
  | { status: "ready"; purchases: Purchase[]; error: null }
  | { status: "error"; purchases: Purchase[]; error: string };

type ProfileFormState = {
  first_name: string;
  last_name: string;
  phone: string;
  date_of_birth: string;
  marketing_opt_in: boolean;
  address: {
    recipient_name: string;
    phone: string;
    address_line1: string;
    address_line2: string;
    subdistrict: string;
    district: string;
    province: string;
    postal_code: string;
    country_code: string;
  };
};

export function CustomerProfileDetail() {
  const router = useRouter();
  const [state, setState] = useState<ProfileState>({
    status: "loading",
    profile: null,
    error: null,
  });
  const [email, setEmail] = useState<string | null>(null);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [form, setForm] = useState<ProfileFormState>(emptyProfileForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [historyState, setHistoryState] = useState<PurchaseHistoryState>({
    status: "loading",
    purchases: [],
    error: null,
  });

  useEffect(() => {
    void loadProfile();
    void loadPurchaseHistory();
  }, []);

  async function loadProfile() {
    const token = tokenStore.get("customer");

    if (!token) {
      setState({
        status: "error",
        profile: null,
        error: "Customer token is missing",
      });
      return;
    }

    setEmail(decodeAccessToken(token)?.email ?? null);

    try {
      const profile = await customerApi.getMe(token);
      setState({ status: "ready", profile, error: null });
      setForm(profileToForm(profile));
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Cannot load profile";

      setState({
        status: isProfileNotFound(err) ? "not-found" : "error",
        profile: null,
        error: message,
      });
    }
  }

  async function updateMe() {
    const token = tokenStore.get("customer");

    if (!token) {
      setState({
        status: "error",
        profile: null,
        error: "Customer token is missing",
      });
      return;
    }

    setIsUpdating(true);

    try {
      const profile = await customerApi.updateMe(token, {
        first_name: null,
        last_name: null,
        phone: null,
        date_of_birth: null,
        marketing_opt_in: false,
      });

      setState({ status: "ready", profile, error: null });
      setForm(profileToForm(profile));
    } catch (err) {
      setState({
        status: "error",
        profile: null,
        error: err instanceof ApiError ? err.message : "Cannot update profile",
      });
    } finally {
      setIsUpdating(false);
    }
  }

  async function saveProfile() {
    const token = tokenStore.get("customer");

    if (!token) {
      setFormError("Customer token is missing");
      return;
    }

    const validationError = validateProfileForm(form);

    if (validationError) {
      setFormError(validationError);
      return;
    }

    setIsUpdating(true);
    setFormError(null);

    try {
      const profile = await customerApi.updateMe(token, {
        first_name: emptyToNull(form.first_name),
        last_name: emptyToNull(form.last_name),
        phone: emptyToNull(form.phone),
        date_of_birth: emptyToNull(form.date_of_birth),
        marketing_opt_in: form.marketing_opt_in,
        address: {
          recipient_name: form.address.recipient_name.trim(),
          phone: form.address.phone.trim(),
          address_line1: form.address.address_line1.trim(),
          address_line2: emptyToNull(form.address.address_line2),
          subdistrict: emptyToNull(form.address.subdistrict),
          district: form.address.district.trim(),
          province: form.address.province.trim(),
          postal_code: form.address.postal_code.trim(),
          country_code: form.address.country_code.trim().toUpperCase() || "TH",
        },
      });

      setState({ status: "ready", profile, error: null });
      setForm(profileToForm(profile));
      setIsEditing(false);
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : "Cannot update profile");
    } finally {
      setIsUpdating(false);
    }
  }

  function startEditing(profile: CustomerProfile) {
    setForm(profileToForm(profile));
    setFormError(null);
    setIsEditing(true);
  }

  function cancelEditing(profile: CustomerProfile) {
    setForm(profileToForm(profile));
    setFormError(null);
    setIsEditing(false);
  }

  async function loadPurchaseHistory() {
    const token = tokenStore.get("customer");

    if (!token) {
      setHistoryState({
        status: "error",
        purchases: [],
        error: "Customer token is missing",
      });
      return;
    }

    try {
      const purchases = await purchasesApi.list(token);
      setHistoryState({ status: "ready", purchases, error: null });
    } catch (err) {
      setHistoryState({
        status: "error",
        purchases: [],
        error: err instanceof ApiError ? err.message : "Cannot load purchases",
      });
    }
  }

  function logout() {
    tokenStore.clear("customer");
    router.replace("/login");
    router.refresh();
  }

  if (state.status === "loading") {
    return <ProfileMessage message="Loading profile..." />;
  }

  if (state.status === "error") {
    return <ProfileMessage message={state.error} tone="error" />;
  }

  if (state.status === "not-found") {
    return (
      <ProfileShell onLogout={logout}>
        <section className="mt-8 rounded-md border border-black/10 bg-white p-6 text-center shadow-soft">
          <h2 className="text-xl font-semibold text-ink">Profile not found</h2>
          <p className="mt-2 text-sm text-black/55">{state.error}</p>
          <button
            type="button"
            onClick={updateMe}
            disabled={isUpdating}
            className="mt-5 inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isUpdating ? "Updating..." : "UpdateMe"}
          </button>
        </section>
      </ProfileShell>
    );
  }

  const { profile } = state;
  const fullName = [profile.first_name, profile.last_name].filter(Boolean).join(" ");
  const completedFields = [
    profile.first_name,
    profile.last_name,
    profile.phone,
    profile.date_of_birth,
  ].filter(Boolean).length;

  return (
    <ProfileShell
      email={email}
      fullName={fullName}
      onLogout={logout}
      orderCount={historyState.purchases.length}
    >
      <section className="mt-6 grid gap-4 sm:grid-cols-3">
        <SummaryCard value={String(historyState.purchases.length)} label="Orders" />
        <SummaryCard value={`${completedFields}/4`} label="Profile fields" />
        <SummaryCard
          value={profile.marketing_opt_in ? "On" : "Off"}
          label="Marketing"
        />
      </section>

      <section className="mt-5 rounded-md border border-black/10 bg-white/90 p-5 shadow-soft">
        <div className="flex flex-wrap items-start justify-between gap-4 border-b border-black/10 pb-4">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.2em] text-clay">
              Account details
            </p>
            <h2 className="mt-2 text-2xl font-semibold text-ink">
              Personal information
            </h2>
          </div>
          <span className="rounded-md bg-[#f7f3ef] px-3 py-2 text-sm font-semibold text-black/55">
            Customer #{profile.user_id}
          </span>
        </div>
        {isEditing ? (
          <ProfileEditForm
            form={form}
            error={formError}
            isSaving={isUpdating}
            onCancel={() => cancelEditing(profile)}
            onChange={(nextForm) => {
              setForm(nextForm);
              setFormError(null);
            }}
            onSave={saveProfile}
          />
        ) : (
          <>
            <dl className="mt-5 grid gap-3 sm:grid-cols-2">
              <Info label="Email" value={email ?? "None"} />
              <Info label="Name" value={fullName || "None"} />
              <Info label="Phone" value={profile.phone ?? "None"} />
              <Info label="Date of birth" value={profile.date_of_birth ?? "None"} />
              <Info
                label="Marketing"
                value={profile.marketing_opt_in ? "Opted in" : "Opted out"}
              />
              <Info label="Address" value={formatAddress(profile.address)} />
            </dl>
            <button
              type="button"
              onClick={() => startEditing(profile)}
              className="mt-5 inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white transition hover:bg-black"
            >
              Edit profile
            </button>
          </>
        )}
      </section>
      <PurchaseHistory state={historyState} />
    </ProfileShell>
  );
}

function ProfileShell({
  children,
  email,
  fullName,
  onLogout,
  orderCount = 0,
}: {
  children: React.ReactNode;
  email?: string | null;
  fullName?: string;
  onLogout: () => void;
  orderCount?: number;
}) {
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
        <div className="flex flex-wrap items-center gap-2">
          <CartCountLink />
          <button
            type="button"
            onClick={onLogout}
            className="inline-flex h-10 items-center rounded-md border border-clay/30 bg-white/90 px-4 text-sm font-semibold text-clay shadow-soft transition hover:border-clay"
          >
            Logout
          </button>
        </div>
      </header>

      <div className="mx-auto w-full max-w-6xl px-5 pb-12 pt-3">
        <section className="rounded-md border border-black/10 bg-white/90 p-5 shadow-soft sm:p-6">
          <div className="flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-4">
              <div className="grid h-16 w-16 shrink-0 place-items-center rounded-md bg-ink text-2xl font-semibold text-white">
                {getInitials(fullName, email)}
              </div>
              <div>
                <p className="text-sm font-semibold uppercase tracking-[0.2em] text-clay">
                  Customer account
                </p>
                <h1 className="mt-2 text-3xl font-semibold text-ink">
                  {fullName || "User profile"}
                </h1>
                <p className="mt-1 text-sm text-black/55">
                  {email ?? "Email not available"}
                </p>
              </div>
            </div>
            <div className="rounded-md bg-[#f7f3ef] px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-[0.16em] text-black/42">
                Items to receive
              </p>
              <p className="mt-1 text-2xl font-semibold text-clay">{orderCount}</p>
            </div>
          </div>
        </section>

        {children}
      </div>
    </main>
  );
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md border border-black/10 bg-[#fbf8f4] p-4">
      <dt className="text-xs font-semibold uppercase tracking-[0.16em] text-black/42">
        {label}
      </dt>
      <dd className="mt-2 font-semibold text-ink">{value}</dd>
    </div>
  );
}

function SummaryCard({ value, label }: { value: string; label: string }) {
  return (
    <div className="rounded-md border border-black/10 bg-white/76 p-4 shadow-soft">
      <p className="text-2xl font-semibold text-ink">{value}</p>
      <p className="mt-1 text-xs font-semibold uppercase tracking-[0.16em] text-black/45">
        {label}
      </p>
    </div>
  );
}

function ProfileEditForm({
  error,
  form,
  isSaving,
  onCancel,
  onChange,
  onSave,
}: {
  error: string | null;
  form: ProfileFormState;
  isSaving: boolean;
  onCancel: () => void;
  onChange: (form: ProfileFormState) => void;
  onSave: () => void;
}) {
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        void onSave();
      }}
      className="mt-5"
    >
      <div className="grid gap-4 sm:grid-cols-2">
        <ProfileInput
          label="First name"
          value={form.first_name}
          onChange={(value) => onChange({ ...form, first_name: value })}
          placeholder="Your first name"
        />
        <ProfileInput
          label="Last name"
          value={form.last_name}
          onChange={(value) => onChange({ ...form, last_name: value })}
          placeholder="Your last name"
        />
        <ProfileInput
          label="Phone"
          value={form.phone}
          onChange={(value) => onChange({ ...form, phone: value })}
          placeholder="0812345678"
        />
        <ProfileInput
          label="Date of birth"
          type="date"
          value={form.date_of_birth}
          onChange={(value) => onChange({ ...form, date_of_birth: value })}
        />
      </div>

      <div className="mt-6 border-t border-black/10 pt-5">
        <p className="text-sm font-semibold uppercase tracking-[0.2em] text-clay">
          Delivery address
        </p>
        <div className="mt-4 grid gap-4 sm:grid-cols-2">
          <ProfileInput
            label="Recipient name"
            value={form.address.recipient_name}
            onChange={(value) =>
              onChange({
                ...form,
                address: { ...form.address, recipient_name: value },
              })
            }
            placeholder="Name for delivery"
          />
          <ProfileInput
            label="Address phone"
            value={form.address.phone}
            onChange={(value) =>
              onChange({ ...form, address: { ...form.address, phone: value } })
            }
            placeholder="0812345678"
          />
          <ProfileInput
            label="Address line 1"
            value={form.address.address_line1}
            onChange={(value) =>
              onChange({
                ...form,
                address: { ...form.address, address_line1: value },
              })
            }
            placeholder="House number, street, building"
          />
          <ProfileInput
            label="Address line 2"
            value={form.address.address_line2}
            onChange={(value) =>
              onChange({
                ...form,
                address: { ...form.address, address_line2: value },
              })
            }
            placeholder="Room, floor, landmark"
          />
          <ProfileInput
            label="Subdistrict"
            value={form.address.subdistrict}
            onChange={(value) =>
              onChange({
                ...form,
                address: { ...form.address, subdistrict: value },
              })
            }
            placeholder="Subdistrict"
          />
          <ProfileInput
            label="District"
            value={form.address.district}
            onChange={(value) =>
              onChange({ ...form, address: { ...form.address, district: value } })
            }
            placeholder="District"
          />
          <ProfileInput
            label="Province"
            value={form.address.province}
            onChange={(value) =>
              onChange({ ...form, address: { ...form.address, province: value } })
            }
            placeholder="Province"
          />
          <div className="grid gap-4 sm:grid-cols-[1fr_96px]">
            <ProfileInput
              label="Postal code"
              value={form.address.postal_code}
              onChange={(value) =>
                onChange({
                  ...form,
                  address: { ...form.address, postal_code: value },
                })
              }
              placeholder="10110"
            />
            <ProfileInput
              label="Country"
              value={form.address.country_code}
              onChange={(value) =>
                onChange({
                  ...form,
                  address: { ...form.address, country_code: value },
                })
              }
              placeholder="TH"
            />
          </div>
        </div>
      </div>

      <label className="mt-4 flex items-start gap-3 rounded-md border border-black/10 bg-[#fbf8f4] p-4">
        <input
          type="checkbox"
          checked={form.marketing_opt_in}
          onChange={(event) =>
            onChange({ ...form, marketing_opt_in: event.target.checked })
          }
          className="mt-1 h-4 w-4 accent-[#426851]"
        />
        <span>
          <span className="block text-sm font-semibold text-ink">
            Marketing updates
          </span>
          <span className="mt-1 block text-sm leading-6 text-black/55">
            Receive store updates, product drops, and style notes.
          </span>
        </span>
      </label>

      {error ? (
        <p className="mt-4 rounded-md bg-clay/10 px-3 py-2 text-sm font-semibold text-clay">
          {error}
        </p>
      ) : null}

      <div className="mt-5 flex flex-wrap gap-2">
        <button
          type="submit"
          disabled={isSaving}
          className="inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isSaving ? "Saving..." : "Save profile"}
        </button>
        <button
          type="button"
          onClick={onCancel}
          disabled={isSaving}
          className="inline-flex h-10 items-center rounded-md border border-black/15 bg-white px-4 text-sm font-semibold text-ink transition hover:border-moss hover:text-moss disabled:cursor-not-allowed disabled:opacity-60"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}

function ProfileInput({
  label,
  onChange,
  placeholder,
  type = "text",
  value,
}: {
  label: string;
  onChange: (value: string) => void;
  placeholder?: string;
  type?: "date" | "text";
  value: string;
}) {
  return (
    <label className="text-sm font-semibold text-ink">
      {label}
      <input
        type={type}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        className="mt-2 h-11 w-full rounded-md border border-black/15 bg-white px-3 text-sm font-medium text-ink outline-none transition placeholder:text-black/35 focus:border-moss"
      />
    </label>
  );
}

function ProfileMessage({
  message,
  tone = "muted",
}: {
  message: string;
  tone?: "muted" | "error";
}) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[#f7f3ef] px-5">
      <div className="rounded-md border border-black/10 bg-white p-6 text-center shadow-soft">
        <p className={tone === "error" ? "font-semibold text-clay" : "text-black/55"}>
          {message}
        </p>
      </div>
    </main>
  );
}

function PurchaseHistory({ state }: { state: PurchaseHistoryState }) {
  if (state.status === "loading") {
    return (
      <section className="mt-5 rounded-md border border-black/10 bg-white/90 p-5 shadow-soft">
        <h2 className="text-2xl font-semibold text-ink">Items to receive</h2>
        <p className="mt-3 text-sm text-black/55">Loading orders...</p>
      </section>
    );
  }

  if (state.status === "error") {
    return (
      <section className="mt-5 rounded-md border border-black/10 bg-white/90 p-5 shadow-soft">
        <h2 className="text-2xl font-semibold text-ink">Items to receive</h2>
        <p className="mt-3 text-sm text-clay">{state.error}</p>
      </section>
    );
  }

  if (state.purchases.length === 0) {
    return (
      <section className="mt-5 rounded-md border border-black/10 bg-white/90 p-8 text-center shadow-soft">
        <h2 className="text-2xl font-semibold text-ink">Items to receive</h2>
        <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-black/55">
          No orders yet. When you checkout, your incoming pieces will appear here.
        </p>
        <Link
          href="/"
          className="mt-5 inline-flex h-10 items-center rounded-md bg-ink px-4 text-sm font-semibold text-white transition hover:bg-black"
        >
          Shop now
        </Link>
      </section>
    );
  }

  return (
    <section className="mt-5 rounded-md border border-black/10 bg-white/90 p-5 shadow-soft">
      <div className="flex flex-wrap items-end justify-between gap-3 border-b border-black/10 pb-4">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-clay">
            Orders
          </p>
          <h2 className="mt-2 text-2xl font-semibold text-ink">Items to receive</h2>
        </div>
        <p className="text-sm font-semibold text-black/50">
          {state.purchases.length} order{state.purchases.length === 1 ? "" : "s"}
        </p>
      </div>
      <div className="mt-5 space-y-4">
        {state.purchases.map((purchase) => (
          <article
            key={purchase.id}
            className="rounded-md border border-black/10 bg-[#fbf8f4] p-4"
          >
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <h3 className="font-semibold text-ink">Order #{purchase.id}</h3>
                <p className="mt-1 text-sm text-black/55">
                  {purchase.created_at ? formatDate(purchase.created_at) : "Order date not available"}
                </p>
              </div>
              <p className="rounded-md bg-white px-3 py-2 font-semibold text-clay">
                {formatMoney(purchase.total_amount)}
              </p>
            </div>

            <div className="mt-4 space-y-3">
              {purchase.items.map((item) => {
                const productHref = `/products/${item.clothes_id}`;

                return (
                  <div
                    key={item.id}
                    className="grid gap-3 rounded-md border border-black/10 bg-white p-3 sm:grid-cols-[76px_1fr_auto]"
                  >
                    <Link
                      href={productHref}
                      className="block h-[76px] w-[76px] overflow-hidden rounded-md bg-[#ded4c9]"
                    >
                      {item.product_image_url ? (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          src={item.product_image_url}
                          alt={item.product_name || `Product #${item.clothes_id}`}
                          className="h-full w-full object-cover"
                        />
                      ) : null}
                    </Link>
                    <div className="min-w-0">
                      <Link href={productHref} className="font-semibold text-ink">
                        {item.product_name || `Product #${item.clothes_id}`}
                      </Link>
                      <p className="mt-1 text-sm text-black/55">
                        Size: {item.size || "No size"} / Qty: {item.quantity}
                      </p>
                      <p className="mt-2 inline-flex rounded-md bg-moss/10 px-2 py-1 text-xs font-semibold text-moss">
                        Waiting to receive
                      </p>
                    </div>
                    <p className="font-semibold text-ink sm:text-right">
                      {formatMoney(item.unit_price * item.quantity)}
                    </p>
                  </div>
                );
              })}
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}

function getInitials(fullName?: string, email?: string | null) {
  const source = fullName || email || "User";
  const parts = source
    .replace(/@.*/, "")
    .split(/\s+/)
    .filter(Boolean);

  return parts
    .slice(0, 2)
    .map((part) => part[0])
    .join("")
    .toUpperCase();
}

const emptyProfileForm: ProfileFormState = {
  first_name: "",
  last_name: "",
  phone: "",
  date_of_birth: "",
  marketing_opt_in: false,
  address: {
    recipient_name: "",
    phone: "",
    address_line1: "",
    address_line2: "",
    subdistrict: "",
    district: "",
    province: "",
    postal_code: "",
    country_code: "TH",
  },
};

function profileToForm(profile: CustomerProfile): ProfileFormState {
  return {
    first_name: profile.first_name ?? "",
    last_name: profile.last_name ?? "",
    phone: profile.phone ?? "",
    date_of_birth: profile.date_of_birth ?? "",
    marketing_opt_in: profile.marketing_opt_in,
    address: {
      recipient_name: profile.address?.recipient_name ?? "",
      phone: profile.address?.phone ?? profile.phone ?? "",
      address_line1: profile.address?.address_line1 ?? "",
      address_line2: profile.address?.address_line2 ?? "",
      subdistrict: profile.address?.subdistrict ?? "",
      district: profile.address?.district ?? "",
      province: profile.address?.province ?? "",
      postal_code: profile.address?.postal_code ?? "",
      country_code: profile.address?.country_code ?? "TH",
    },
  };
}

function emptyToNull(value: string): string | null {
  const trimmedValue = value.trim();
  return trimmedValue ? trimmedValue : null;
}

function validateProfileForm(form: ProfileFormState): string | null {
  const phone = form.phone.trim();

  if (phone && !/^[0-9+\-\s()]{8,20}$/.test(phone)) {
    return "Please enter a valid phone number.";
  }

  const addressPhone = form.address.phone.trim();
  if (!form.address.recipient_name.trim()) {
    return "Recipient name is required.";
  }
  if (!addressPhone || !/^[0-9+\-\s()]{8,20}$/.test(addressPhone)) {
    return "Please enter a valid address phone number.";
  }
  if (!form.address.address_line1.trim()) {
    return "Address line 1 is required.";
  }
  if (!form.address.district.trim()) {
    return "District is required.";
  }
  if (!form.address.province.trim()) {
    return "Province is required.";
  }
  if (!form.address.postal_code.trim()) {
    return "Postal code is required.";
  }
  if (form.address.country_code.trim().length !== 2) {
    return "Country code must be 2 characters.";
  }

  if (form.date_of_birth) {
    const birthday = new Date(form.date_of_birth);
    const today = new Date();

    if (Number.isNaN(birthday.getTime()) || birthday > today) {
      return "Date of birth cannot be in the future.";
    }
  }

  return null;
}

function formatAddress(address: CustomerProfile["address"]): string {
  if (!address) {
    return "None";
  }

  return [
    address.recipient_name,
    address.phone,
    address.address_line1,
    address.address_line2,
    address.subdistrict,
    address.district,
    address.province,
    address.postal_code,
    address.country_code,
  ]
    .filter(Boolean)
    .join(", ");
}

function isProfileNotFound(err: unknown) {
  return (
    err instanceof ApiError &&
    (err.code === 404101 || err.message === "customer profile not found")
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
