import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#151515",
        linen: "#f7f3ef",
        moss: "#426851",
        clay: "#b45f43",
      },
      boxShadow: {
        soft: "0 14px 40px rgba(21, 21, 21, 0.08)",
      },
    },
  },
  plugins: [],
};

export default config;
