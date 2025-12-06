/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/eventsub",
        destination: "https://gempbot-api.gempir.com/api/eventsub",
      },
    ];
  },
};

export default nextConfig;
