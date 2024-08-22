/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
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
    ],
  },
};

export default nextConfig;
