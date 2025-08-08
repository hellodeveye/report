import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { themeService } from './utils/themeService.js'

// 初始化主题（在应用挂载前）
themeService.initTheme()

createApp(App).mount('#app') 