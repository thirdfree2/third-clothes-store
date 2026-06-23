import {
  accessTokenCookieNames,
  accessTokenMaxAgeSeconds,
  accessTokenStorageKeys,
} from "./constants";
import type { AuthArea } from "./types";

const defaultArea: AuthArea = "customer";

export const tokenStore = {
  get(area: AuthArea = defaultArea) {
    if (typeof window === "undefined") {
      return null;
    }

    return window.localStorage.getItem(accessTokenStorageKeys[area]);
  },

  set(token: string, area: AuthArea = defaultArea) {
    window.localStorage.setItem(accessTokenStorageKeys[area], token);
    document.cookie = `${accessTokenCookieNames[area]}=${token}; path=/; max-age=${accessTokenMaxAgeSeconds}; SameSite=Lax`;
  },

  clear(area: AuthArea = defaultArea) {
    window.localStorage.removeItem(accessTokenStorageKeys[area]);
    document.cookie = `${accessTokenCookieNames[area]}=; path=/; max-age=0; SameSite=Lax`;
  },
};
