/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "standalone",
  assetPrefix: process.env.NODE_ENV === "production" ? "/" : "",
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
    NEXT_PUBLIC_APP_URL: process.env.NEXT_PUBLIC_APP_URL,
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "wow.zamimg.com",
      },
      {
        protocol: "https",
        hostname: "render.worldofwarcraft.com",
      },
      {
        protocol: "https",
        hostname: "assets.rpglogs.com",
      },
      {
        protocol: "https",
        hostname: "cdn.raiderio.net",
      },
    ],
  },
  // Configuration for watching files in Docker
  experimental: {
    webpackBuildWorker: false,
  },
};

export default nextConfig;
