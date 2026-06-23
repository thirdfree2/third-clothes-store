"use client";

import { useParams } from "next/navigation";
import { ClothesDetail } from "@/components/backoffice/clothes-detail";

export default function BackofficeClothesDetailPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);

  if (!Number.isInteger(id) || id <= 0) {
    return (
      <main className="mx-auto flex min-h-[50vh] w-full max-w-4xl items-center justify-center px-5">
        <p className="text-clay">Invalid clothes id</p>
      </main>
    );
  }

  return <ClothesDetail id={id} />;
}
