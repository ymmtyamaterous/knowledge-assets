import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "アセナレ",
  description: "お金の知識を体系的に学べる学習プラットフォーム",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
