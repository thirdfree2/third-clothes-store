"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { authApi } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { tokenStore } from "@/lib/auth/token-store";

export default function LoginPage() {
  const router = useRouter();
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("third_test@gmail.com");
  const [password, setPassword] = useState("1234");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      if (mode === "register") {
        if (password !== confirmPassword) {
          setError("Passwords do not match");
          return;
        }

        await authApi.register({
          email,
          password,
        });
      }

      const response = await authApi.login({
        email,
        password,
        userType: "CUSTOMER",
      });

      tokenStore.set(response.access_token, "customer");
      router.replace(returnToPath());
      router.refresh();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError(mode === "register" ? "Cannot register right now" : "Cannot login right now");
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
          Clothes Store
        </p>
        <h1 className="mt-2 text-2xl font-semibold text-ink">
          {mode === "register" ? "Create account" : "Customer login"}
        </h1>

        <div className="mt-5 grid grid-cols-2 rounded-md border border-black/10 bg-stone-50 p-1">
          <button
            type="button"
            onClick={() => {
              setMode("login");
              setError(null);
            }}
            className={`h-9 rounded px-3 text-sm font-semibold ${
              mode === "login"
                ? "bg-white text-ink shadow-soft"
                : "text-black/55 hover:text-ink"
            }`}
          >
            Sign in
          </button>
          <button
            type="button"
            onClick={() => {
              setMode("register");
              setError(null);
            }}
            className={`h-9 rounded px-3 text-sm font-semibold ${
              mode === "register"
                ? "bg-white text-ink shadow-soft"
                : "text-black/55 hover:text-ink"
            }`}
          >
            Register
          </button>
        </div>

        <label className="mt-5 block text-sm font-medium text-ink">
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
            autoComplete={mode === "register" ? "new-password" : "current-password"}
            required
            className="mt-2 h-11 w-full rounded-md border border-black/15 px-3 outline-none focus:border-moss"
          />
        </label>

        {mode === "register" ? (
          <label className="mt-4 block text-sm font-medium text-ink">
            Confirm password
            <input
              value={confirmPassword}
              onChange={(event) => setConfirmPassword(event.target.value)}
              type="password"
              autoComplete="new-password"
              required
              className="mt-2 h-11 w-full rounded-md border border-black/15 px-3 outline-none focus:border-moss"
            />
          </label>
        ) : null}

        {error ? <p className="mt-4 text-sm text-clay">{error}</p> : null}

        <button
          type="submit"
          disabled={isSubmitting}
          className="mt-6 h-11 w-full rounded-md bg-ink px-4 font-medium text-white hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isSubmitting
            ? mode === "register"
              ? "Creating account..."
              : "Signing in..."
            : mode === "register"
              ? "Create account"
              : "Sign in"}
        </button>

        <div className="mt-4 flex items-center justify-between gap-3 text-sm">
          <Link href="/" className="font-medium text-moss">
            Back to shop
          </Link>
        </div>
      </form>
    </main>
  );
}

function returnToPath(): string {
  const searchParams = new URLSearchParams(window.location.search);
  const returnTo = searchParams.get("returnTo");

  if (
    !returnTo ||
    !returnTo.startsWith("/") ||
    returnTo.startsWith("//") ||
    returnTo.startsWith("/login")
  ) {
    return "/";
  }

  return returnTo;
}
