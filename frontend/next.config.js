/** @type {import('next').NextConfig} */
const isProd = process.env.NODE_ENV === "production";

const nextConfig = {
  output: isProd ? "export" : undefined,
  distDir: isProd ? "build" : ".next",
  images: {
    unoptimized: true,
  },
  trailingSlash: false,
  devIndicators: false,
};

module.exports = nextConfig;
