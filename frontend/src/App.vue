<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import TiptapEditor from './components/TiptapEditor.vue';
import VueTailwindDatepicker from 'vue-tailwind-datepicker';

// --- State Management ---

const isAuthenticated = ref(true);
const isLoading = ref(false);

const templates = ref([
  {
    id: 'weekly',
    name: '周报模板',
    fields: [
      { id: 'summary', label: '本周工作总结', type: 'tiptap', placeholder: '请输入本周工作总结...' },
      { id: 'plan', label: '下周工作计划', type: 'tiptap', placeholder: '请输入下周工作计划...' },
      { id: 'risk', label: '风险与建议', type: 'tiptap', placeholder: '请输入风险与建议...' },
    ],
  },
  {
    id: 'monthly',
    name: '月报模板',
    fields: [
      { id: 'kpi', label: '本月关键绩效', type: 'text', placeholder: '请输入本月关键绩效...' },
      { id: 'achievements', label: '主要成就和产出', type: 'tiptap', placeholder: '请输入主要成就和产出...' },
      { id: 'learnings', label: '心得与反思', type: 'tiptap', placeholder: '请输入心得与反思...' },
      { id: 'next_month_goals', label: '下月核心目标', type: 'tiptap', placeholder: '请输入下月核心目标...' },
    ],
  },
  {
    id: 'comprehensive',
    name: '综合报告模板',
    fields: [
      { id: 'project_name', label: '项目名称', type: 'text', placeholder: '请输入项目名称...' },
      { id: 'progress_rating', label: '项目进度 (1-5)', type: 'number', placeholder: '请输入1-5的数字' },
      { id: 'status', label: '当前状态', type: 'dropdown', options: [
          { value: 'on_track', text: '正常推进' },
          { value: 'at_risk', text: '存在风险' },
          { value: 'delayed', text: '已延期' },
        ],
      },
      { id: 'stakeholders', label: '关键干系人', type: 'multiSelect', options: [
          { value: 'product', text: '产品' },
          { value: 'design', text: '设计' },
          { value: 'dev', text: '开发' },
          { value: 'qa', text: '测试' },
        ],
      },
      { id: 'summary', label: '图文总结', type: 'tiptap', placeholder: '请输入图文总结...' },
      { id: 'location', label: '办公地址', type: 'address', placeholder: '请输入详细地址...' },
      { id: 'meeting_time', label: '下次会议时间', type: 'datetime' },
      { id: 'screenshot', label: '效果图', type: 'image', maxCount: 99, maxSize: 20 * 1024 * 1024 }, // 20MB
      { id: 'log_file', label: '日志文件', type: 'attachment', maxCount: 9, maxSize: 50 * 1024 * 1024 }, // 50MB
    ],
  },
]);
const selectedSourceTemplateId = ref('weekly');
const selectedTemplateId = ref('comprehensive');
const formValues = ref({});
const fileInputRefs = ref({});
const dateValue = ref({
  startDate: '',
  endDate: '',
});

const isProfileMenuOpen = ref(false);
const profileMenuNode = ref(null);
const notifications = ref([]);

onMounted(() => {
  document.addEventListener('click', (event) => {
    if (profileMenuNode.value && !profileMenuNode.value.contains(event.target)) {
      isProfileMenuOpen.value = false;
    }
  });
});

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
  const today = new Date();
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  const last7Days = new Date();
  last7Days.setDate(last7Days.getDate() - 6);
  const last30Days = new Date();
  last30Days.setDate(last30Days.getDate() - 29);
  const thisMonthStart = new Date(today.getFullYear(), today.getMonth(), 1);
  const lastMonthEnd = new Date(today.getFullYear(), today.getMonth(), 0);
  const lastMonthStart = new Date(today.getFullYear(), today.getMonth() - 1, 1);

  return [
    { label: '今天', at: [today, today] },
    { label: '昨天', at: [yesterday, yesterday] },
    { label: '过去7天', at: [last7Days, today] },
    { label: '过去30天', at: [last30Days, today] },
    { label: '本月', at: [thisMonthStart, today] },
    { label: '上个月', at: [lastMonthStart, lastMonthEnd] },
  ]
}

const currentTemplate = computed(() => {
  return templates.value.find(t => t.id === selectedTemplateId.value);
});

const sourceReports = ref([]);
const generatedForm = ref([]); // This will hold the structured form fields for the right panel

// --- Methods ---

const toggleReportDetail = (report) => {
    report.isCollapsed = !report.isCollapsed;
};

const getReports = () => {
  sourceReports.value = [
    { 
      id: 'rep-comp-1', 
      title: '综合报告示例 - Q2', 
      isCollapsed: true,
      fields: [
        { name: '项目名称', value: '飞书报告助手', type: 'text' },
        { name: '项目进度 (1-5)', value: 4, type: 'number' },
        { name: '当前状态', value: '正常推进', type: 'dropdown' },
        { name: '关键干系人', value: ['产品', '开发'], type: 'multiSelect' },
        { name: '办公地址', value: '未来城A座-501', type: 'address' },
        { name: '下次会议时间', value: '2025-08-01T10:00', type: 'datetime' },
        { name: '效果图', value: [{ name: 'placeholder.png', size: 818200, url: 'https://template.tiptap.dev/images/placeholder-image.png' }], type: 'image' },
        { name: '日志文件', value: [{ name: 'mindmap.png', size: 818200, url: 'https://template.tiptap.dev/images/placeholder-image.png' }], type: 'attachment' },
        { name: '图文总结', value: '<p>这是富文本总结，包含图片。</p><img src="https://template.tiptap.dev/images/placeholder-image.png"/>', type: 'tiptap' },
      ] 
    },
    { 
      id: 'rep-2', 
      title: '客户端团队周报 2024-07-01 ~ 2024-07-07', 
      isCollapsed: true,
      fields: [
        { name: 'iOS端进展', value: '新版本已提交审核。' },
        { name: 'Android端进展', value: '正在进行性能优化。' }
      ]
    },
    { 
      id: 'rep-3', 
      title: '服务端架构升级讨论纪要', 
      isCollapsed: true,
      fields: [
        { name: '讨论决议', value: '<li>采用微服务架构。</li><li>数据库方案选型为PostgreSQL。</li>' },
        { name: '后续任务', value: '由王五负责输出详细设计文档。' }
      ]
    },
  ];
};

const generateDraft = () => {
  if (!currentTemplate.value) return;
  isLoading.value = true;

  setTimeout(() => {
    const newValues = {};
    currentTemplate.value.fields.forEach(field => {
      const randomContent = `这是为"${field.label}"随机生成的内容。现在是 ${new Date().toLocaleTimeString()}。`;
      switch (field.type) {
        case 'tiptap':
          newValues[field.id] = `<p>${randomContent}</p>`;
          break;
        case 'text':
        case 'address':
          newValues[field.id] = randomContent;
          break;
        case 'number':
          newValues[field.id] = Math.floor(Math.random() * 5) + 1;
          break;
        case 'dropdown':
          newValues[field.id] = field.options[Math.floor(Math.random() * field.options.length)].value;
          break;
        case 'multiSelect':
          // Select 1 to N options randomly
          newValues[field.id] = field.options
            .filter(() => Math.random() > 0.5)
            .map(opt => opt.value);
          if (newValues[field.id].length === 0 && field.options.length > 0) {
            newValues[field.id].push(field.options[0].value); // ensure at least one is selected
          }
          break;
        case 'datetime':
          newValues[field.id] = new Date(Date.now() + Math.random() * 1000 * 3600 * 24 * 7).toISOString().substring(0, 16);
          break;
        case 'image':
        case 'attachment':
          newValues[field.id] = [{
            name: 'placeholder.png',
            size: 818200,
            url: 'https://template.tiptap.dev/images/placeholder-image.png',
            id: Date.now()
          }];
          break;
        default:
          newValues[field.id] = '';
      }
    });
    formValues.value = newValues;
    isLoading.value = false;
    addNotification('草稿已生成', '已为您填充好表单内容。', 'success');
  }, 1000);
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
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" /></svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </transition-group>
    </div>
  </div>

  <div class="min-h-screen w-full flex items-center justify-center p-4 bg-gray-900/10">
    <div class="w-full max-w-screen-2xl h-[90vh] bg-white/60 backdrop-blur-xl rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-white/30 relative">
      
      <!-- Loading Overlay -->
      <div v-if="isLoading" class="absolute inset-0 bg-gray-900/30 backdrop-blur-sm flex items-center justify-center z-50 rounded-2xl">
        <div class="flex flex-col items-center">
          <svg class="animate-spin -ml-1 mr-3 h-10 w-10 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="mt-4 text-white text-lg font-semibold">正在生成草稿...</span>
        </div>
      </div>
      
      <header class="flex-shrink-0 flex items-center justify-between border-b border-white/30 shadow-sm z-10 p-4">
        <h1 class="text-xl font-bold text-gray-800">飞书报告助手</h1>
        
        <!-- User Profile Section -->
        <div class="flex items-center space-x-3">
          <div class="relative" ref="profileMenuNode">
            <button @click="isProfileMenuOpen = !isProfileMenuOpen" 
                    class="flex items-center space-x-2 p-1 rounded-md transition-colors duration-200"
                    :class="[isProfileMenuOpen ? 'bg-gray-500/20' : 'hover:bg-gray-500/10']">
              <span class="inline-flex items-center justify-center h-8 w-8 rounded-full bg-indigo-500 text-white font-bold">
                A
              </span>
              <span class="text-gray-700 font-semibold">Admin</span>
              <svg class="w-5 h-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
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
              <div v-if="isProfileMenuOpen" class="absolute right-0 mt-2 w-48 bg-white/80 backdrop-blur-sm rounded-md shadow-lg py-1 z-20 border border-white/30">
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-white/50">个人资料</a>
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-white/50">设置</a>
                <div class="border-t border-gray-200/50"></div>
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-white/50">退出登录</a>
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
                  <button @click="getReports" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-white/20 text-gray-700 backdrop-blur-sm border border-white/30 shadow-lg hover:bg-white/30 hover:text-gray-900 focus:outline-none">
                    获取报告
                  </button>
                  <button @click="generateDraft" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-indigo-500/20 text-indigo-800 backdrop-blur-sm border border-indigo-500/30 shadow-lg hover:bg-indigo-500/40 hover:text-indigo-900 focus:outline-none">
                    生成草稿
                  </button>
                  <button @click="addNotification('发送成功', '报告已提交，请在飞书中查收。')" class="w-full px-4 py-2 rounded-lg font-semibold transition-all duration-300 bg-blue-500/20 text-blue-800 backdrop-blur-sm border border-blue-500/30 shadow-lg hover:bg-blue-500/40 hover:text-blue-900 focus:outline-none">
                    发送到飞书
                  </button>
                </div>
            </div>
        </section>
        
        <div class="flex-grow grid grid-cols-1 md:grid-cols-2 gap-4 p-4 overflow-hidden">
            <!-- Left Column: Source Reports -->
            <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md flex flex-col overflow-hidden border border-white/20">
                <h2 class="text-lg font-semibold p-4 border-b border-white/20">源报告</h2>
                <div class="overflow-y-auto flex-grow p-4 space-y-2">
                    <div v-if="sourceReports.length === 0" class="text-gray-500 text-center pt-10">源报告将在此处显示。</div>
                    <div v-for="report in sourceReports" :key="report.id" class="border rounded-md">
                        <div @click="toggleReportDetail(report)" class="p-3 flex justify-between items-center cursor-pointer hover:bg-gray-50">
                            <h3 class="font-semibold">{{ report.title }}</h3>
                            <svg class="w-5 h-5 transition-transform" :class="{'rotate-180': !report.isCollapsed}" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"><path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"/></svg>
                        </div>
                        <div v-if="!report.isCollapsed" class="p-4 border-t border-white/20 prose max-w-none text-sm">
                          <div v-for="field in report.fields" :key="field.name" class="mb-3">
                              <b class="text-gray-700">{{ field.name }}:</b>
                              
                              <div v-if="field.type === 'image'" class="mt-2 grid grid-cols-4 gap-2">
                                  <div v-for="file in field.value" :key="file.name" class="relative">
                                    <img :src="file.url" alt="图片预览" class="w-full h-auto rounded-md shadow-md border border-gray-200" />
                                  </div>
                              </div>
                              <div v-else-if="field.type === 'attachment'" class="mt-1 space-y-2">
                                  <div v-for="file in field.value" :key="file.name">
                                      <a :href="file.url" target="_blank" class="text-indigo-600 hover:underline">{{ file.name }} ({{ formatFileSize(file.size) }})</a>
                                  </div>
                              </div>
                              <div v-else-if="field.type === 'multiSelect'" class="inline-flex flex-wrap gap-2 mt-1">
                                <span v-for="item in field.value" :key="item" class="bg-gray-200 text-gray-700 px-2 py-0.5 rounded-full text-xs">{{ item }}</span>
                              </div>
                              <div v-else-if="field.type === 'tiptap'" class="mt-1 border rounded-md p-2 bg-gray-50/50" v-html="field.value"></div>
                              <span v-else class="ml-2 text-gray-800">{{ field.value }}</span>
                          </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Right Column: Generated Form -->
            <section class="bg-white/50 backdrop-blur-sm rounded-lg shadow-md flex flex-col overflow-hidden border border-white/20">
                 <h2 class="text-lg font-semibold p-4 border-b border-white/20">生成的草稿表单</h2>
                 <div class="flex-grow overflow-y-auto p-4 space-y-4">
                    <div v-if="currentTemplate" v-for="field in currentTemplate.fields" :key="field.id" class="space-y-2">
                        <label :for="field.id" class="font-semibold text-gray-700">{{ field.label }}</label>
                        
                        <!-- Rich Text -->
                        <TiptapEditor v-if="field.type === 'tiptap'" v-model="formValues[field.id]" :placeholder="field.placeholder" />
                        
                        <!-- Number -->
                        <input v-else-if="field.type === 'number'" :id="field.id" type="number" v-model.number="formValues[field.id]" :placeholder="field.placeholder" class="form-input" />

                        <!-- Dropdown -->
                        <select v-else-if="field.type === 'dropdown'" :id="field.id" v-model="formValues[field.id]" class="form-input">
                            <option disabled value="">请选择</option>
                            <option v-for="opt in field.options" :key="opt.value" :value="opt.value">{{ opt.text }}</option>
                        </select>

                        <!-- Multi-Select Checkboxes -->
                        <div v-else-if="field.type === 'multiSelect'" class="flex flex-wrap gap-x-4 gap-y-2 pt-1">
                           <label v-for="opt in field.options" :key="opt.value" class="flex items-center space-x-2 cursor-pointer">
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
                              <div v-for="file in formValues[field.id]" :key="file.id" class="relative w-24 h-24">
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
                              <div v-for="file in formValues[field.id]" :key="file.id" class="flex items-center justify-between p-2 bg-gray-100/80 rounded-md">
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
    </div>
  </div>
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

</style> 