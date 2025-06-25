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
    const response = await this.callAPI(prompt, text, { ...options, stream: true });
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