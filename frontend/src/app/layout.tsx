import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

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
      <body
        className={`${inter.className} bg-gradient-to-b from-[#010517] to-[#020C3A] min-h-screen`}
      >
        <main>{children}</main>
      </body>
    </html>
  );
}
