export type ApiErrorDetail = {
  code: number;
  message: string;
  details?: Record<string, unknown>;
};

export type ApiEnvelope<T> = {
  success: boolean;
  payload: T;
  error?: ApiErrorDetail;
};

export type PaginationMeta = {
  page: number;
  perPage: number;
  total: number;
  total_pages: number;
};

export type PaginatedPayload<T> = {
  items: T[];
  pagination: PaginationMeta;
};

export type Color = {
  id: number;
  name: string;
  hex_code: string;
};

export type Category = {
  id: number;
  name: string;
};

export type ClothesImage = {
  id: number;
  clothes_id: number;
  bucket?: string;
  object_key?: string;
  image_url: string;
  original_filename?: string;
  content_type?: string;
  size_bytes?: number;
  created_at?: string;
};

export type Clothes = {
  id: number;
  name: string;
  price: number;
  color_id: number | null;
  color?: Color;
  categories: Category[];
  images: ClothesImage[];
  created_at: string;
  updated_at: string;
};

export type CartItemClothes = {
  id: number;
  name: string;
  price: number;
  color_id: number | null;
  images: Pick<ClothesImage, "id" | "clothes_id" | "image_url">[];
};

export type CartItem = {
  id: number;
  clothes_id: number;
  clothes?: CartItemClothes;
  quantity: number;
  size: string;
};

export type Cart = {
  id: number;
  user_id: number;
  status: string;
  items: CartItem[];
};

export type CustomerProfile = {
  user_id: number;
  first_name: string | null;
  last_name: string | null;
  phone: string | null;
  date_of_birth: string | null;
  marketing_opt_in: boolean;
  address: CustomerAddress | null;
};

export type UpdateCustomerProfileInput = {
  first_name: string | null;
  last_name: string | null;
  phone: string | null;
  date_of_birth: string | null;
  marketing_opt_in: boolean;
  address?: UpdateCustomerAddressInput;
};

export type CustomerAddress = {
  id: number;
  recipient_name: string;
  phone: string;
  address_line1: string;
  address_line2: string | null;
  subdistrict: string | null;
  district: string;
  province: string;
  postal_code: string;
  country_code: string;
  is_default: boolean;
};

export type UpdateCustomerAddressInput = {
  recipient_name: string;
  phone: string;
  address_line1: string;
  address_line2: string | null;
  subdistrict: string | null;
  district: string;
  province: string;
  postal_code: string;
  country_code: string;
};

export type PurchaseItem = {
  id: number;
  purchase_id: number;
  clothes_id: number;
  product_name: string;
  product_image_url: string | null;
  quantity: number;
  unit_price: number;
  size: string;
};

export type Purchase = {
  id: number;
  user_id: number;
  total_amount: number;
  items: PurchaseItem[];
  created_at?: string;
  updated_at?: string;
};

export type LoginResponse = {
  access_token: string;
  token_type: "Bearer";
};

export type RegisterResponse = {
  id: number;
  email: string;
  user_type: string;
  status: string;
};
