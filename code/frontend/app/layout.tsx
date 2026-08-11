import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Todo List App v2",
  description: "A simple database-backed todo list app."
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
