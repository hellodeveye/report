<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import TiptapEditor from './components/TiptapEditor.vue';
import TiptapViewer from './components/TiptapViewer.vue';
import VueTailwindDatepicker from 'vue-tailwind-datepicker';
import LoginPage from './components/LoginPage.vue';
import AuthCallback from './components/AuthCallback.vue';
import SettingsPage from './components/SettingsPage.vue';
import { feishuApiService, reportSummarizer, aiService } from './utils/aiService.js';
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
const fileInputRefs = ref({});
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
  
  // 添加键盘事件监听器
  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
      if (showSettingsPage.value) {
        showSettingsPage.value = false;
      }
    }
  });
  
  // 检查认证状态
  await checkAuthStatus();
});

// 检查认证状态
const checkAuthStatus = async () => {
  try {
    // 检查URL是否包含飞书OAuth回调
    if (window.location.pathname === '/auth/callback' || window.location.search.includes('code=')) {
      showAuthCallback.value = true;
      return; // 不继续执行其他逻辑
    }
    
    // 检查是否已登录
    if (authService.isAuthenticated()) {
      try {
        // 优先从本地存储获取用户信息
        let user = authService.getUser();
        
        // 如果本地没有用户信息，从服务器获取
        if (!user) {
          user = await authService.getCurrentUser();
        }
        
        currentUser.value = user;
        isAuthenticated.value = true;
        
        // 加载模板列表
        await loadTemplates();
        
        // 只在首次登录时显示欢迎消息
        if (!sessionStorage.getItem('login_welcomed')) {
          addNotification('登录成功', `欢迎回来，${user.name || '用户'}！`, 'success');
          sessionStorage.setItem('login_welcomed', 'true');
        }
      } catch (error) {
        console.error('获取用户信息失败:', error);
        authService.clearToken();
        isAuthenticated.value = false;
        currentUser.value = null;
      }
    } else {
      isAuthenticated.value = false;
      currentUser.value = null;
    }
  } catch (error) {
    console.error('认证状态检查失败:', error);
    isAuthenticated.value = false;
    currentUser.value = null;
  }
};

// 退出登录
const handleLogout = async () => {
  try {
    // 清除欢迎标记
    sessionStorage.removeItem('login_welcomed');
    // 调用authService的logout方法（会自动清理状态并跳转）
    await authService.logout();
  } catch (error) {
    console.error('退出登录失败:', error);
    // 即使API调用失败，也清除本地状态
    authService.clearToken();
    isAuthenticated.value = false;
    currentUser.value = null;
    // 手动跳转到登录页
    window.location.href = '/';
  }
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

// 获取平台头像背景色
const getProviderAvatarClass = (provider) => {
  switch (provider) {
    case 'feishu':
      return 'bg-indigo-500';
    case 'dingtalk':
      return 'bg-blue-500';
    default:
      return 'bg-gray-500';
  }
};

// 获取平台默认邮箱后缀
const getProviderDefaultEmail = (provider) => {
  switch (provider) {
    case 'feishu':
      return 'feishu@user';
    case 'dingtalk':
      return 'dingtalk@user';
    default:
      return 'unknown@user';
  }
};

// 加载模板列表
const loadTemplates = async () => {
  try {
    loadingText.value = '正在加载模板...';
    isLoading.value = true;
    const templatesData = await feishuApiService.getAllTemplates();
    templates.value = templatesData;
    
    // 设置默认选中的模板
    if (templatesData.length > 0) {
      selectedSourceTemplateId.value = templatesData[0].id;
      selectedTemplateId.value = templatesData[0].id;
    }
    
    console.log(`成功加载${templatesData.length}个模板`);
  } catch (error) {
    console.error('加载模板失败:', error);
    addNotification('模板加载失败', error.message, 'error');
    
    // 如果加载失败，使用固定的模板列表
    const fixedTemplates = feishuApiService.getFixedTemplateList();
    const fallbackTemplates = fixedTemplates.map(t => 
      feishuApiService.getDefaultTemplate(t.id, t.name)
    );
    templates.value = fallbackTemplates;
    
    if (fallbackTemplates.length > 0) {
      selectedSourceTemplateId.value = fallbackTemplates[0].id;
      selectedTemplateId.value = fallbackTemplates[0].id;
    }
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
        const dayOfWeek = today.getDay();
        const diffToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
        const thisWeekStart = new Date(today);
        thisWeekStart.setDate(today.getDate() - diffToMonday);
        return [thisWeekStart, today];
      },
    },
    {
      label: "上周",
      atClick: () => {
        const today = new Date();
        const dayOfWeek = today.getDay();
        const diffToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
        const thisWeekStart = new Date(today);
        thisWeekStart.setDate(today.getDate() - diffToMonday);
        const lastWeekEnd = new Date(thisWeekStart);
        lastWeekEnd.setDate(thisWeekStart.getDate() - 1);
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
    past: (period) => `过去${period}天`,
    currentMonth: "本月",
    pastMonth: "上月",
  },
  footer: {
    apply: "应用",
    cancel: "取消",
  },
});

const currentTemplate = computed(() => {
  return templates.value.find(t => t.id === selectedTemplateId.value);
});

const sourceReports = ref([]);
const generatedForm = ref([]); // This will hold the structured form fields for the right panel

// --- Methods ---

const toggleReportDetail = (report) => {
    report.isCollapsed = !report.isCollapsed;
};

const getReports = async () => {
  try {
    loadingText.value = '正在获取报告...';
    isLoading.value = true;
    
    // 构建查询参数
    const params = {};
    
    // 如果选择了源模板，添加模板过滤
    if (selectedSourceTemplateId.value) {
      params.rule_id = selectedSourceTemplateId.value;
    }
    
    // 如果选择了日期范围，添加时间过滤
    if (dateValue.value.startDate && dateValue.value.endDate) {
      // 开始时间：当天的开始（00:00:00）
      const startDate = new Date(dateValue.value.startDate);
      params.start_time = Math.floor(startDate.getTime() / 1000);
      
      // 结束时间：当天的最晚时间（23:59:59）
      const endDate = new Date(dateValue.value.endDate);
      endDate.setHours(23, 59, 59, 999); // 设置为当天的最晚时间
      params.end_time = Math.floor(endDate.getTime() / 1000);
    }
    
    // 添加详细的调试日志
    if (params.start_time && params.end_time) {
      console.log('调用 /reports 接口，参数:', params);
      console.log('时间范围:', {
        startDate: dateValue.value.startDate,
        endDate: dateValue.value.endDate,
        startTime: new Date(params.start_time * 1000).toLocaleString(),
        endTime: new Date(params.end_time * 1000).toLocaleString()
      });
    } else {
      console.log('调用 /reports 接口，参数:', params);
    }
    
    // 获取原始模板数据用于字段类型映射
    const templateData = selectedSourceTemplateId.value ? 
      feishuApiService.getRawRuleById(selectedSourceTemplateId.value) : null;
    
    const reportsData = await feishuApiService.getReports(params, templateData);
    sourceReports.value = reportsData;
    
    addNotification('报告获取成功', `成功获取${reportsData.length}条报告`, 'success');
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
  
  // 检查 API Key
  if (!aiService.hasApiKey()) {
    addNotification(
      '需要配置AI模型',
      '请先前往“设置”页面配置您的AI提供商和API Key。',
      'error',
      5000
    );
    return;
  }
  
  loadingText.value = '正在生成草稿...';
  isLoading.value = true;
  
  // 检查是否有左侧报告数据
  const hasSourceReports = sourceReports.value && sourceReports.value.length > 0;
  
  try {
    if (hasSourceReports) {
      
      // 使用AI汇总生成草稿
      loadingText.value = '正在分析源报告并生成智能草稿...';
      addNotification('正在生成', '正在分析源报告并生成智能草稿...', 'success');
      
      const summary = await reportSummarizer.summarizeReports(
        sourceReports.value,
        currentTemplate.value
      );
      
      formValues.value = summary;
      addNotification('AI草稿已生成', '已基于左侧报告内容智能生成草稿', 'success');
    } else {
      // 没有源报告时，使用AI生成通用草稿
      loadingText.value = '正在生成AI智能草稿...';
      addNotification('正在生成', '正在生成AI智能草稿...', 'success');
      
      const newValues = {};
      
      // 为每个字段生成AI内容
      for (const field of currentTemplate.value.fields) {
        try {
          if (field.type === 'tiptap' || field.type === 'text' || field.type === 'address') {
            // 对文本类型字段使用AI生成
            const prompt = `请为"${field.label}"字段生成合适的${currentTemplate.value.name}内容。要求：
1. 内容专业且实用
2. 符合工作报告的语气
3. 字数控制在50-200字之间
4. ${field.type === 'tiptap' ? '使用HTML格式' : '纯文本格式'}`;
            
            const content = await aiService.streamProcess(
              prompt,
              `字段名称：${field.label}`,
              null,
              { stream: false }
            );
            
            newValues[field.id] = content || `请填写${field.label}`;
          } else if (field.type === 'number') {
            newValues[field.id] = Math.floor(Math.random() * 5) + 1;
          } else if (field.type === 'dropdown') {
            if (field.options && field.options.length > 0) {
              newValues[field.id] = field.options[Math.floor(Math.random() * field.options.length)].value;
            } else {
              newValues[field.id] = '';
            }
          } else if (field.type === 'multiSelect') {
            if (field.options && field.options.length > 0) {
              newValues[field.id] = field.options
                .filter(() => Math.random() > 0.5)
                .map(opt => opt.value);
              if (newValues[field.id].length === 0) {
                newValues[field.id].push(field.options[0].value);
              }
            } else {
              newValues[field.id] = [];
            }
          } else if (field.type === 'datetime') {
            newValues[field.id] = new Date(Date.now() + Math.random() * 1000 * 3600 * 24 * 7).toISOString().substring(0, 16);
          } else if (field.type === 'image' || field.type === 'attachment') {
            newValues[field.id] = [{
              name: 'placeholder.png',
              size: 818200,
              url: 'https://template.tiptap.dev/images/placeholder-image.png',
              id: Date.now()
            }];
          } else {
            newValues[field.id] = '';
          }
        } catch (error) {
          console.warn(`生成字段 ${field.label} 失败:`, error);
          newValues[field.id] = `请填写${field.label}`;
        }
      }
      
      formValues.value = newValues;
      addNotification('AI草稿已生成', '已生成AI智能草稿，您可以根据需要进行修改', 'success');
    }
  } catch (error) {
    console.error('生成草稿失败:', error);
    addNotification('生成失败', error.message, 'error');
  } finally {
    isLoading.value = false;
  }
};

// Initialize form values when component mounts or template changes
const initializeFormValues = () => {
  const values = {};
  if (currentTemplate.value) {
    currentTemplate.value.fields.forEach(field => {
      // Set default value based on type
      if (field.type === 'multiSelect' || field.type === 'image' || field.type === 'attachment') {
        values[field.id] = [];
      } else {
        values[field.id] = '';
      }
    });
  }
  formValues.value = values;
};

watch(selectedTemplateId, initializeFormValues, { immediate: true });

const triggerFileInput = (fieldId) => {
  fileInputRefs.value[fieldId]?.click();
};

const handleFileSelect = (event, field) => {
  const files = event.target.files;
  if (!files) return;

  const currentFiles = formValues.value[field.id] || [];
  if (currentFiles.length + files.length > field.maxCount) {
    alert(`最多只能上传 ${field.maxCount} 个文件。`);
    return;
  }

  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    if (file.size > field.maxSize) {
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
  // Reset file input to allow selecting the same file again
  event.target.value = '';
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
}

const handleShowApiKeyConfig = () => {
  addNotification(
    '需要配置AI模型',
    '请先前往“设置”页面配置您的AI提供商和API Key。',
    'error',
    5000
  );
};

</script>

<template>
  <!-- Notifications container -->
  <div aria-live="assertive" class="pointer-events-none fixed inset-0 flex items-end px-4 py-6 sm:items-start sm:p-6 z-50">
    <div class="flex w-full flex-col items-center sm:items-end">
      <transition-group
        name="notification"
        tag="div"
        class="w-full space-y-4 flex flex-col items-center sm:items-end"
        enter-active-class="transform ease-out duration-300 transition"
        enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
        enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
        leave-active-class="transition ease-in duration-100"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
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

  <!-- OAuth回调处理页面 -->
  <AuthCallback v-if="showAuthCallback" />
  
  <!-- 登录页面 -->
  <LoginPage v-else-if="!isAuthenticated" />
  
  <!-- 主应用界面 -->
  <div v-else class="min-h-screen w-full flex items-center justify-center p-4 bg-gray-900/10">
    <div class="w-full max-w-screen-2xl h-[90vh] bg-white/60 backdrop-blur-xl rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-white/30 relative">
      
      <!-- Settings Page -->
      <SettingsPage
        v-if="showSettingsPage"
        :current-user="currentUser"
        @close="showSettingsPage = false"
        @notify="(notification) => addNotification(notification.title, notification.description, notification.type)"
      />
      
      <!-- Main Content -->
      <template v-else>
        <!-- Loading Overlay -->
        <div v-if="isLoading" class="absolute inset-0 bg-gray-900/30 backdrop-blur-sm flex items-center justify-center z-50 rounded-2xl">
          <div class="flex flex-col items-center bg-white/90 backdrop-blur-md rounded-2xl p-8 shadow-xl border border-white/30">
            <svg class="animate-spin h-12 w-12 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span class="mt-4 text-gray-800 text-lg font-semibold">{{ loadingText }}</span>
          </div>
        </div>
        
        <header class="flex-shrink-0 flex items-center justify-between border-b border-white/30 shadow-sm z-10 p-4">
          <h1 class="text-xl font-bold text-gray-800">{{ getProviderDisplayName(currentUser?.provider) }}报告助手</h1>
          
          <!-- User Profile Section -->
          <div class="flex items-center space-x-3">
            <div class="relative" ref="profileMenuNode">
              <button @click="isProfileMenuOpen = !isProfileMenuOpen" 
                      class="flex items-center space-x-2 p-2 rounded-lg transition-colors duration-200"
                      :class="[isProfileMenuOpen ? 'bg-gray-500/20' : 'hover:bg-gray-500/10']">
              <!-- 用户头像 -->
              <img v-if="currentUser?.avatar_url" 
                   :src="currentUser.avatar_url" 
                   :alt="currentUser.name"
                   class="h-8 w-8 rounded-full object-cover border border-gray-200">
              <span v-else 
                    class="inline-flex items-center justify-center h-8 w-8 rounded-full text-white font-bold text-sm"
                    :class="getProviderAvatarClass(currentUser?.provider)">
                {{ currentUser?.name?.charAt(0)?.toUpperCase() || 'U' }}
              </span>
              <div class="text-left">
                <div class="text-sm font-semibold text-gray-700">{{ currentUser?.name || '用户' }}</div>
                <div class="text-xs text-gray-500">{{ getProviderDisplayName(currentUser?.provider) }}用户</div>
              </div>
              <svg class="w-4 h-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
            
            <!-- Dropdown menu -->
            <transition
              enter-active-class="transition ease-out duration-200"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div v-if="isProfileMenuOpen" class="absolute right-0 mt-2 w-64 bg-white/90 backdrop-blur-sm rounded-lg shadow-lg py-2 z-20 border border-white/30">
                <!-- 用户信息区域 -->
                                  <div class="px-4 py-3 border-b border-gray-200/50">
                    <div class="flex items-center space-x-3 min-w-0">
                    <!-- 用户头像 -->
                    <img v-if="currentUser?.avatar_url" 
                         :src="currentUser.avatar_url" 
                         :alt="currentUser.name"
                         class="h-10 w-10 rounded-full object-cover border border-gray-200">
                    <span v-else 
                          class="inline-flex items-center justify-center h-10 w-10 rounded-full text-white font-bold"
                          :class="getProviderAvatarClass(currentUser?.provider)">
                      {{ currentUser?.name?.charAt(0)?.toUpperCase() || 'U' }}
                    </span>
                    <div class="flex-1 min-w-0">
                      <div class="text-sm font-semibold text-gray-900 truncate">{{ currentUser?.name || '用户' }}</div>
                      <div class="text-xs text-gray-500 break-all max-w-[180px]" :title="currentUser?.email || getProviderDefaultEmail(currentUser?.provider)">
                        {{ currentUser?.email || getProviderDefaultEmail(currentUser?.provider) }}
                      </div>
                    </div>
                  </div>
                </div>
                
                <!-- 菜单项 -->
                <div class="py-1">
                  <a href="#" @click.prevent="showSettingsPage = true; isProfileMenuOpen = false" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-white/50 transition-colors">
                    <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826 3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    </svg>
                    <span>设置</span>
                  </a>
                </div>
                
                <div class="border-t border-gray-200/50"></div>
                <div class="py-1">
                  <a href="#" @click.prevent="handleLogout" class="flex items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors">
                    <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                    </svg>
                    <span>退出登录</span>
                  </a>
                </div>
              </div>
            </transition>
          </div>
          </div>
        </header>
        
        <main class="flex-grow flex flex-col overflow-hidden">
          <!-- Top Controls Section -->
          <section class="flex-shrink-0 border-b border-white/30 p-4">
              <div class="grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
                  <div>
                      <label for="source_template" class="text-sm font-medium text-gray-700">源模板:</label>
                      <select id="source_template" v-model="selectedSourceTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                          <option v-for="template in templates" :key="template.id" :value="template.id">
                              {{ template.name }}
                          </option>
                      </select>
                  </div>
                  <div>
                    <label for="date" class="text-sm font-medium text-gray-700">报告周期:</label>
                    <VueTailwindDatepicker
                        id="date"
                        v-model="dateValue"
                        :formatter="{ date: 'YYYY-MM-DD', month: 'MMM' }"
                        i18n="zh-cn"
                        placeholder="选择日期范围"
                        :shortcuts="datePickerShortcuts"
                        :options="datePickerOptions"
                        :auto-apply="true"
                        class="mt-1 w-full"
                    />
                  </div>
                  <div>
                      <label for="template" class="text-sm font-medium text-gray-700">目标模板</label>
                      <select id="template" v-model="selectedTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                          <option v-for="template in templates" :key="template.id" :value="template.id">
                              {{ template.name }}
                          </option>
                      </select>
                  </div>
                  <div class="md:col-span-2 flex space-x-2 items-end">
                    <button @click="getReports" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-white/30 text-gray-800 backdrop-blur-sm border border-white/40 shadow-lg hover:bg-white/50 hover:text-gray-900 focus:outline-none">
                      获取报告
                    </button>
                    <button @click="generateDraft" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-indigo-500/20 text-indigo-800 backdrop-blur-sm border border-indigo-500/30 shadow-lg hover:bg-indigo-500/40 hover:text-indigo-900 focus:outline-none">
                      生成草稿
                    </button>
                    <!-- <button @click="addNotification('发送成功', '报告已提交，请在飞书中查收。')" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-blue-500/20 text-blue-800 backdrop-blur-sm border border-blue-500/30 shadow-lg hover:bg-blue-500/40 hover:text-blue-900 focus:outline-none">
                      发送到飞书
                    </button> -->
                  </div>
              </div>
          </section>
          
          <div class="flex-grow grid grid-cols-1 md:grid-cols-2 gap-4 p-4 overflow-hidden">
              <!-- Left Column: Source Reports -->
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
                            <div v-for="field in report.fields" :key="field.name" class="mb-3">
                                <b class="text-gray-700">{{ field.name }}:</b>
                                
                                <div v-if="field.type === 'image'" class="mt-2 grid grid-cols-4 gap-2">
                                    <div v-for="file in (field.value || [])" :key="file.name" class="relative">
                                      <img :src="file.url" alt="图片预览" class="w-full h-auto rounded-md shadow-md border border-gray-200" />
                                    </div>
                                </div>
                                <div v-else-if="field.type === 'attachment'" class="mt-1 space-y-2">
                                    <div v-for="file in (field.value || [])" :key="file.name">
                                        <a :href="file.url" target="_blank" class="text-indigo-600 hover:underline">{{ file.name }} ({{ formatFileSize(file.size) }})</a>
                                    </div>
                                </div>
                                <div v-else-if="field.type === 'multiSelect'" class="inline-flex flex-wrap gap-2 mt-1">
                                  <span v-for="item in (field.value || [])" :key="item" class="bg-gray-200 text-gray-700 px-2 py-0.5 rounded-full text-xs">{{ item }}</span>
                                </div>
                                <div v-else-if="field.type === 'tiptap'" class="mt-1">
                                  <TiptapViewer :content="field.value" />
                                </div>
                                <span v-else class="ml-2 text-gray-800">{{ field.value }}</span>
                            </div>
                          </div>
                      </div>
                  </div>
              </section>

              <!-- Right Column: Generated Form -->
              <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md flex flex-col overflow-hidden border border-white/20">
                   <h2 class="text-lg font-semibold p-4 border-b border-white/20">{{ currentTemplate?.name || '生成的草稿' }}</h2>
                   <div class="flex-grow overflow-y-auto p-4 space-y-4">
                      <div v-if="currentTemplate" v-for="field in (currentTemplate.fields || [])" :key="field.id" class="space-y-2">
                          <label :for="field.id" class="font-semibold text-gray-700">{{ field.label }}</label>
                          
                          <!-- Rich Text -->
                          <TiptapEditor v-if="field.type === 'tiptap'" 
                                        v-model="formValues[field.id]" 
                                        :placeholder="field.placeholder" 
                                        @showApiKeyConfig="handleShowApiKeyConfig" />
                          
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
                                      <svg class="h-6 w-6 text-yellow-500" fill="currentColor" viewBox="0 0 20 20"><path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"></path></svg>
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
  <!-- 主应用界面结束 -->
</template>

<style>
.prose:focus {
  outline: none;
}
.tiptap {
  height: 100%;
  background-color: rgba(255, 255, 255, 0.8);
}
.tiptap:focus {
  outline: none;
}

/* Common form input style */
.form-input {
  @apply w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80;
}

/* Style for Vue-Tailwind-Datepicker input */
[data-headlessui-state="open"] ~ div, input {
  background-color: rgba(255, 255, 255, 0.8) !important;
}

/* Custom override for Today's date border */
.vtd-today:not(.vtd-start-date):not(.vtd-end-date) span {
  border: 1px solid #F97316 !important; /* orange-500 */
}

/* Custom scrollbar styles */
* {
  scrollbar-width: thin;
  scrollbar-color: #CBD5E1 #F1F5F9;
}

*::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

*::-webkit-scrollbar-track {
  background: #F1F5F9;
  border-radius: 4px;
}

*::-webkit-scrollbar-thumb {
  background: #CBD5E1;
  border-radius: 4px;
  border: 1px solid #F1F5F9;
}

*::-webkit-scrollbar-thumb:hover {
  background: #94A3B8;
}

*::-webkit-scrollbar-corner {
  background: #F1F5F9;
}
</style> 