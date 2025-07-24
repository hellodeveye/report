<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-gray-900/10">
    <div class="max-w-md w-full bg-white/60 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/30 p-8">
      <!-- Logo 和标题 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-indigo-500 mb-4">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-gray-800 mb-2">报告助手</h1>
        <p class="text-gray-600">请选择登录方式继续使用</p>
      </div>

      <!-- 登录按钮 -->
      <div class="space-y-4">
        <!-- 飞书登录 -->
        <button @click="handleLogin('feishu')"
                class="w-full flex items-center justify-center bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-4 rounded-lg shadow-md transition-transform transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-opacity-75">
          <svg v-if="isLoading && currentProvider === 'feishu'" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-5 h-5 mr-3" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12.59 13.92L7.74 19.4a.9.9 0 01-1.27 0l-1.9-1.9a.9.9 0 010-1.28l10.3-10.3a.9.9 0 011.27 0l1.9 1.9a.9.9 0 010 1.27L12.59 13.92z" />
            <path d="M18.91 4.5l-1.9-1.9a.9.9 0 00-1.27 0L5.44 12.9a.9.9 0 000 1.28l1.9 1.9" />
            <path d="M14.73 3.23l1.9 1.9a.9.9 0 010 1.27L8.98 14.05" />
          </svg>
          {{ (isLoading && currentProvider === 'feishu') ? '正在跳转...' : '使用飞书登录' }}
        </button>

        <!-- 钉钉登录按钮 -->
        <button 
          @click="() => handleLogin('dingtalk')" 
          :disabled="isLoading"
          class="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-lg text-white font-semibold bg-blue-500 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg v-if="isLoading && currentProvider === 'dingtalk'" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-5 h-5 mr-3" fill="currentColor" viewBox="0 0 24 24">
            <path d="M22.46,6C21.64,6.35,20.76,6.58,19.83,6.69C20.79,6.11,21.5,5.26,21.84,4.24C21,4.72,20.06,5.08,19.07,5.32C18.25,4.45,17,4,15.61,4C13,4,10.9,6.1,10.9,8.71c0,0.36,0.04,0.71,0.12,1.06C7.23,9.58,3.84,7.84,1.5,5.13C1.1,5.78,0.89,6.54,0.89,7.38c0,1.56,0.81,2.94,2.03,3.76 C2.17,11.08,1.47,10.85,0.87,10.5v0.03c0,2.18,1.54,4,3.57,4.42C4.12,15,3.75,15.06,3.36,15.06c-0.28,0-0.56-0.03-0.82-0.08 c0.57,1.78,2.24,3.08,4.22,3.12C5.29,19.34,3.6,19.98,1.77,19.98c-0.3,0-0.6-0.02-0.89-0.05C2.8,20.89,5.01,21.6,7.38,21.6 c7.8,0,12.07-6.46,12.07-12.07c0-0.18,0-0.37-0.01-0.55C20.35,8.24,21.01,7.52,21.5,6.72c-0.74,0.33-1.53,0.55-2.36,0.65 c0.85-0.51,1.5-1.32,1.81-2.27L21.5,6.72z" />
          </svg>
          {{ (isLoading && currentProvider === 'dingtalk') ? '正在跳转...' : '使用钉钉登录' }}
        </button>

        <!-- 登录提示 -->
        <div class="text-center">
          <p class="text-sm text-gray-500">
            点击登录后将跳转到对应平台的认证页面
          </p>
          <p class="text-xs text-gray-400 mt-1">
            您也可以使用手机扫码快速登录
          </p>
        </div>
      </div>

              <!-- 功能说明 -->
        <div class="mt-8 pt-6 border-t border-gray-200/50">
          <h3 class="text-sm font-semibold text-gray-700 mb-3">产品特性</h3>
          <ul class="space-y-2 text-xs text-gray-600">
            <li class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
              </svg>
              自动汇总历史报告内容
            </li>
            <li class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
              </svg>
              AI智能生成报告草稿
            </li>
            <li class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
              </svg>
              支持多种报告模板
            </li>
            <li class="flex items-center">
              <svg class="w-4 h-4 mr-2 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
              </svg>
              支持飞书和钉钉平台集成
            </li>
          </ul>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { authService } from '../utils/authService.js';

const isLoading = ref(false);
const currentProvider = ref('');

const handleLogin = async (provider) => {
  try {
    isLoading.value = true;
    currentProvider.value = provider;
    await authService.login(provider);
  } catch (error) {
    console.error('Login error:', error);
    isLoading.value = false;
    currentProvider.value = '';
  }
};
</script>

<style scoped>
/* 背景渐变动画 */
div.min-h-screen::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image:
    radial-gradient(circle at 15% 50%, #6f7bf7 0%, transparent 25%),
    radial-gradient(circle at 85% 30%, #4f5bce 0%, transparent 25%),
    radial-gradient(circle at 60% 80%, #a855f7 0%, transparent 25%);
  filter: blur(80px);
  z-index: -1;
  opacity: 0.5;
  animation: gradient-shift 10s ease-in-out infinite alternate;
}

@keyframes gradient-shift {
  0% {
    filter: blur(80px) hue-rotate(0deg);
  }
  100% {
    filter: blur(80px) hue-rotate(30deg);
  }
}
</style> 