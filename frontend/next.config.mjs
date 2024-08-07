/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  images: {
    domains: [
      "wow.zamimg.com",
      "render.worldofwarcraft.com",
      "assets.rpglogs.com",
    ],
  },
};

export default nextConfig;
