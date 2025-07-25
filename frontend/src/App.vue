<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import TiptapEditor from './components/TiptapEditor.vue';
import TiptapViewer from './components/TiptapViewer.vue';
import VueTailwindDatepicker from 'vue-tailwind-datepicker';
import LoginPage from './components/LoginPage.vue';
import AuthCallback from './components/AuthCallback.vue';
import SettingsPage from './components/SettingsPage.vue';
import { reportSummarizer, aiService } from './utils/aiUtils.js';
import { apiService } from './utils/apiService.js';
import { authService } from './utils/authService.js';

// --- State Management ---
const isAuthenticated = ref(false);
const currentUser = ref(null);
const showAuthCallback = ref(false);
const showSettingsPage = ref(false);
const isLoading = ref(false);
const loadingText = ref('正在加载...');

const templates = ref([]);
const selectedSourceTemplateId = ref('');
const selectedTemplateId = ref('');
const formValues = ref({});
const dateValue = ref({
  startDate: '',
  endDate: '',
});

const isProfileMenuOpen = ref(false);
const profileMenuNode = ref(null);
const notifications = ref([]);

onMounted(async () => {
  document.addEventListener('click', (event) => {
    if (profileMenuNode.value && !profileMenuNode.value.contains(event.target)) {
      isProfileMenuOpen.value = false;
    }
  });
  
  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
      if (showSettingsPage.value) {
        showSettingsPage.value = false;
      }
    }
  });
  
  await checkAuthStatus();
});

const checkAuthStatus = async () => {
  isLoading.value = true;
  loadingText.value = '正在验证身份...';
  
  if (window.location.pathname === '/auth/callback' || window.location.search.includes('code=')) {
    showAuthCallback.value = true;
    isAuthenticated.value = false; // Hide main UI
    isLoading.value = false;
    return;
  }

  showAuthCallback.value = false;
  if (authService.isAuthenticated()) {
    try {
      currentUser.value = authService.getUser();
      isAuthenticated.value = true;
      await loadTemplates();
       if (!sessionStorage.getItem('login_welcomed')) {
          addNotification('登录成功', `欢迎回来，${currentUser.value.name || '用户'}！`, 'success');
          sessionStorage.setItem('login_welcomed', 'true');
       }
    } catch (error) {
      console.error('获取用户信息或模板失败:', error);
      handleLogout();
    }
  } else {
    isAuthenticated.value = false;
    currentUser.value = null;
  }
  isLoading.value = false;
};

const handleLoginSuccess = () => {
    window.history.pushState({}, document.title, "/");
    showAuthCallback.value = false;
    checkAuthStatus();
};

const handleLogout = () => {
  authService.logout();
  isAuthenticated.value = false;
  currentUser.value = null;
  sessionStorage.removeItem('login_welcomed');
};

const getProviderDisplayName = (provider) => {
  return provider === 'feishu' ? '飞书' : (provider === 'dingtalk' ? '钉钉' : '未知');
};

const getProviderAvatarClass = (provider) => {
    return provider === 'feishu' ? 'bg-indigo-500' : (provider === 'dingtalk' ? 'bg-blue-500' : 'bg-gray-500');
};



const loadTemplates = async () => {
  try {
    loadingText.value = '正在加载模板...';
    isLoading.value = true;
    const templatesData = await apiService.getTemplates();
    
    const detailedTemplates = await Promise.all(
      templatesData.map(async (t) => {
        try {
          return await apiService.getTemplateDetail(t.name, t.id);
        } catch (e) {
          console.warn(`获取模板 ${t.name} 详情失败`, e);
          return { ...t, fields: [] };
        }
      })
    );

    templates.value = detailedTemplates.filter(t => t.fields && t.fields.length > 0);
    if (templates.value.length > 0) {
      selectedSourceTemplateId.value = templates.value[0].id;
      selectedTemplateId.value = templates.value[0].id;
    }
  } catch (error) {
    console.error('加载模板失败:', error);
    addNotification('模板加载失败', error.message, 'error');
  } finally {
    isLoading.value = false;
  }
};

const removeNotification = (id) => {
  notifications.value = notifications.value.filter(n => n.id !== id);
};

const addNotification = (title, description = '', type = 'success', duration = 3000) => {
  const id = Date.now() + Math.random();
  notifications.value.push({ id, title, description, type });
  setTimeout(() => {
    removeNotification(id);
  }, duration);
};

const datePickerShortcuts = () => {
  return [
    {
      label: "今天",
      atClick: () => {
        const today = new Date();
        return [today, today];
      },
    },
    {
      label: "昨天",
      atClick: () => {
        const yesterday = new Date();
        yesterday.setDate(yesterday.getDate() - 1);
        return [yesterday, yesterday];
      },
    },
    {
      label: "本周",
      atClick: () => {
        const today = new Date();
        const dayOfWeek = today.getDay() || 7; // Make Sunday 7
        const thisWeekStart = new Date(today);
        thisWeekStart.setDate(today.getDate() - dayOfWeek + 1);
        return [thisWeekStart, today];
      },
    },
    {
      label: "上周",
       atClick: () => {
        const today = new Date();
        const dayOfWeek = today.getDay() || 7;
        const lastWeekEnd = new Date(today);
        lastWeekEnd.setDate(today.getDate() - dayOfWeek);
        const lastWeekStart = new Date(lastWeekEnd);
        lastWeekStart.setDate(lastWeekEnd.getDate() - 6);
        return [lastWeekStart, lastWeekEnd];
      },
    },
    {
      label: "本月",
      atClick: () => {
        const today = new Date();
        const thisMonthStart = new Date(today.getFullYear(), today.getMonth(), 1);
        return [thisMonthStart, today];
      },
    },
    {
      label: "上月",
      atClick: () => {
        const today = new Date();
        const lastMonthEnd = new Date(today.getFullYear(), today.getMonth(), 0);
        const lastMonthStart = new Date(lastMonthEnd.getFullYear(), lastMonthEnd.getMonth(), 1);
        return [lastMonthStart, lastMonthEnd];
      },
    },
  ];
};

const datePickerOptions = ref({
  shortcuts: {
    today: "今天",
    yesterday: "昨天",
    past: (period) => `过去 ${period} 天`,
    currentMonth: "本月",
    pastMonth: "上月",
  },
  footer: { apply: "应用", cancel: "取消" },
});

const currentTemplate = computed(() => {
  return templates.value.find(t => t.id === selectedTemplateId.value);
});

const sourceTemplate = computed(() => {
    return templates.value.find(t => t.id === selectedSourceTemplateId.value);
});


const sourceReports = ref([]);

const getReports = async () => {
  if (!selectedSourceTemplateId.value) return;
  try {
    loadingText.value = '正在获取报告...';
    isLoading.value = true;
    
    const params = { 
        rule_id: selectedSourceTemplateId.value,
        template_name: sourceTemplate.value?.name,
    };
    if (dateValue.value.startDate && dateValue.value.endDate) {
      const startDate = new Date(dateValue.value.startDate);
      params.start_time = Math.floor(startDate.getTime() / 1000);
      const endDate = new Date(dateValue.value.endDate);
      endDate.setHours(23, 59, 59, 999);
      params.end_time = Math.floor(endDate.getTime() / 1000);
    }
    
    sourceReports.value = await apiService.getReports(params, sourceTemplate.value);
    
    if (sourceReports.value.length > 0) {
      addNotification('报告获取成功', `成功获取${sourceReports.value.length}条报告`, 'success');
    } else {
      addNotification('未找到报告', '在指定的时间范围内没有找到相关报告。', 'info');
    }
  } catch (error) {
    console.error('获取报告失败:', error);
    addNotification('获取报告失败', error.message, 'error');
    sourceReports.value = [];
  } finally {
    isLoading.value = false;
  }
};

const generateDraft = async () => {
  if (!currentTemplate.value) return;
  if (!aiService.hasApiKey()) {
    addNotification('需要配置AI模型', '请先前往“设置”页面配置您的AI提供商和API Key。', 'error', 5000);
    return;
  }
  
  loadingText.value = '正在生成草稿...';
  isLoading.value = true;
  
  try {
    const summary = await reportSummarizer.summarizeReports(
      sourceReports.value,
      currentTemplate.value
    );
    formValues.value = summary;
    addNotification('AI草稿已生成', '已基于左侧报告内容智能生成草稿', 'success');
  } catch (error) {
    console.error('生成草稿失败:', error);
    addNotification('生成失败', error.message, 'error');
  } finally {
    isLoading.value = false;
  }
};

const initializeFormValues = () => {
  const values = {};
  if (currentTemplate.value) {
    currentTemplate.value.fields.forEach(field => {
      values[field.id] = field.type === 'multiSelect' || field.type === 'image' || field.type === 'attachment' ? [] : '';
    });
  }
  formValues.value = values;
};

watch(selectedTemplateId, initializeFormValues, { immediate: true });

const toggleReportDetail = (report) => {
    report.isCollapsed = !report.isCollapsed;
};

const fileInputRefs = ref({});

const triggerFileInput = (fieldId) => {
  fileInputRefs.value[fieldId]?.click();
};

const handleFileSelect = (event, field) => {
  const files = event.target.files;
  if (!files) return;

  const currentFiles = formValues.value[field.id] || [];
  if (field.maxCount && currentFiles.length + files.length > field.maxCount) {
    alert(`最多只能上传 ${field.maxCount} 个文件。`);
    return;
  }

  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    if (field.maxSize && file.size > field.maxSize) {
      alert(`文件 ${file.name} 大小超过了 ${formatFileSize(field.maxSize)} 的限制。`);
      continue;
    }
    const fileObject = {
      id: Date.now() + i,
      name: file.name,
      size: file.size,
      raw: file,
      url: URL.createObjectURL(file)
    };
    currentFiles.push(fileObject);
  }

  formValues.value[field.id] = currentFiles;
  event.target.value = ''; // Reset file input
};

const removeFile = (fieldId, fileId) => {
  const currentFiles = formValues.value[fieldId] || [];
  formValues.value[fieldId] = currentFiles.filter(f => f.id !== fileId);
};

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
</script>

<template>
  <!-- Notifications container -->
  <div aria-live="assertive" class="pointer-events-none fixed inset-0 flex items-end px-4 py-6 sm:items-start sm:p-6 z-50">
    <div class="flex w-full flex-col items-center sm:items-end">
      <transition-group name="notification" tag="div" class="w-full space-y-4 flex flex-col items-center sm:items-end" enter-active-class="transform ease-out duration-300 transition" enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2" enter-to-class="translate-y-0 opacity-100 sm:translate-x-0" leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100" leave-to-class="opacity-0">
        <div v-for="notification in notifications" :key="notification.id" class="pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg bg-white/80 backdrop-blur-md shadow-lg ring-1 ring-black ring-opacity-5 border border-white/30">
          <div class="p-4">
            <div class="flex items-start">
              <div class="flex-shrink-0">
                <svg v-if="notification.type === 'success'" class="h-6 w-6 text-green-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                <svg v-if="notification.type === 'error'" class="h-6 w-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" /></svg>
              </div>
              <div class="ml-3 flex-1">
                <p class="text-sm font-medium text-gray-900">{{ notification.title }}</p>
                <p v-if="notification.description" class="mt-1 text-sm text-gray-500">{{ notification.description }}</p>
              </div>
              <div class="ml-4 flex flex-shrink-0">
                <button @click="removeNotification(notification.id)" type="button" class="inline-flex rounded-md bg-white/0 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                  <span class="sr-only">Close</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.414 1.414L10 11.06l3.72 3.72a.75.75 0 101.414-1.414L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" /></svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </transition-group>
    </div>
  </div>

  <AuthCallback v-if="showAuthCallback" @login-success="handleLoginSuccess" />
  
  <LoginPage v-else-if="!isAuthenticated" />

  <div v-else class="min-h-screen w-full flex items-center justify-center p-4 bg-gray-900/10">
    <div class="w-full max-w-screen-2xl h-[90vh] bg-white/60 backdrop-blur-xl rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-white/30 relative">
      
      <SettingsPage v-if="showSettingsPage" :current-user="currentUser" @close="showSettingsPage = false" @notify="(n) => addNotification(n.title, n.description, n.type)" />
      
      <template v-else>
        <div v-if="isLoading" class="absolute inset-0 bg-gray-900/30 backdrop-blur-sm flex items-center justify-center z-50 rounded-2xl">
          <div class="flex flex-col items-center">
            <svg class="animate-spin h-12 w-12 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            <span class="mt-4 text-white text-lg font-semibold">{{ loadingText }}</span>
          </div>
        </div>
        
        <header class="flex-shrink-0 flex items-center justify-between border-b border-white/30 shadow-sm z-10 p-4">
          <h1 class="text-xl font-bold text-gray-800">{{ getProviderDisplayName(currentUser?.provider) }}报告助手</h1>
          
          <div class="flex items-center space-x-3">
            <div class="relative" ref="profileMenuNode">
              <button @click="isProfileMenuOpen = !isProfileMenuOpen" class="flex items-center space-x-2 p-2 rounded-lg transition-colors duration-200" :class="[isProfileMenuOpen ? 'bg-gray-500/20' : 'hover:bg-gray-500/10']">
                <img v-if="currentUser?.avatar_url" :src="currentUser.avatar_url" :alt="currentUser.name" class="h-8 w-8 rounded-full object-cover border border-gray-200">
                <span v-else class="inline-flex items-center justify-center h-8 w-8 rounded-full text-white font-bold text-sm" :class="getProviderAvatarClass(currentUser?.provider)">{{ currentUser?.name?.charAt(0)?.toUpperCase() || 'U' }}</span>
                <div class="text-left">
                  <div class="text-sm font-semibold text-gray-700">{{ currentUser?.name || '用户' }}</div>
                  <div class="text-xs text-gray-500">{{ getProviderDisplayName(currentUser?.provider) }}用户</div>
                </div>
                <svg class="w-4 h-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" /></svg>
              </button>
              
              <transition enter-active-class="transition ease-out duration-200" enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
                <div v-if="isProfileMenuOpen" class="absolute right-0 mt-2 w-64 bg-white/90 backdrop-blur-sm rounded-lg shadow-lg py-2 z-20 border border-white/30">
                  <div class="px-4 py-3 border-b border-gray-200/50">
                    <div class="flex items-center space-x-3 min-w-0">
                      <img v-if="currentUser?.avatar_url" :src="currentUser.avatar_url" :alt="currentUser.name" class="h-10 w-10 rounded-full object-cover border border-gray-200">
                      <span v-else class="inline-flex items-center justify-center h-10 w-10 rounded-full text-white font-bold" :class="getProviderAvatarClass(currentUser?.provider)">
                        {{ currentUser?.name?.charAt(0)?.toUpperCase() || 'U' }}
                      </span>
                      <div class="flex-1 min-w-0">
                        <div class="text-sm font-semibold text-gray-900 truncate">{{ currentUser?.name || '用户' }}</div>
                        <div class="text-xs text-gray-500 break-all max-w-[180px]" :title="currentUser?.email || '--'">
                          {{ currentUser?.email || '--' }}
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="py-1">
                    <a href="#" @click.prevent="showSettingsPage = true; isProfileMenuOpen = false" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-white/50 transition-colors">
                      <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826 3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
                      <span>设置</span>
                    </a>
                  </div>
                  <div class="border-t border-gray-200/50"></div>
                  <div class="py-1">
                    <a href="#" @click.prevent="handleLogout" class="flex items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors">
                      <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
                      <span>退出登录</span>
                    </a>
                  </div>
                </div>
              </transition>
            </div>
          </div>
        </header>
        
        <main class="flex-grow flex flex-col overflow-hidden">
          <section class="flex-shrink-0 border-b border-white/30 p-4">
              <div class="grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
                  <div>
                      <label for="source_template" class="text-sm font-medium text-gray-700">源模板:</label>
                      <select id="source_template" v-model="selectedSourceTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                          <option v-for="template in templates" :key="template.id" :value="template.id">{{ template.name }}</option>
                      </select>
                  </div>
                  <div>
                    <label for="date" class="text-sm font-medium text-gray-700">报告周期:</label>
                    <VueTailwindDatepicker id="date" v-model="dateValue" i18n="zh-cn" placeholder="选择日期范围" :shortcuts="datePickerShortcuts" :options="datePickerOptions" use-range :formatter="{ date: 'YYYY-MM-DD', month: 'MMM' }" class="mt-1 w-full" />
                  </div>
                  <div>
                      <label for="template" class="text-sm font-medium text-gray-700">目标模板</label>
                      <select id="template" v-model="selectedTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                          <option v-for="template in templates" :key="template.id" :value="template.id">{{ template.name }}</option>
                      </select>
                  </div>
                  <div class="md:col-span-2 flex space-x-2 items-end">
                    <button @click="getReports" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-white/30 text-gray-800 backdrop-blur-sm border border-white/40 shadow-lg hover:bg-white/50 hover:text-gray-900 focus:outline-none">获取报告</button>
                    <button @click="generateDraft" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-indigo-500/20 text-indigo-800 backdrop-blur-sm border border-indigo-500/30 shadow-lg hover:bg-indigo-500/40 hover:text-indigo-900 focus:outline-none">生成草稿</button>
                  </div>
              </div>
          </section>
          
          <div class="flex-grow grid grid-cols-1 md:grid-cols-2 gap-4 p-4 overflow-hidden">
              <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md flex flex-col overflow-hidden border border-white/20">
                  <h2 class="text-lg font-semibold p-4 border-b border-white/20">报告内容</h2>
                  <div class="overflow-y-auto flex-grow p-4 space-y-2">
                      <div v-if="sourceReports.length === 0" class="text-gray-500 text-center pt-10">报告内容将在此处显示。</div>
                      <div v-for="report in sourceReports" :key="report.id" class="border rounded-md">
                          <div @click="toggleReportDetail(report)" class="p-3 flex justify-between items-center cursor-pointer hover:bg-gray-50">
                              <h3 class="font-semibold">{{ report.title }}</h3>
                              <svg class="w-5 h-5 transition-transform" :class="{'rotate-180': !report.isCollapsed}" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"><path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"/></svg>
                          </div>
                          <div v-if="!report.isCollapsed" class="p-4 border-t border-white/20 prose max-w-none text-sm">
                            <TiptapViewer :content="report.fields.map(f => `<h3>${f.name}</h3><div>${f.value}</div>`).join('')" />
                          </div>
                      </div>
                  </div>
              </section>

              <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md flex flex-col overflow-hidden border border-white/20">
                   <h2 class="text-lg font-semibold p-4 border-b border-white/20">{{ currentTemplate?.name || '生成的草稿' }}</h2>
                   <div class="flex-grow overflow-y-auto p-4 space-y-4">
                      <div v-if="currentTemplate" v-for="field in (currentTemplate.fields || [])" :key="field.id" class="space-y-2">
                          <label :for="field.id" class="font-semibold text-gray-700">{{ field.label }}</label>
                           
                           <!-- Rich Text -->
                           <TiptapEditor v-if="field.type === 'tiptap'" v-model="formValues[field.id]" :placeholder="field.placeholder" @showApiKeyConfig="showSettingsPage = true" />
                           
                           <!-- Number -->
                           <input v-else-if="field.type === 'number'" :id="field.id" type="number" v-model.number="formValues[field.id]" :placeholder="field.placeholder" class="form-input" />
 
                           <!-- Dropdown -->
                           <select v-else-if="field.type === 'dropdown'" :id="field.id" v-model="formValues[field.id]" class="form-input">
                               <option disabled value="">请选择</option>
                               <option v-for="opt in (field.options || [])" :key="opt.value" :value="opt.value">{{ opt.text }}</option>
                           </select>
 
                           <!-- Multi-Select Checkboxes -->
                           <div v-else-if="field.type === 'multiSelect'" class="flex flex-wrap gap-x-4 gap-y-2 pt-1">
                              <label v-for="opt in (field.options || [])" :key="opt.value" class="flex items-center space-x-2 cursor-pointer">
                                <input type="checkbox" :value="opt.value" v-model="formValues[field.id]" class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500">
                                <span>{{ opt.text }}</span>
                              </label>
                           </div>
                           
                           <!-- Address -->
                           <textarea v-else-if="field.type === 'address'" :id="field.id" v-model="formValues[field.id]" :placeholder="field.placeholder" rows="3" class="form-input"></textarea>
 
                           <!-- DateTime -->
                           <input v-else-if="field.type === 'datetime'" :id="field.id" type="datetime-local" v-model="formValues[field.id]" class="form-input" />
                           
                           <!-- Image Uploader -->
                           <div v-else-if="field.type === 'image'">
                              <div class="grid grid-cols-5 gap-4">
                                 <div v-for="file in (formValues[field.id] || [])" :key="file.id" class="relative w-24 h-24">
                                   <img :src="file.url" :alt="file.name" class="w-full h-full object-cover rounded-lg shadow-md">
                                   <button @click="removeFile(field.id, file.id)" class="absolute -top-1 -right-1 bg-gray-700 text-white rounded-full h-5 w-5 flex items-center justify-center text-xs">&times;</button>
                                 </div>
                                 <div @click="triggerFileInput(field.id)" class="w-24 h-24 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center cursor-pointer hover:bg-gray-50/50">
                                   <span class="text-3xl text-gray-400">+</span>
                                 </div>
                              </div>
                              <p class="text-xs text-gray-500 mt-2">单张图片不超过 {{ formatFileSize(field.maxSize) }}，最多上传 {{ field.maxCount }} 张</p>
                              <input type="file" multiple accept="image/*" class="hidden" :ref="(el) => fileInputRefs[field.id] = el" @change="handleFileSelect($event, field)">
                           </div>
 
                           <!-- Attachment Uploader -->
                           <div v-else-if="field.type === 'attachment'">
                              <button @click="triggerFileInput(field.id)" class="px-4 py-2 text-sm font-semibold border border-gray-300 rounded-md hover:bg-gray-50/50">+ 添加附件</button>
                              <p class="text-xs text-gray-500 mt-2">单个附件不超过 {{ formatFileSize(field.maxSize) }}，最多上传 {{ field.maxCount }} 个</p>
                              <div class="mt-4 space-y-2">
                                 <div v-for="file in (formValues[field.id] || [])" :key="file.id" class="flex items-center justify-between p-2 bg-gray-100/80 rounded-md">
                                    <div class="flex items-center space-x-2">
                                       <svg class="h-6 w-6 text-yellow-500" fill="currentColor" viewBox="0 0 20 20"><path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2-2H4a2 2 0 01-2-2V6z"></path></svg>
                                       <div>
                                         <p class="text-sm font-medium">{{ file.name }}</p>
                                         <p class="text-xs text-gray-500">{{ formatFileSize(file.size) }}</p>
                                       </div>
                                    </div>
                                    <button @click="removeFile(field.id, file.id)" class="p-1 text-gray-500 hover:text-red-600">
                                       <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>
                                    </button>
                                 </div>
                              </div>
                              <input type="file" multiple class="hidden" :ref="(el) => fileInputRefs[field.id] = el" @change="handleFileSelect($event, field)">
                           </div>
 
                           <!-- User Picker (Placeholder) -->
                           <div v-else-if="field.type === 'user-picker'">
                               <div class="form-input flex items-center bg-gray-100/80 cursor-not-allowed">
                                   <div class="flex items-center space-x-2 text-gray-500">
                                       <svg class="h-6 w-6 text-gray-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"></path></svg>
                                       <span>选择人员 (暂不可用)</span>
                                   </div>
                               </div>
                           </div>

                           <!-- Plain Text -->
                           <input v-else :id="field.id" type="text" v-model="formValues[field.id]" :placeholder="field.placeholder" class="form-input" />
                      </div>
                   </div>
              </section>
          </div>
        </main>
      </template>
    </div>
  </div>
</template>

<style>
/* Add any necessary global styles here */
.form-input {
    @apply mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80;
}
</style> 