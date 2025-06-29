<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-gray-900/10">
    <div class="max-w-md w-full bg-white/60 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/30 p-8">
      <div class="text-center">
        <!-- 加载状态 -->
        <div v-if="isProcessing" class="space-y-4">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-500 mb-4">
            <svg class="animate-spin w-8 h-8 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-gray-800 mb-2">正在处理登录</h1>
          <p class="text-gray-600">请稍候，正在验证您的身份...</p>
        </div>

        <!-- 成功状态 -->
        <div v-else-if="isSuccess" class="space-y-4">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-500 mb-4">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-gray-800 mb-2">登录成功</h1>
          <p class="text-gray-600 mb-4">欢迎回来，{{ userInfo?.name || '用户' }}！</p>
          <p class="text-sm text-gray-500">使用{{ getProviderDisplayName(userInfo?.provider) }}账号登录</p>
          <p class="text-sm text-gray-500">正在跳转到主页...</p>
        </div>

        <!-- 错误状态 -->
        <div v-else class="space-y-4">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-500 mb-4">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-gray-800 mb-2">登录失败</h1>
          <p class="text-gray-600 mb-4">{{ errorMessage }}</p>
          
          <div class="bg-red-50 border border-red-200 rounded-md p-3 mb-4">
            <p class="text-sm text-red-800">错误详情:</p>
            <p class="text-xs text-red-600 break-all">{{ errorDetails }}</p>
          </div>

          <button @click="retryLogin" class="w-full px-4 py-2 bg-indigo-500 text-white rounded-lg hover:bg-indigo-600">
            重新登录
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { authService } from '../utils/authService.js';

const isProcessing = ref(true);
const isSuccess = ref(false);
const userInfo = ref(null);
const errorMessage = ref('');
const errorDetails = ref('');

onMounted(async () => {
  try {
    console.log('Processing auth callback...');
    
    // 处理授权回调
    const authData = await authService.handleAuthCallback();
    
    console.log('Auth callback successful:', authData);
    
    isProcessing.value = false;
    isSuccess.value = true;
    userInfo.value = authData.user;
    
    // 2秒后跳转到主页
    setTimeout(() => {
      window.location.href = '/';
    }, 2000);
    
  } catch (error) {
    console.error('Auth callback failed:', error);
    
    isProcessing.value = false;
    isSuccess.value = false;
    errorMessage.value = '登录验证失败，请重试';
    errorDetails.value = error.message || '未知错误';
  }
});

const retryLogin = () => {
  window.location.href = '/';
};

// 获取平台显示名称
const getProviderDisplayName = (provider) => {
  switch (provider) {
    case 'feishu':
      return '飞书';
    case 'dingtalk':
      return '钉钉';
    default:
      return '未知';
  }
};
</script> 