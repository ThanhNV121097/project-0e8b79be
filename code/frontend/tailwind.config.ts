import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{js,ts,jsx,tsx,mdx}"],
  theme: {
    extend: {
      colors: {
        background: "#F8FAFC",
        surface: "#FFFFFF",
        text: "#0F172A",
        muted: "#64748B",
        primary: "#2563EB",
        primaryHover: "#1D4ED8",
        primarySoft: "#EFF6FF",
        success: "#10B981",
        danger: "#EF4444"
      },
      boxShadow: {
        app: "0 24px 70px rgba(37,99,235,.14)"
      },
      borderRadius: {
        app: "24px"
      }
    }
  },
  plugins: []
};

export default config;
