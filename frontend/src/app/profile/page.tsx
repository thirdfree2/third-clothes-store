import { RouteGuard } from "@/components/route-guard";
import { CustomerProfileDetail } from "@/components/shop/customer-profile-detail";

export default function ProfilePage() {
  return (
    <RouteGuard permissions={["profile:read"]}>
      <CustomerProfileDetail />
    </RouteGuard>
  );
}
