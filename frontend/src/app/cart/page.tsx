import { RouteGuard } from "@/components/route-guard";
import { CartDetail } from "@/components/shop/cart-detail";

export default function CartPage() {
  return (
    <RouteGuard permissions={["cart:read"]}>
      <CartDetail />
    </RouteGuard>
  );
}
