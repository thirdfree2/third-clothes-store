import type { ApiEnvelope } from "./types";

// const defaultBaseURL = "http://localhost:9080";
// const defaultBaseURL = "https://third-shop.duckdns.org";
const defaultBaseURL = "http://192.168.1.102:9080";


export class ApiError extends Error {
  readonly status: number;
  readonly code?: number;
  readonly details?: Record<string, unknown>;

  constructor(
    message: string,
    options: {
      status: number;
      code?: number;
      details?: Record<string, unknown>;
    },
  ) {
    super(message);
    this.name = "ApiError";
    this.status = options.status;
    this.code = options.code;
    this.details = options.details;
  }
}

export type ApiRequestOptions = Omit<RequestInit, "body"> & {
  token?: string | null;
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined | null>;
};

export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions = {},
): Promise<T> {
  const { token, body, query, headers, ...init } = options;
  const url = buildURL(path, query);

  const response = await fetch(url, {
    ...init,
    headers: {
      Accept: "application/json",
      ...(body === undefined ? {} : { "Content-Type": "application/json" }),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });

  const envelope = (await response.json().catch(() => null)) as
    | ApiEnvelope<T>
    | null;

  if (!response.ok || !envelope?.success) {
    throw new ApiError(
      envelope?.error?.message ?? `Request failed with status ${response.status}`,
      {
        status: response.status,
        code: envelope?.error?.code,
        details: envelope?.error?.details,
      },
    );
  }

  return envelope.payload;
}

export type ApiFormDataRequestOptions = Omit<RequestInit, "body"> & {
  token?: string | null;
  body: FormData;
  query?: ApiRequestOptions["query"];
};

export async function apiFormDataFetch<T>(
  path: string,
  options: ApiFormDataRequestOptions,
): Promise<T> {
  const { token, body, query, headers, ...init } = options;
  const url = buildURL(path, query);

  const response = await fetch(url, {
    ...init,
    headers: {
      Accept: "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
    body,
  });

  const envelope = (await response.json().catch(() => null)) as
    | ApiEnvelope<T>
    | null;

  if (!response.ok || !envelope?.success) {
    throw new ApiError(
      envelope?.error?.message ?? `Request failed with status ${response.status}`,
      {
        status: response.status,
        code: envelope?.error?.code,
        details: envelope?.error?.details,
      },
    );
  }

  return envelope.payload;
}

function buildURL(
  path: string,
  query?: ApiRequestOptions["query"],
): string {
  const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL ?? defaultBaseURL;
  const url = new URL(path, baseURL);

  if (query) {
    for (const [key, value] of Object.entries(query)) {
      if (value === undefined || value === null) {
        continue;
      }

      url.searchParams.set(key, String(value));
    }
  }

  return url.toString();
}
