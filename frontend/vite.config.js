import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    // This is useful for docker-compose setups
    host: '0.0.0.0',
    port: 5173,
    // Proxy API requests to the backend server
    proxy: {
      '/api': {
        // 根据环境变量选择目标地址
        // 本地开发: http://localhost:8080
        // Docker环境: http://backend:8080
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        // 添加请求日志
        configure: (proxy, options) => {
          proxy.on('proxyReq', (proxyReq, req, res) => {
            console.log('Proxying request:', req.method, req.url, '-> ', options.target + req.url);
          });
          proxy.on('error', (err, req, res) => {
            console.error('Proxy error:', err);
          });
        }
      },
    },
  }
}) 