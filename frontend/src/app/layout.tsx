import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import ReactQueryProvider from "@/providers/ReactQueryProvider";
import Layout from "@/components/Layout";
import { AuthProvider } from "@/providers/AuthContext";
import { Toaster } from "react-hot-toast";

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
      <body className={`${inter.className} bg-[#002440]`}>
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
