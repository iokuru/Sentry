/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  distDir: 'build',
  images: {
    unoptimized: true,
  },
  trailingSlash: false,
  devIndicators: false,
};

module.exports = nextConfig;
