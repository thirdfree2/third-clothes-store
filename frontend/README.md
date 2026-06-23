# Clothes Store Frontend

Next.js frontend for the clothes-store services.

## Setup

```bash
npm install
npm run dev
```

Copy `.env.example` to `.env.local` if the API gateway is not running at `http://localhost:9080`.

## API Lib

- `src/lib/api/auth.ts`
- `src/lib/api/clothes.ts`
- `src/lib/api/cart.ts`
- `src/lib/api/purchases.ts`

Protected APIs expect the JWT access token from `authApi.login`.
