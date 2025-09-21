import type { NextConfig } from 'next'

// eslint-disable-next-line node/prefer-global/process
const basePath = process.env.NODE_ENV === 'production'
  ? '/specialstandard-frontend'
  : ''

const nextConfig: NextConfig = {
  /* config options here */
  basePath,
}

export default nextConfig
