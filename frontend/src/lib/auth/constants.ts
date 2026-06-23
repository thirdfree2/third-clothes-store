export const accessTokenMaxAgeSeconds = 60 * 60;

export const accessTokenStorageKeys = {
  customer: "clothes_store_customer_access_token",
  admin: "clothes_store_admin_access_token",
} as const;

export const accessTokenCookieNames = {
  customer: "clothes_store_customer_access_token",
  admin: "clothes_store_admin_access_token",
} as const;
