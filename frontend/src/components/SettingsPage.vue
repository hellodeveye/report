<script setup>
import { ref, defineEmits, defineProps, onMounted, computed } from 'vue';
import { aiService, MODELS_CONFIG } from '../utils/aiUtils.js';

const props = defineProps({
  currentUser: {
    type: Object,
    default: () => null
  }
});

const emit = defineEmits(['close', 'save-settings']);

const activeTab = ref('account');
const tabs = [
  { id: 'account', name: '账户资料', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' },
  { id: 'model', name: '模型设置', icon: 'M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 10-3.86.517l-.318.158a6 6 0 00-3.86.517L6.05 15.428a2 2 0 00-1.82.875 2 2 0 00.166 2.304l.558.837a6 6 0 1011.9.001l.558-.836a2 2 0 00.166-2.304z' },
  { id: 'appearance', name: '外观', icon: 'M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.373 3.373 0 0012 18.443V21' },
  { id: 'notifications', name: '通知', icon: 'M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9' },
  { id: 'language', name: '语言', icon: 'M3 5h12M9 3v2m4.25 16a4.25 4.25 0 000-8.5M15 12a4.25 4.25 0 00-8.5 0M12.5 17h-5' },
];

const aiConfig = ref({
  provider: 'deepseek',
  apiKey: '',
  model: 'deepseek-chat',
});

const settings = ref({
  darkMode: false,
  notifications: {
    email: true,
    push: false,
  },
  language: 'zh-cn',
});

const onProviderChange = () => {
  const provider = aiConfig.value.provider;
  if (MODELS_CONFIG[provider] && MODELS_CONFIG[provider].models.length > 0) {
    aiConfig.value.model = MODELS_CONFIG[provider].models[0].id;
  }
};

const saveSettings = () => {
  aiService.saveSettings(aiConfig.value);
  // Here you would also save other settings (appearance, notifications, etc.)
  emit('save-settings', { ...settings.value, ai: aiConfig.value });
  emit('close');
};

const selectedProviderModels = computed(() => {
  return MODELS_CONFIG[aiConfig.value.provider]?.models || [];
});

onMounted(() => {
  aiConfig.value = aiService.getSettings();
});
</script>

<template>
  <div class="flex-grow flex flex-col overflow-hidden bg-white/50">
    <div class="p-4 border-b border-white/30 flex-shrink-0 flex items-center">
      <button @click="emit('close')" class="flex items-center text-sm font-semibold text-gray-600 hover:text-gray-800 transition-colors">
        <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
        返回主界面
      </button>
    </div>
    <div class="flex-grow flex overflow-hidden">
      <!-- Left: Tab navigation -->
      <aside class="w-56 flex-shrink-0 p-4 border-r border-white/30 overflow-y-auto">
        <nav class="space-y-2">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'w-full flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors',
              activeTab === tab.id
                ? 'bg-indigo-100/80 text-indigo-700 shadow-sm'
                : 'text-gray-600 hover:bg-gray-500/10 hover:text-gray-800',
            ]"
          >
            <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="tab.icon"></path></svg>
            <span>{{ tab.name }}</span>
          </button>
        </nav>
      </aside>
      
      <!-- Right: Content -->
      <main class="flex-grow overflow-y-auto">
        <div class="p-8 max-w-3xl mx-auto w-full">
          <transition name="fade" mode="out-in">
            <div :key="activeTab">
              <!-- Account Settings -->
              <div v-if="activeTab === 'account'">
                <h3 class="text-xl font-bold text-gray-800 mb-6">账户资料</h3>
                <div class="bg-white/80 rounded-lg shadow p-8 border border-white/30">
                  <div class="flex items-center space-x-6">
                    <!-- Avatar -->
                    <img v-if="props.currentUser?.avatar_url"
                         :src="props.currentUser.avatar_url"
                         :alt="props.currentUser.name"
                         class="h-24 w-24 rounded-full object-cover border-4 border-white shadow-lg">
                    <span v-else
                          class="inline-flex items-center justify-center h-24 w-24 rounded-full text-white text-3xl font-bold bg-indigo-500">
                      {{ props.currentUser?.name?.charAt(0)?.toUpperCase() || 'U' }}
                    </span>
                    <!-- User Info -->
                    <div class="flex-grow">
                      <h4 class="text-2xl font-bold text-gray-800">{{ props.currentUser?.name || '用户' }}</h4>
                      <p class="text-gray-500">{{ props.currentUser?.email || '无可用邮箱' }}</p>
                      <span v-if="props.currentUser?.provider" class="mt-2 inline-block px-3 py-1 text-xs font-semibold rounded-full"
                            :class="{
                              'bg-indigo-100 text-indigo-800': props.currentUser.provider === 'feishu',
                              'bg-blue-100 text-blue-800': props.currentUser.provider === 'dingtalk',
                              'bg-gray-100 text-gray-800': !props.currentUser.provider
                            }">
                        {{ props.currentUser.provider === 'feishu' ? '飞书用户' : (props.currentUser.provider === 'dingtalk' ? '钉钉用户' : '未知来源') }}
                      </span>
                    </div>
                  </div>
                  <div class="mt-8 border-t border-gray-200/60 pt-6">
                    <dl class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                      <div class="col-span-1">
                        <dt class="text-sm font-medium text-gray-500">UserID</dt>
                        <dd class="mt-1 text-sm text-gray-900">{{ props.currentUser?.user_id || '--' }}</dd>
                      </div>
                      <div class="col-span-1">
                        <dt class="text-sm font-medium text-gray-500">OpenID</dt>
                        <dd class="mt-1 text-sm text-gray-900">{{ props.currentUser?.open_id || '--' }}</dd>
                      </div>
                    </dl>
                  </div>
                </div>
              </div>

              <!-- Model Settings -->
              <div v-if="activeTab === 'model'">
                <h3 class="text-xl font-bold text-gray-800 mb-6">模型设置</h3>
                <div class="space-y-6">
                  <div class="bg-white/80 rounded-lg shadow p-6 border border-white/30 space-y-4">
                    <!-- AI Provider -->
                    <div>
                      <label for="ai-provider" class="block text-sm font-semibold text-gray-700 mb-1">AI 提供商</label>
                      <select id="ai-provider" v-model="aiConfig.provider" @change="onProviderChange" class="form-input">
                        <option v-for="(config, provider) in MODELS_CONFIG" :key="provider" :value="provider">
                          {{ config.label }}
                        </option>
                      </select>
                    </div>
                    <!-- API Key -->
                    <div>
                      <label for="api-key" class="block text-sm font-semibold text-gray-700 mb-1">API Key</label>
                      <input id="api-key" type="password" v-model="aiConfig.apiKey" placeholder="请输入您的 API Key" class="form-input" />
                    </div>
                    <!-- Model Selection -->
                    <div>
                      <label for="ai-model" class="block text-sm font-semibold text-gray-700 mb-1">模型</label>
                      <select id="ai-model" v-model="aiConfig.model" class="form-input">
                        <option v-for="model in selectedProviderModels" :key="model.id" :value="model.id">
                          {{ model.name }}
                        </option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Appearance Settings -->
              <div v-if="activeTab === 'appearance'">
                <h3 class="text-xl font-bold text-gray-800 mb-6">外观</h3>
                <div class="space-y-6">
                  <div class="bg-white/80 rounded-lg shadow p-6 border border-white/30">
                    <div class="flex items-center justify-between">
                      <div>
                        <p class="font-semibold text-gray-700">深色模式</p>
                        <p class="text-sm text-gray-500 mt-1">为界面启用或禁用深色主题。</p>
                      </div>
                      <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" v-model="settings.darkMode" class="sr-only peer">
                        <div class="w-11 h-6 bg-gray-200 rounded-full peer peer-focus:ring-4 peer-focus:ring-indigo-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
                      </label>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Notifications Settings -->
              <div v-if="activeTab === 'notifications'">
                <h3 class="text-xl font-bold text-gray-800 mb-6">通知</h3>
                <div class="space-y-6">
                  <div class="bg-white/80 rounded-lg shadow p-6 border border-white/30 space-y-4 divide-y divide-gray-200/50">
                    <div class="flex items-center justify-between pt-4 first:pt-0">
                       <div>
                        <p class="font-semibold text-gray-700">邮件通知</p>
                        <p class="text-sm text-gray-500 mt-1">接收关于账户活动和更新的邮件。</p>
                      </div>
                      <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" v-model="settings.notifications.email" class="sr-only peer">
                        <div class="w-11 h-6 bg-gray-200 rounded-full peer peer-focus:ring-4 peer-focus:ring-indigo-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
                      </label>
                    </div>
                    <div class="flex items-center justify-between pt-4 first:pt-0">
                      <div>
                        <p class="font-semibold text-gray-400">推送通知</p>
                        <p class="text-sm text-gray-400 mt-1">通过您的设备接收实时推送通知。</p>
                      </div>
                      <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" v-model="settings.notifications.push" class="sr-only peer" disabled>
                        <div class="w-11 h-6 bg-gray-200 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600 cursor-not-allowed"></div>
                        <span class="ml-3 text-sm font-medium text-gray-400">（即将推出）</span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Language Settings -->
              <div v-if="activeTab === 'language'">
                <h3 class="text-xl font-bold text-gray-800 mb-6">语言</h3>
                 <div class="space-y-6">
                  <div class="bg-white/80 rounded-lg shadow p-6 border border-white/30">
                    <div>
                      <p class="font-semibold text-gray-700 mb-1">界面语言</p>
                      <p class="text-sm text-gray-500 mb-4">选择您希望在应用中显示的语言。</p>
                      <select v-model="settings.language" class="w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                        <option value="zh-cn">简体中文</option>
                        <option value="en" disabled>English (即将推出)</option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Save Button -->
              <div class="mt-10 pt-6 border-t border-gray-200/60 flex justify-end">
                <button @click="saveSettings" class="px-6 py-2 bg-indigo-600 text-white font-semibold rounded-lg shadow-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors">
                  保存设置
                </button>
              </div>
            </div>
          </transition>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.1s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style> 