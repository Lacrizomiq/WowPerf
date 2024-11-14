/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "standalone",
  // Important for working with Traefik
  assetPrefix: process.env.NODE_ENV === "production" ? "/_next" : "",
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
    domains: ["localhost", "api.localhost"],
  },
  // Configuration for development in HTTPS
  webpack: (config, { dev, isServer }) => {
    if (dev && !isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        https: require.resolve("https-browserify"),
        http: require.resolve("stream-http"),
      };
    }
    return config;
  },
  // Configuration for development in HTTPS
  webpackDevMiddleware: (config) => {
    config.watchOptions = {
      poll: 1000,
      aggregateTimeout: 300,
    };
    return config;
  },
  // Configuration of headers
  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          {
            key: "X-Frame-Options",
            value: "SAMEORIGIN",
          },
          {
            key: "Access-Control-Allow-Origin",
            value:
              process.env.NODE_ENV === "development"
                ? "https://localhost"
                : process.env.NEXT_PUBLIC_FRONTEND_URL,
          },
          {
            key: "Access-Control-Allow-Credentials",
            value: "true",
          },
        ],
      },
    ];
  },
};

export default nextConfig;
