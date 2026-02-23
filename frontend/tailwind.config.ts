import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        sakura: "#f8b4c8",
      },
    },
  },
  plugins: [],
};

export default config;
