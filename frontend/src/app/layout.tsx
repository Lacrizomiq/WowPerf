import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import ReactQueryProvider from "@/providers/ReactQueryProvider";
import Layout from "@/components/Layout";
import { AuthProvider } from "@/providers/AuthContext";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "WowPerf",
  description:
    "WowPerf is a web application that allows users to track their performance in World of Warcraft.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className}`}>
        <ReactQueryProvider>
          <AuthProvider>
            <Layout>{children}</Layout>
          </AuthProvider>
        </ReactQueryProvider>
      </body>
    </html>
  );
}
