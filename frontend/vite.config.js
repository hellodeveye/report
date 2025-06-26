import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [vue()],
    server: {
      // This is useful for docker-compose setups
      host: '0.0.0.0',
      port: 5173,
      // 只在开发环境使用代理
      proxy: mode === 'dev' ? {
        '/api': {
          target: 'http://backend:8080',
          changeOrigin: true,
        },
      } : undefined,
    },
    // 定义全局常量
    define: {
      __API_BASE_URL__: JSON.stringify(env.VITE_API_BASE_URL || 'http://localhost:8080/api'),
    },
  }
}) 