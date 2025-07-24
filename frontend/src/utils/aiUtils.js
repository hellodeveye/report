import { authService } from './authService.js';
import { apiService } from './apiService.js';

// --- 模型配置 ---
export const MODELS_CONFIG = {
  deepseek: {
    label: 'DeepSeek',
    baseURL: 'https://api.deepseek.com',
    models: [
      { id: 'deepseek-chat', name: 'DeepSeek Chat (推荐)' },
      { id: 'deepseek-reasoner', name: 'DeepSeek Reasoner' },
    ],
    // 将标准输入转换为提供商特定的API格式
    transformPayload: (payload) => ({
      model: payload.model,
      messages: payload.messages,
      stream: payload.stream,
      temperature: payload.temperature,
      max_tokens: payload.max_tokens,
    }),
    // 从提供商的响应中提取内容
    extractContent: (data) => data.choices?.[0]?.delta?.content,
  },
  doubao: {
    label: '火山方舟 (豆包)',
    baseURL: 'https://ark.cn-beijing.volces.com/api/v3',
    models: [
      { id: 'kimi-k2-250711', name: 'Kimi-K2' },
      { id: 'doubao-1-5-pro-32k-250115', name: 'Doubao Pro 32k' },
      { id: 'doubao-1-5-lite-32k-250115', name: 'Doubao Lite 32k' },
    ],
    transformPayload: (payload) => ({
      model: payload.model,
      messages: payload.messages,
      stream: payload.stream,
      temperature: payload.temperature,
    }),
    extractContent: (data) => data.choices?.[0]?.delta?.content,
  },
};

// AI 服务工具类
export class AIService {
  constructor() {
    this.settings = this.getSettings();
  }

  // 获取所有AI设置
  getSettings() {
    try {
      const settingsStr = localStorage.getItem('ai_settings');
      if (settingsStr) {
        return JSON.parse(settingsStr);
      }
    } catch (e) {
      console.error('Failed to parse AI settings from localStorage', e);
    }
    // Return default settings if nothing is stored or parsing fails
    return {
      provider: 'deepseek',
      apiKey: '',
      model: MODELS_CONFIG.deepseek.models[0].id,
    };
  }

  // 保存所有AI设置
  saveSettings(settings) {
    this.settings = settings;
    localStorage.setItem('ai_settings', JSON.stringify(settings));
  }

  // 检查 API Key 是否已设置
  hasApiKey() {
    return !!this.settings.apiKey;
  }

  // 调用 AI API
  async callAPI(prompt, text, options = {}) {
    if (!this.hasApiKey()) {
      throw new Error('请先在设置中配置AI模型的API Key');
    }

    const providerConfig = MODELS_CONFIG[this.settings.provider];
    if (!providerConfig) {
      throw new Error(`未知的AI提供商: ${this.settings.provider}`);
    }

    const standardPayload = {
      model: this.settings.model,
      messages: [
        {
          role: 'system',
          content: options.systemPrompt || '你是一个专业的文本编辑助手，请根据用户的要求对文本进行处理。直接返回处理后的结果，不要添加额外的解释或格式。'
        },
        {
          role: 'user',
          content: `${prompt}\n\n${text}`
        }
      ],
      stream: options.stream !== false,
      temperature: options.temperature || 0.7,
      max_tokens: options.maxTokens || 2000,
    };

    const apiPayload = providerConfig.transformPayload(standardPayload);

    const response = await fetch(`${providerConfig.baseURL}/chat/completions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.settings.apiKey}`
      },
      body: JSON.stringify(apiPayload)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(`API 调用失败: ${response.status} ${response.statusText} - ${errorData.error?.message || ''}`);
    }

    return response;
  }

  // 流式处理文本
  async streamProcess(prompt, text, onChunk, options = {}) {
    const useStream = options.stream !== false;
    const response = await this.callAPI(prompt, text, { ...options, stream: useStream });
    
    const providerConfig = MODELS_CONFIG[this.settings.provider];

    if (!useStream) {
      const result = await response.json();
      const content = result.choices?.[0]?.message?.content || '';
      return content;
    }

    // 流式处理
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let accumulatedText = '';

    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        const lines = chunk.split('\n');

        for (const line of lines) {
          if (line.startsWith('data: ') && line !== 'data: [DONE]') {
            try {
              const data = JSON.parse(line.slice(6));
              const content = providerConfig.extractContent(data);
              
              if (content) {
                accumulatedText += content;
                if (onChunk) {
                  onChunk(content, accumulatedText);
                }
              }
            } catch (e) {
              console.error('解析 SSE 数据错误:', e);
            }
          }
        }
      }
    } finally {
      reader.releaseLock();
    }

    return accumulatedText;
  }
}

// 报告汇总功能
export class ReportSummarizer {
  constructor(aiService) {
    this.aiService = aiService;
  }

  // 智能汇总多个报告
  async summarizeReports(sourceReports, targetTemplate) {
    if (!sourceReports || sourceReports.length === 0) {
      throw new Error('没有源报告数据');
    }

    const summary = {};
    
    // 对目标模板的每个字段生成内容
    for (const field of targetTemplate.fields) {
      try {
        const content = await this.generateFieldContent(sourceReports, field);
        summary[field.id] = content;
      } catch (error) {
        console.warn(`生成字段 ${field.label} 失败:`, error);
        summary[field.id] = `请手动填写 ${field.label}`;
      }
    }

    return summary;
  }

  // 为特定字段生成内容
  async generateFieldContent(sourceReports, targetField) {
    const reportContents = this.extractReportContents(sourceReports);
    const prompt = this.buildPrompt(targetField, reportContents);
    
    return await this.aiService.streamProcess(
      prompt,
      reportContents,
      null,
      { stream: false }
    );
  }

  // 提取报告内容
  extractReportContents(reports) {
    return reports.map(report => {
      const content = {
        title: report.title,
        fields: {}
      };
      
      report.fields.forEach(field => {
        content.fields[field.name] = field.value;
      });
      
      return content;
    }).map(report => {
      return `【${report.title}】\n${Object.entries(report.fields)
        .map(([key, value]) => `${key}: ${value}`)
        .join('\n')}\n`;
    }).join('\n---\n');
  }

  // 构建针对性prompt
  buildPrompt(targetField, reportContents) {
    const fieldMappings = {
      '本月总结': '请基于以下日报内容，撰写本月工作总结，突出主要成就和完成的工作',
      '本月工作总结': '请基于以下日报内容，撰写本月工作总结，突出主要成就和完成的工作',
      '主要成就': '请从以下日报中提取并总结主要成就和亮点',
      '进展同步': '请基于以下日报，总结项目进展情况',
      '下月计划': '请基于以下日报中的计划内容，制定下月工作计划',
      '下月工作计划': '请基于以下日报中的计划内容，制定下月工作计划',
      '复盘总结': '请基于以下日报，进行工作复盘，分析经验教训',
      '遇到的挑战': '请从以下日报中提取遇到的问题和挑战',
      '团队反馈': '请基于以下日报，总结团队协作和反馈情况'
    };

    const specificPrompt = fieldMappings[targetField.label] || 
      `请基于以下报告内容，为"${targetField.label}"生成合适的内容`;

    return `${specificPrompt}。要求：
1. 内容简洁明了，突出重点
2. 保持专业的工作报告语气
3. 如果是富文本字段，可以使用适当的HTML格式
4. 字数控制在100-300字之间

以下是源报告内容：`;
  }
}

// AI 提示词配置
export const AI_PROMPTS = {
  '重构': {
    prompt: '请重新组织和改进以下文本的结构和表达，使其更加清晰、逻辑性更强：',
    description: '重新组织文本结构，提高逻辑性和可读性'
  },
  '博客化': {
    prompt: '请将以下内容改写成适合博客发布的风格，要求有吸引力的标题、清晰的段落结构和引人入胜的表达：',
    description: '转换为博客风格，增加吸引力和可读性'
  },
  '提取要点': {
    prompt: '请从以下文本中提取出主要要点，用简洁的条目列出：',
    description: '提取并列出文本的核心要点'
  },
  '改写': {
    prompt: '请用不同的表达方式重新表述以下内容，保持原意但改变用词和句式：',
    description: '保持原意的情况下重新表述内容'
  },
  '缩短': {
    prompt: '请将以下内容压缩成更简洁的版本，保留核心信息：',
    description: '压缩文本长度，保留核心信息'
  },
  '扩写': {
    prompt: '请展开以下内容，添加更多细节、例子或解释，使其更加丰富完整：',
    description: '增加细节和例子，丰富内容'
  },
  '总结': {
    prompt: '请对以下内容进行总结，突出主要观点和结论：',
    description: '总结主要观点和结论'
  },
  '简化': {
    prompt: '请将以下内容简化，使用更简单易懂的语言表达：',
    description: '使用简单易懂的语言重新表达'
  },
  '修正拼写': {
    prompt: '请检查并修正以下文本中的拼写、语法和标点错误：',
    description: '检查并修正语法、拼写和标点错误'
  },
  '继续写作': {
    prompt: '请基于以下内容继续写作，保持相同的风格和主题：',
    description: '基于现有内容继续写作'
  },
  '使用激动语气': {
    prompt: '请将以下内容改写成充满激情和活力的语调：',
    description: '转换为充满激情的语调'
  },
  '添加表情 🙂': {
    prompt: '请在以下文本中适当位置添加表情符号，使其更生动有趣：',
    description: '添加表情符号，增加趣味性'
  },
  '去除表情': {
    prompt: '请从以下文本中移除所有表情符号，保持正式的文本风格：',
    description: '移除表情符号，保持正式风格'
  },
  '翻译成瑞典语': {
    prompt: 'Please translate the following text to Swedish:',
    description: '翻译为瑞典语'
  },
  '翻译成德语': {
    prompt: 'Please translate the following text to German:',
    description: '翻译为德语'
  },
  '翻译成英文': {
    prompt: 'Please translate the following text to English:',
    description: '翻译为英文'
  },
  '翻译成老挝语': {
    prompt: 'Please translate the following text to Lao (Laotian):',
    description: '翻译为老挝语'
  },
  '翻译成中文': {
    prompt: 'Please translate the following text to Chinese (Simplified):',
    description: '翻译为中文'
  },
  '一句话总结': {
    prompt: '请用一句话总结以下内容的核心观点：',
    description: '用一句话概括核心观点'
  }
};

// 创建单例实例
export const aiService = new AIService();
export const reportSummarizer = new ReportSummarizer(aiService); 