import { type ReactNode } from "react";
import { BackofficeShell } from "@/components/backoffice/backoffice-shell";

export default function BackofficeProtectedLayout({
  children,
}: {
  children: ReactNode;
}) {
  return <BackofficeShell>{children}</BackofficeShell>;
}
