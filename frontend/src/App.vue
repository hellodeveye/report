<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import TiptapEditor from './components/TiptapEditor.vue';
import TiptapViewer from './components/TiptapViewer.vue';
import VueTailwindDatepicker from 'vue-tailwind-datepicker';
import LoginPage from './components/LoginPage.vue';
import AuthCallback from './components/AuthCallback.vue';
import SettingsPage from './components/SettingsPage.vue';
import TurndownService from 'turndown';
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
const initialSettingsTab = ref('account');

const templates = ref([]);
const selectedSourceTemplateId = ref('');
const selectedTemplateId = ref('');
const formValues = ref({});

const getThisWeekRange = () => {
    const today = new Date();
    const dayOfWeek = today.getDay() || 7; // Sunday is 7, Monday is 1
    const thisWeekStart = new Date(today);
    thisWeekStart.setDate(today.getDate() - dayOfWeek + 1);
    
    const formatDate = (date) => {
        const d = new Date(date);
        let month = '' + (d.getMonth() + 1);
        let day = '' + d.getDate();
        const year = d.getFullYear();

        if (month.length < 2) 
            month = '0' + month;
        if (day.length < 2) 
            day = '0' + day;

        return [year, month, day].join('-');
    }

    return {
        startDate: formatDate(thisWeekStart),
        endDate: formatDate(today)
    };
};

const dateValue = ref(getThisWeekRange());

const isProfileMenuOpen = ref(false);
const profileMenuNode = ref(null);
const isActionMenuOpen = ref(false);
const actionMenuNode = ref(null);
const isContentMenuOpen = ref(false);
const contentMenuNode = ref(null);
const notifications = ref([]);

onMounted(async () => {
  document.addEventListener('click', (event) => {
    if (profileMenuNode.value && !profileMenuNode.value.contains(event.target)) {
      isProfileMenuOpen.value = false;
    }
    if (actionMenuNode.value && !actionMenuNode.value.contains(event.target)) {
      isActionMenuOpen.value = false;
    }
    if (contentMenuNode.value && !contentMenuNode.value.contains(event.target)) {
      isContentMenuOpen.value = false;
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

    detailedTemplates.sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'));

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

const saveDraft = () => {
  if (!currentTemplate.value) {
    addNotification('无法保存草稿', '请先选择一个模板。', 'error');
    return;
  }
  try {
    const draftData = JSON.parse(JSON.stringify(formValues.value)); // Deep copy to prevent side effects

    // Do not save file objects as they are not serializable
    // and cannot be restored from localStorage.
    if (currentTemplate.value) {
      currentTemplate.value.fields.forEach(field => {
        if (field.type === 'image' || field.type === 'attachment') {
          delete draftData[field.id];
        }
      });
    }

    const draft = {
      formValues: draftData,
      timestamp: new Date().toISOString()
    };
    localStorage.setItem(`draft_${currentTemplate.value.id}`, JSON.stringify(draft));
    addNotification('草稿已保存', '草稿已成功保存到本地。', 'success');
  } catch (e) {
    console.error('保存草稿失败:', e);
    addNotification('保存草稿失败', '无法保存草稿，可能是本地存储已满。', 'error');
  } finally {
    isActionMenuOpen.value = false;
  }
};

const exportContent = () => {
  if (sourceReports.value.length === 0) {
    addNotification('没有可导出的内容', '请先获取报告内容。', 'info');
    isContentMenuOpen.value = false;
    return;
  }

  try {
    const turndownService = new TurndownService({
      headingStyle: 'atx',
      codeBlockStyle: 'fenced'
    });
    
    let markdownContent = `# 报告内容\n\n`;

    sourceReports.value.forEach(report => {
      markdownContent += `## ${report.title}\n\n`;
      report.fields.forEach(field => {
        markdownContent += `### ${field.name}\n`;
        // Convert HTML value to Markdown
        const markdownValue = turndownService.turndown(field.value);
        markdownContent += `${markdownValue}\n\n`;
      });
    });

    const blob = new Blob([markdownContent], { type: 'text/markdown;charset=utf-8' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    
    const date = new Date().toISOString().slice(0, 10);
    link.download = `报告内容-${date}.md`;
    
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    
    addNotification('导出成功', 'Markdown文件已开始下载。', 'success');
  } catch (error) {
    console.error('导出Markdown失败:', error);
    addNotification('导出失败', '生成Markdown文件时发生错误。', 'error');
  } finally {
    isContentMenuOpen.value = false;
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

  // Try to load a saved draft
  if (currentTemplate.value) {
    const savedDraft = localStorage.getItem(`draft_${currentTemplate.value.id}`);
    if (savedDraft) {
      try {
        const draft = JSON.parse(savedDraft);
        // Merge draft values, this will not overwrite file fields as they were not saved.
        formValues.value = { ...formValues.value, ...draft.formValues };
        addNotification('草稿已加载', '已自动加载上次保存的草稿。', 'info');
      } catch (e) {
        console.error('加载草稿失败:', e);
        localStorage.removeItem(`draft_${currentTemplate.value.id}`); // Remove corrupted draft
      }
    }
  }
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

const sendReport = async () => {
  if (!currentTemplate.value) {
    addNotification('请选择一个模板', '', 'error');
    return;
  }

  // A simple validation to check that all fields are filled.
  for (const field of currentTemplate.value.fields) {
    const value = formValues.value[field.id];
    if (value === undefined || value === null || (Array.isArray(value) && value.length === 0) || String(value).trim() === '') {
        addNotification('请填写所有必填字段', `字段 "${field.label}" 不能为空。`, 'error');
        return;
    }
  }

  loadingText.value = '正在发送报告...';
  isLoading.value = true;
  try {
    const reportData = {
      template_id: currentTemplate.value.id,
      template_name: currentTemplate.value.name,
      contents: [],
    };

    for (const field of currentTemplate.value.fields) {
      reportData.contents.push({
        key: field.label, // Use label as key, which maps to field_name
        value: formValues.value[field.id],
      });
    }

    await apiService.sendDingTalkReport(reportData);
    addNotification('报告发送成功', '你的报告已成功提交', 'success');
    initializeFormValues(); // Optionally clear the form
  } catch (error) {
    console.error('发送报告失败:', error);
    addNotification('发送报告失败', error.message || 'An unknown error occurred', 'error');
  } finally {
    isLoading.value = false;
  }
};

const openSettings = () => {
  initialSettingsTab.value = 'model';
  showSettingsPage.value = true;
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
      
      <SettingsPage v-if="showSettingsPage" :current-user="currentUser" :initial-tab="initialSettingsTab" @close="showSettingsPage = false" @notify="(n) => addNotification(n.title, n.description, n.type)" />
      
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
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-3 text-gray-400">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                      </svg>
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
                      <label for="source_template" class="text-sm font-medium text-gray-700">日志模板:</label>
                      <select id="source_template" v-model="selectedSourceTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                          <option v-for="template in templates" :key="template.id" :value="template.id">{{ template.name }}</option>
                      </select>
                  </div>
                  <div>
                    <label for="date" class="text-sm font-medium text-gray-700">报告周期:</label>
                    <VueTailwindDatepicker id="date" v-model="dateValue" i18n="zh-cn" placeholder="选择日期范围" :shortcuts="datePickerShortcuts" :options="datePickerOptions" use-range :formatter="{ date: 'YYYY-MM-DD', month: 'MMM' }" class="mt-1 w-full" />
                  </div>
                  <div>
                      <label for="template" class="text-sm font-medium text-gray-700">生成模板</label>
                      <select id="template" v-model="selectedTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 bg-white/80">
                          <option v-for="template in templates" :key="template.id" :value="template.id">{{ template.name }}</option>
                      </select>
                  </div>
                  <div class="md:col-span-2 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
                    <button @click="getReports" class="w-full flex items-center justify-center px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-white/30 text-gray-800 backdrop-blur-sm border border-white/40 shadow-lg hover:bg-white/50 hover:text-gray-900 focus:outline-none">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5 mr-2">
                        <path fill-rule="evenodd" d="M10.5 3.75a6.75 6.75 0 1 0 0 13.5 6.75 6.75 0 0 0 0-13.5ZM2.25 10.5a8.25 8.25 0 1 1 14.59 5.28l4.69 4.69a.75.75 0 1 1-1.06 1.06l-4.69-4.69A8.25 8.25 0 0 1 2.25 10.5Z" clip-rule="evenodd" />
                      </svg>
                      获取报告
                    </button>
                    <button @click="generateDraft" class="w-full flex items-center justify-center px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-gradient-to-br from-purple-500 to-indigo-600 text-white backdrop-blur-sm border border-transparent shadow-lg hover:shadow-indigo-500/40 hover:from-purple-600 hover:to-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-2"><path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.898 20.572L16.5 21.75l-.398-1.178a3.375 3.375 0 00-2.456-2.456L12.5 18l1.178-.398a3.375 3.375 0 002.456-2.456L16.5 14.25l.398 1.178a3.375 3.375 0 002.456 2.456L20.25 18l-1.178.398a3.375 3.375 0 00-2.456 2.456z" /></svg>
                      生成草稿
                    </button>
                    <button v-if="currentUser?.provider === 'dingtalk'" @click="sendReport" class="w-full flex items-center justify-center px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-blue-500/20 text-blue-800 backdrop-blur-sm border border-blue-500/30 shadow-lg hover:bg-blue-500/40 hover:text-blue-900 focus:outline-none">
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" /></svg>
                      发送报告
                    </button>
                  </div>
              </div>
          </section>
          
          <div class="flex-grow md:grid md:grid-cols-2 gap-4 p-4 overflow-y-auto custom-scrollbar space-y-4 md:space-y-0">
              <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md border border-white/20 md:flex md:flex-col md:overflow-hidden">
                  <div class="flex items-center justify-between p-4 border-b border-white/20">
                    <h2 class="text-lg font-semibold">报告内容</h2>
                    <div class="relative flex-shrink-0" ref="contentMenuNode">
                      <button @click="isContentMenuOpen = !isContentMenuOpen" class="p-2 rounded-full text-gray-500 hover:bg-gray-500/10 hover:text-gray-800 focus:outline-none transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                          <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
                        </svg>
                      </button>
                      <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
                        <div v-if="isContentMenuOpen" class="absolute right-0 top-full mt-2 w-48 bg-white/90 backdrop-blur-sm rounded-md shadow-lg py-1 z-20 border border-white/30">
                          <a href="#" @click.prevent="exportContent" class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-500/10 transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-3 text-gray-500">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9.75v6.75m0 0l-3-3m3 3l3-3m-8.25 6a4.5 4.5 0 01-1.41-8.775 5.25 5.25 0 0110.233-2.33 3 3 0 013.758 3.848A3.752 3.752 0 0118 19.5H6.75z" />
                            </svg>
                            <span>导出</span>
                          </a>
                        </div>
                      </transition>
                    </div>
                  </div>
                  <div class="p-4 space-y-2 custom-scrollbar md:overflow-y-auto md:flex-grow">
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

              <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md border border-white/20 md:flex md:flex-col md:overflow-hidden">
                  <div class="flex items-center justify-between p-4 border-b border-white/20">
                    <h2 class="text-lg font-semibold">{{ currentTemplate?.name || '生成的草稿' }}</h2>
                    <div class="relative flex-shrink-0" ref="actionMenuNode">
                      <button @click="isActionMenuOpen = !isActionMenuOpen" class="p-2 rounded-full text-gray-500 hover:bg-gray-500/10 hover:text-gray-800 focus:outline-none transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                          <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
                        </svg>
                      </button>
                      <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
                          <div v-if="isActionMenuOpen" class="absolute right-0 top-full mt-2 w-48 bg-white/90 backdrop-blur-sm rounded-md shadow-lg py-1 z-20 border border-white/30">
                              <a href="#" @click.prevent="saveDraft" class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-500/10 transition-colors">
                                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-3 text-gray-500">
                                      <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
                                  </svg>
                                  <span>保存草稿</span>
                              </a>
                          </div>
                      </transition>
                    </div>
                  </div>
                   <div class="p-4 space-y-4 custom-scrollbar md:flex-grow md:overflow-y-auto">
                      <div v-if="currentTemplate" v-for="field in (currentTemplate.fields || [])" :key="field.id" class="space-y-2">
                          <label :for="field.id" class="font-semibold text-gray-700">{{ field.label }}</label>
                           
                           <!-- Rich Text -->
                           <TiptapEditor v-if="field.type === 'tiptap'" v-model="formValues[field.id]" :placeholder="field.placeholder" @openSettings="openSettings" />
                           
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

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background-color: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  border-radius: 10px;
  background-color: transparent; /* Hide by default */
  transition: background-color .2s;
}

.custom-scrollbar:hover::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.15); /* Show on container hover */
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.3); /* Darker on thumb hover */
}

/* For Firefox */
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent; /* Hide by default */
}
.custom-scrollbar:hover {
  scrollbar-color: rgba(0, 0, 0, 0.15) transparent; /* Show on hover */
}
</style> 