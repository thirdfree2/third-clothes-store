export default function UnauthorizedPage() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md flex-col justify-center px-5">
      <div className="rounded-md border border-black/10 bg-white p-6 shadow-soft">
        <p className="text-sm font-medium uppercase tracking-wide text-clay">
          Unauthorized
        </p>
        <h1 className="mt-2 text-2xl font-semibold text-ink">Access denied</h1>
        <p className="mt-3 text-sm leading-6 text-black/60">
          Your account does not have permission to open this page.
        </p>
      </div>
    </main>
  );
}
