// app/layout.tsx ou src/app/layout.tsx
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import ReactQueryProvider from "@/providers/ReactQueryProvider";
import Layout from "@/components/Layout";
import { AuthProvider } from "@/providers/AuthContext";
import { Toaster } from "react-hot-toast";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
});

export const metadata: Metadata = {
  title: "WoW Perf - Insight. Optimize. Conquer.",
  description:
    "Your ultimate companion for Mythic+, Raids and PvP analytics and character improvement.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <body
        className={`${inter.className} antialiased bg-background text-foreground`}
      >
        <ReactQueryProvider>
          <AuthProvider>
            <Layout>{children}</Layout>
            <Toaster position="bottom-right" />
          </AuthProvider>
        </ReactQueryProvider>
      </body>
    </html>
  );
}
