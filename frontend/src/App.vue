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
]);
const selectedSourceTemplateId = ref('weekly');
const selectedTemplateId = ref('weekly');
const formValues = ref({});
const dateValue = ref({
  startDate: '',
  endDate: '',
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

const generateDraft = async () => {
  if (isLoading.value) return;
  isLoading.value = true;
  
  // 1. Mock fetching source reports
  await new Promise(resolve => setTimeout(resolve, 500));
  sourceReports.value = [
      { 
        id: 1, title: '周一的报告', isCollapsed: true,
        fields: [
            { name: '完成的任务', value: '<li>搭建项目 A。</li>' },
            { name: '花费小时数', value: 8 },
        ]
      },
      { 
        id: 2, title: '周二的报告', isCollapsed: true,
        fields: [
            { name: '完成的任务', value: '<li>参加计划会议。</li>' },
            { name: '花费小时数', value: 4 },
        ]
      },
      { 
        id: 3, title: '周三的报告', isCollapsed: true,
        fields: [
            { name: '完成的任务', value: '<li>修复 #1024 错误。</li>' },
            { name: '花费小时数', value: 6 },
        ]
      },
  ];

  // 2. Mock generating the form structure for the weekly report
  await new Promise(resolve => setTimeout(resolve, 500));
  
  // Aggregate data from source reports
  const allTasks = sourceReports.value.map(r => r.fields.find(f => f.name === '完成的任务')?.value || '').join('');
  const totalHours = sourceReports.value.reduce((sum, r) => sum + (r.fields.find(f => f.name === '花费小时数')?.value || 0), 0);

  generatedForm.value = [
    { 
      id: 'summary', name: '本周总结', type: 'text', 
      value: `<h2>工作总结</h2><p>本周工作内容主要包括项目搭建、计划会议和错误修复。</p><ul>${allTasks}</ul>`
    },
    { 
      id: 'total_hours', name: '本周总工时', type: 'number',
      value: totalHours
    },
    {
      id: 'status', name: '项目状态', type: 'dropdown',
      value: 'on_track',
      options: [
        { value: 'on_track', text: '正常' },
        { value: 'at_risk', text: '有风险' },
        { value: 'off_track', text: '已延期' },
      ]
    },
    {
      id: 'screenshot', name: '关键截图', type: 'image', value: null
    },
     {
      id: 'deployment_log', name: '部署日志', type: 'attachment', value: null
    }
  ];

  isLoading.value = false;
};

// Initialize form values when component mounts or template changes
const initializeFormValues = () => {
  const values = {};
  if (currentTemplate.value) {
    currentTemplate.value.fields.forEach(field => {
      // Set default value based on type
      values[field.id] = '';
    });
  }
  formValues.value = values;
};

watch(selectedTemplateId, initializeFormValues, { immediate: true });

</script>

<template>
  <div class="bg-gray-100 min-h-screen font-sans text-gray-800">
    <div class="flex flex-col h-screen">
      
      <header class="flex-shrink-0 bg-white border-b shadow-sm z-10 p-4">
        <h1 class="text-xl font-bold text-gray-800">飞书报告助手</h1>
      </header>
      
      <main class="flex-grow flex flex-col overflow-hidden">
        <!-- Top Controls Section -->
        <section class="flex-shrink-0 bg-white border-b p-4">
            <div class="grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
                <div>
                    <label for="source_template" class="text-sm font-medium text-gray-700">源模板:</label>
                    <select id="source_template" v-model="selectedSourceTemplateId" class="w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500">
                        <option v-for="template in templates" :key="template.id" :value="template.id">
                            {{ template.name }}
                        </option>
                    </select>
                </div>
                <div>
                    <label for="template" class="text-sm font-medium text-gray-700">目标模板</label>
                    <select id="template" v-model="selectedTemplateId" class="mt-1 w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500">
                        <option v-for="template in templates" :key="template.id" :value="template.id">
                            {{ template.name }}
                        </option>
                    </select>
                </div>
                <div>
                    <label for="date" class="text-sm font-medium text-gray-700">报告周期</label>
                    <VueTailwindDatepicker
                        v-model="dateValue"
                        :formatter="{ date: 'YYYY-MM-DD', month: 'MMM' }"
                        as-single
                        i18n="zh-cn"
                        placeholder="选择日期范围"
                        class="w-full"
                    />
                </div>
                <button @click="generateDraft" class="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-3 rounded-lg flex items-center justify-center h-10" :disabled="isLoading">
                    <span>{{ isLoading ? '生成中...' : '生成' }}</span>
                </button>
            </div>
        </section>
        
        <div class="flex-grow grid grid-cols-1 md:grid-cols-2 gap-4 p-4 overflow-hidden">
            <!-- Left Column: Source Reports -->
            <section class="bg-white rounded-lg shadow-md flex flex-col overflow-hidden">
                <h2 class="text-lg font-semibold p-4 border-b">源报告</h2>
                <div class="overflow-y-auto flex-grow p-4 space-y-2">
                    <div v-if="sourceReports.length === 0" class="text-gray-500 text-center pt-10">源报告将在此处显示。</div>
                    <div v-for="report in sourceReports" :key="report.id" class="border rounded-md">
                        <div @click="toggleReportDetail(report)" class="p-3 flex justify-between items-center cursor-pointer hover:bg-gray-50">
                            <h3 class="font-semibold">{{ report.title }}</h3>
                            <svg class="w-5 h-5 transition-transform" :class="{'rotate-180': report.isCollapsed}" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"><path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"/></svg>
                        </div>
                        <div v-if="!report.isCollapsed" class="p-3 border-t prose max-w-none" v-html="report.fields.map(f => `<b>${f.name}:</b> ${f.value}`).join('<br>')"></div>
                    </div>
                </div>
            </section>

            <!-- Right Column: Generated Form -->
            <section class="bg-white rounded-lg shadow-md flex flex-col overflow-hidden">
                 <h2 class="text-lg font-semibold p-4 border-b">生成的草稿表单</h2>
                 <div class="flex-grow overflow-y-auto p-4 space-y-4">
                    <div v-if="currentTemplate" v-for="field in currentTemplate.fields" :key="field.id" class="space-y-2">
                        <label :for="field.id" class="font-semibold text-gray-700">{{ field.label }}</label>
                        <TiptapEditor v-if="field.type === 'tiptap'" v-model="formValues[field.id]" :placeholder="field.placeholder" />
                        <input
                            v-if="field.type === 'text'"
                            :id="field.id"
                            type="text"
                            v-model="formValues[field.id]"
                            :placeholder="field.placeholder"
                            class="w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                        />
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
}
.tiptap:focus {
  outline: none;
}
</style> 