import { authService } from './authService.js';

// AI 服务工具类
export class AIService {
  constructor() {
    this.apiKey = localStorage.getItem('deepseek_api_key') || '';
    this.baseURL = 'https://api.deepseek.com';
  }

  // 设置 API Key
  setApiKey(key) {
    this.apiKey = key;
    localStorage.setItem('deepseek_api_key', key);
  }

  // 获取 API Key
  getApiKey() {
    return this.apiKey;
  }

  // 检查 API Key 是否已设置
  hasApiKey() {
    return !!this.apiKey;
  }

  // 调用 DeepSeek API
  async callAPI(prompt, text, options = {}) {
    if (!this.hasApiKey()) {
      throw new Error('请先配置 DeepSeek API Key');
    }

    const payload = {
      model: options.model || 'deepseek-chat',
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
      max_tokens: options.maxTokens || 2000
    };

    const response = await fetch(`${this.baseURL}/chat/completions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiKey}`
      },
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(`API 调用失败: ${response.status} ${response.statusText} - ${errorData.error?.message || ''}`);
    }

    return response;
  }

  // 流式处理文本
  async streamProcess(prompt, text, onChunk, options = {}) {
    // 根据options.stream决定是否使用流式，默认为true
    const useStream = options.stream !== false;
    
    const response = await this.callAPI(prompt, text, { ...options, stream: useStream });

    if (!useStream) {
      // 非流式处理：直接返回完整结果
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
              const content = data.choices?.[0]?.delta?.content;
              
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

// 飞书API服务类
export class FeishuApiService {
  constructor() {
    // 使用相对路径，让Vite代理处理
    this.baseURL = '/api';
    this.rawRulesCache = []; // 缓存原始规则数据
  }

  // 固定的模板列表
  getFixedTemplateList() {
    return [
      {
        id: 'monthly_report',
        name: '工作月报'
      },
      {
        id: 'daily_report',
        name: '技术部-工作日报'
      },
      {
        id: 'daily_scrum',
        name: '每日站会'
      },
      {
        id: 'complex_template',
        name: '复杂模板'
      }
    ];
  }

  // 获取所有模板（固定列表 + API内容）
  async getAllTemplates() {
    try {
      const fixedTemplates = this.getFixedTemplateList();
      const templates = [];

      // 清空缓存
      this.rawRulesCache = [];

      for (const template of fixedTemplates) {
        try {
          // 通过API获取模板内容
          const templateContent = await this.getTemplateContent(template.name);
          templates.push(templateContent);
        } catch (error) {
          console.warn(`获取模板"${template.name}"内容失败，使用默认结构:`, error);
          // 如果API获取失败，使用默认模板结构
          const defaultTemplate = this.getDefaultTemplate(template.id, template.name);
          templates.push(defaultTemplate);
        }
      }

      return templates;
    } catch (error) {
      console.error('获取模板列表失败:', error);
      throw error;
    }
  }

  // 通过API获取特定模板内容
  async getTemplateContent(templateName) {
    const url = `${this.baseURL}/rules?name=${encodeURIComponent(templateName)}`;
    const response = await authService.authenticatedFetch(url);
    
    if (!response.ok) {
      throw new Error(`获取模板"${templateName}"失败: ${response.status} ${response.statusText}`);
    }
    
    const rules = await response.json();
    
    if (!rules || rules.length === 0) {
      throw new Error(`模板"${templateName}"未找到`);
    }

    const rule = rules[0]; // 取第一个匹配的模板
    
    // 添加到缓存
    this.rawRulesCache.push(rule);

    // 转换为组件格式
    return {
      id: rule.rule_id || templateName.replace(/[^a-zA-Z0-9]/g, '_'),
      name: rule.name,
      fields: rule.form_schema.map((field, index) => ({
        id: `field_${rule.rule_id || templateName}_${index}`,
        label: field.name,
        type: this.mapFieldType(field.type),
        placeholder: `请输入${field.name}...`,
        // 为特定类型添加额外配置
        ...(field.type === 'dropdown' ? { options: [
            { value: 'option1', text: '选项1' },
            { value: 'option2', text: '选项2' }
          ] } : {}),
        ...(field.type === 'multiSelect' ? { options: [
            { value: 'option1', text: '选项1' },
            { value: 'option2', text: '选项2' }
          ] } : {}),
        ...(field.type === 'image' ? { maxCount: 99, maxSize: 20 * 1024 * 1024 } : {}),
        ...(field.type === 'attachment' ? { maxCount: 9, maxSize: 50 * 1024 * 1024 } : {})
      }))
    };
  }

  // 获取默认模板结构（当API失败时使用）
  getDefaultTemplate(templateId, templateName) {
    const defaultTemplates = {
      'monthly_report': {
        id: 'monthly_report',
        name: '工作月报',
        fields: [
          { id: 'monthly_summary', label: '本月工作总结', type: 'tiptap', placeholder: '请输入本月工作总结...' },
          { id: 'key_achievements', label: '主要成就', type: 'tiptap', placeholder: '请输入主要成就...' },
          { id: 'challenges', label: '遇到的挑战', type: 'tiptap', placeholder: '请输入遇到的挑战...' },
          { id: 'next_month_plan', label: '下月工作计划', type: 'tiptap', placeholder: '请输入下月工作计划...' },
          { id: 'kpi_data', label: '关键指标数据', type: 'text', placeholder: '请输入关键指标数据...' },
          { id: 'team_feedback', label: '团队反馈', type: 'tiptap', placeholder: '请输入团队反馈...' }
        ]
      },
      'daily_report': {
        id: 'daily_report',
        name: '技术部-工作日报',
        fields: [
          { id: 'today_summary', label: '今日总结', type: 'tiptap', placeholder: '请输入今日总结...' },
          { id: 'tomorrow_plan', label: '明日计划', type: 'tiptap', placeholder: '请输入明日计划...' },
          { id: 'need_help', label: '需要协调与帮助', type: 'tiptap', placeholder: '请输入需要协调与帮助的内容...' }
        ]
      },
      'daily_scrum': {
        id: 'daily_scrum',
        name: '每日站会',
        fields: [
          { id: 'yesterday_work', label: '昨日工作', type: 'tiptap', placeholder: '请输入昨日工作内容...' },
          { id: 'today_plan', label: '今日计划', type: 'tiptap', placeholder: '请输入今日计划...' },
          { id: 'blockers', label: '遇到的问题', type: 'tiptap', placeholder: '请输入遇到的问题...' }
        ]
      }
    };

    const defaultTemplate = defaultTemplates[templateId];
    if (defaultTemplate) {
      // 添加默认模板到缓存
      this.rawRulesCache.push({
        rule_id: templateId,
        name: templateName,
        form_schema: defaultTemplate.fields.map(field => ({
          name: field.label,
          type: field.type === 'tiptap' ? 'text' : field.type
        }))
      });
      return defaultTemplate;
    }

    // 如果没有默认模板，返回基础结构
    return {
      id: templateId,
      name: templateName,
      fields: [
        { id: 'content', label: '内容', type: 'tiptap', placeholder: '请输入内容...' }
      ]
    };
  }

  // 根据规则ID获取原始规则数据
  getRawRuleById(ruleId) {
    return this.rawRulesCache.find(rule => rule.rule_id === ruleId);
  }

  // 获取指定名称的模板
  async getRuleByName(templateName) {
    try {
      return await this.getTemplateContent(templateName);
    } catch (error) {
      console.error(`获取模板"${templateName}"失败:`, error);
      // 如果API失败，尝试从固定列表中查找并返回默认模板
      const fixedTemplates = this.getFixedTemplateList();
      const fixedTemplate = fixedTemplates.find(t => t.name === templateName);
      if (fixedTemplate) {
        return this.getDefaultTemplate(fixedTemplate.id, fixedTemplate.name);
      }
      throw error;
    }
  }

     // 获取报告列表  
   async getReports(params = {}, templateData = null) {
     try {
       const queryParams = new URLSearchParams();
       if (params.rule_id) queryParams.append('rule_id', params.rule_id);
       if (params.start_time) queryParams.append('start_time', params.start_time);
       if (params.end_time) queryParams.append('end_time', params.end_time);
       
       const url = `${this.baseURL}/reports${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
       const response = await authService.authenticatedFetch(url);
       
       if (!response.ok) {
         throw new Error(`获取报告失败: ${response.status} ${response.statusText}`);
       }
       
       const data = await response.json();
       
       // 检查数据格式并提供默认值
       if (!data || !data.items || !Array.isArray(data.items)) {
         console.warn('API返回的数据格式不正确:', data);
         return [];
       }
       
       // 转换为组件所需的格式
       return data.items.map(report => {
         // 如果有模板数据，尝试匹配字段类型
         const fieldTypeMap = {};
         if (templateData && templateData.rule_id === report.rule_id) {
           templateData.form_schema.forEach(schemaField => {
             fieldTypeMap[schemaField.name] = this.mapFieldType(schemaField.type);
           });
         }
         
         return {
           id: report.task_id,
           title: `${report.rule_name} - ${report.from_user_name} (${this.formatTime(report.commit_time)})`,
           isCollapsed: true,
           fields: (report.form_contents || []).map(content => ({
             name: content.field_name,
             value: this.formatFieldValue(content.field_value, fieldTypeMap[content.field_name] || 'tiptap'),
             type: fieldTypeMap[content.field_name] || 'tiptap'
           }))
         };
       });
     } catch (error) {
       console.error('获取报告失败:', error);
       throw error;
     }
   }

   // 格式化字段值
   formatFieldValue(value, type) {
     if (!value) return value;
     
     switch (type) {
       case 'tiptap':
         // 如果是纯文本，包装成HTML段落
         if (typeof value === 'string' && !value.includes('<')) {
           return `<p>${value}</p>`;
         }
         return value;
       case 'multiSelect':
         // 如果是数组字符串，尝试解析
         if (typeof value === 'string' && value.startsWith('[')) {
           try {
             return JSON.parse(value);
           } catch {
             return [value];
           }
         }
         return Array.isArray(value) ? value : [value];
       case 'image':
       case 'attachment':
         // 如果是文件类型，确保返回数组格式
         if (typeof value === 'string') {
           try {
             const parsed = JSON.parse(value);
             return Array.isArray(parsed) ? parsed : [parsed];
           } catch {
             return [];
           }
         }
         return Array.isArray(value) ? value : [];
       default:
         return value;
     }
   }

     // 字段类型映射
   mapFieldType(apiType) {
     const typeMap = {
       'text': 'tiptap', // 默认使用富文本编辑器
       'number': 'number', 
       'dropdown': 'dropdown',
       'image': 'image',
       'attachment': 'attachment',
       'multiSelect': 'multiSelect',
       'address': 'address',
       'datetime': 'datetime'
     };
     return typeMap[apiType] || 'tiptap';
   }

  // 时间格式化
  formatTime(timestamp) {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString('zh-CN');
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
export const feishuApiService = new FeishuApiService();
export const reportSummarizer = new ReportSummarizer(aiService); 