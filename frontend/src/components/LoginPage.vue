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
        <!-- 飞书登录按钮 -->
        <button 
          @click="() => handleLogin('feishu')" 
          :disabled="isLoading"
          class="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-lg text-white font-semibold bg-indigo-500 hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg v-if="isLoading && currentProvider === 'feishu'" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-5 h-5 mr-3" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm0 22C6.486 22 2 17.514 2 12S6.486 2 12 2s10 4.486 10 10-4.486 10-10 10z"/>
            <path d="M8 7h8v2H8V7zm0 4h8v2H8v-2zm0 4h8v2H8v-2z"/>
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
          <svg v-else class="w-5 h-5 mr-3" viewBox="0 0 24 24" fill="currentColor">
            <path d="M22 12.65C22 6.45 17.55 2 11.35 2S.7 6.45.7 12.65C.7 18.85 5.15 23.3 11.35 23.3S22 18.85 22 12.65zM3.5 12.65C3.5 8 7.35 4.15 12 4.15s8.5 3.85 8.5 8.5-3.85 8.5-8.5 8.5-8.5-3.85-8.5-8.5z"/>
            <path d="M16.5 12.2c-.8-.4-1.3-.6-1.3-1.1 0-.3.2-.6.8-.6.5 0 1 .2 1.5.5l.7-1.3c-.7-.4-1.5-.7-2.3-.7-1.3 0-2.3.8-2.3 2 0 1.3 1 1.8 2 2.3.8.4 1.2.6 1.2 1.1 0 .4-.3.7-.9.7-.6 0-1.2-.3-1.7-.6l-.8 1.3c.7.5 1.6.8 2.5.8 1.4 0 2.4-.8 2.4-2.1 0-1.4-1-1.9-1.8-2.3z"/>
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