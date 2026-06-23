export default function BackofficeOrdersPage() {
  return (
    <section className="mx-auto w-full max-w-7xl px-5 py-6">
      <div>
        <p className="text-sm font-medium uppercase tracking-wide text-moss">
          Sales
        </p>
        <h1 className="mt-2 text-3xl font-semibold">Orders</h1>
      </div>

      <div className="mt-6 overflow-hidden rounded-md border border-black/10 bg-white shadow-soft">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[760px] border-collapse text-left text-sm">
            <thead className="bg-stone-100 text-xs uppercase tracking-wide text-black/55">
              <tr>
                <th className="px-4 py-3 font-semibold">Order</th>
                <th className="px-4 py-3 font-semibold">Customer</th>
                <th className="px-4 py-3 font-semibold">Items</th>
                <th className="px-4 py-3 text-right font-semibold">Total</th>
                <th className="px-4 py-3 text-right font-semibold">Created</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-black/55">
                  Orders are not available in backoffice yet.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
}
