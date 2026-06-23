"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { authApi } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { tokenStore } from "@/lib/auth/token-store";

export default function BackofficeLoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("admin_editor@example.com");
  const [password, setPassword] = useState("1234");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      const response = await authApi.login({
        email,
        password,
        userType: "ADMIN",
      });

      tokenStore.set(response.access_token, "admin");
      router.replace(returnToPath());
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Cannot login right now");
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md flex-col justify-center px-5">
      <form
        onSubmit={handleSubmit}
        className="rounded-md border border-black/10 bg-white p-6 shadow-soft"
      >
        <p className="text-sm font-medium uppercase tracking-wide text-moss">
          Backoffice
        </p>
        <h1 className="mt-2 text-2xl font-semibold text-ink">Admin login</h1>

        <label className="mt-6 block text-sm font-medium text-ink">
          Email
          <input
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            type="email"
            autoComplete="email"
            required
            className="mt-2 h-11 w-full rounded-md border border-black/15 px-3 outline-none focus:border-moss"
          />
        </label>

        <label className="mt-4 block text-sm font-medium text-ink">
          Password
          <input
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            type="password"
            autoComplete="current-password"
            required
            className="mt-2 h-11 w-full rounded-md border border-black/15 px-3 outline-none focus:border-moss"
          />
        </label>

        {error ? <p className="mt-4 text-sm text-clay">{error}</p> : null}

        <button
          type="submit"
          disabled={isSubmitting}
          className="mt-6 h-11 w-full rounded-md bg-ink px-4 font-medium text-white disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isSubmitting ? "Signing in..." : "Sign in"}
        </button>
      </form>
    </main>
  );
}

function returnToPath(): string {
  const searchParams = new URLSearchParams(window.location.search);
  const returnTo = searchParams.get("returnTo");

  if (!returnTo || returnTo.startsWith("/backoffice/login")) {
    return "/backoffice";
  }

  return returnTo;
}
