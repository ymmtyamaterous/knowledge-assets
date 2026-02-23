import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/features/auth/AuthContext";
import Header from "@/components/Header";

export const metadata: Metadata = {
  title: "アセナレ",
  description: "お金の知識を体系的に学べる学習プラットフォーム",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body>
        <AuthProvider>
          <Header />
          <div className="mx-auto max-w-5xl px-4 py-6">{children}</div>
        </AuthProvider>
      </body>
    </html>
  );
}
