<template>
  <div class="border border-gray-200 rounded-lg relative">
    <div v-if="editor" class="flex items-center flex-wrap p-2 border-b border-gray-200 bg-gray-50 rounded-t-lg">
      <!-- History -->
      <button @click="editor.chain().focus().undo().run()" :disabled="!editor.can().undo()" class="toolbar-button">撤销</button>
      <button @click="editor.chain().focus().redo().run()" :disabled="!editor.can().redo()" class="toolbar-button">重做</button>
      <div class="divider"></div>
      
      <!-- Headings -->
      <button @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }" class="toolbar-button font-bold">H2</button>
      <button @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }" class="toolbar-button font-bold">H3</button>
      <div class="divider"></div>

      <!-- Basic Formatting -->
      <button @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }" class="toolbar-button font-bold">B</button>
      <button @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }" class="toolbar-button italic">I</button>
      <button @click="editor.chain().focus().toggleUnderline().run()" :class="{ 'is-active': editor.isActive('underline') }" class="toolbar-button underline">U</button>
      <button @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }" class="toolbar-button line-through">S</button>
      <button @click="editor.chain().focus().toggleCode().run()" :class="{ 'is-active': editor.isActive('code') }" class="toolbar-button font-mono">代码</button>
      <button @click="editor.chain().focus().toggleHighlight().run()" :class="{ 'is-active': editor.isActive('highlight') }" class="toolbar-button bg-yellow-200">H</button>
      <div class="divider"></div>

      <!-- Scripting and Links -->
      <button @click="editor.chain().focus().toggleSubscript().run()" :class="{ 'is-active': editor.isActive('subscript') }" class="toolbar-button">下标</button>
      <button @click="editor.chain().focus().toggleSuperscript().run()" :class="{ 'is-active': editor.isActive('superscript') }" class="toolbar-button">上标</button>
      <button @click="setLink" :class="{ 'is-active': editor.isActive('link') }" class="toolbar-button">链接</button>
      <div class="divider"></div>
      
      <!-- Alignment -->
      <button @click="editor.chain().focus().setTextAlign('left').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'left' }) }" class="toolbar-button">居左</button>
      <button @click="editor.chain().focus().setTextAlign('center').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'center' }) }" class="toolbar-button">居中</button>
      <button @click="editor.chain().focus().setTextAlign('right').run()" :class="{ 'is-active': editor.isActive({ textAlign: 'right' }) }" class="toolbar-button">居右</button>
      <div class="divider"></div>

      <!-- Lists and Images -->
      <button @click="editor.chain().focus().toggleBulletList().run()" :class="{ 'is-active': editor.isActive('bulletList') }" class="toolbar-button">列表</button>
      <button @click="addImage" class="toolbar-button">图片</button>
      <div class="divider"></div>
      <div class="relative">
        <button @click="toggleAiDropdown" :disabled="editor.state.selection.empty || isAiProcessing" class="toolbar-button flex items-center">
          <span v-if="isAiProcessing" class="animate-spin mr-1">⏳</span>
          AI 优化
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
        </button>
        <div v-if="isAiDropdownOpen" class="absolute z-10 mt-1 bg-white rounded-md border border-gray-200 w-64 max-h-60 overflow-y-auto">
          <ul class="py-1">
            <li v-for="option in aiOptions" :key="option">
              <a href="#" @click.prevent="applyAiAction(option)" 
                 class="block px-4 py-2 text-sm hover:bg-gray-100"
                 :class="{ 'opacity-50 cursor-not-allowed': isAiProcessing }"
                 :title="AI_PROMPTS[option]?.description">
                <div class="font-medium text-gray-900">{{ option }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ AI_PROMPTS[option]?.description }}</div>
              </a>
            </li>
          </ul>
        </div>
      </div>
      

    </div>
    <EditorContent :editor="editor" />
    
    <!-- AI 生成状态提示 -->
    <div v-if="isAiProcessing" class="absolute bottom-4 left-1/2 transform -translate-x-1/2 z-20">
      <div class="flex items-center space-x-2 bg-white border border-gray-200 rounded-full px-3 py-1 shadow-sm">
        <div class="flex space-x-1">
          <div class="w-1.5 h-1.5 bg-indigo-400 rounded-full animate-bounce opacity-90"></div>
          <div class="w-1.5 h-1.5 bg-indigo-500 rounded-full animate-bounce opacity-90" style="animation-delay: 0.1s"></div>
          <div class="w-1.5 h-1.5 bg-indigo-600 rounded-full animate-bounce opacity-90" style="animation-delay: 0.2s"></div>
        </div>
        <span class="text-gray-600 text-sm font-medium">AI 生成中...</span>
      </div>
    </div>
    

  </div>
</template>

<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import Underline from '@tiptap/extension-underline';
import Highlight from '@tiptap/extension-highlight';
import Link from '@tiptap/extension-link';
import Subscript from '@tiptap/extension-subscript';
import Superscript from '@tiptap/extension-superscript';
import TextAlign from '@tiptap/extension-text-align';
import Image from '@tiptap/extension-image';
import Placeholder from '@tiptap/extension-placeholder';
import { watch, ref } from 'vue';
import { aiService, AI_PROMPTS } from '../utils/aiService.js';

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '请输入...' },
});

const emit = defineEmits(['update:modelValue', 'showApiKeyConfig']);

const isAiDropdownOpen = ref(false);
const isAiProcessing = ref(false);

const aiOptions = Object.keys(AI_PROMPTS);

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    Underline,
    Highlight,
    Link.configure({ openOnClick: false }),
    Subscript,
    Superscript,
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Image,
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
  ],
  editorProps: {
    attributes: { class: 'prose max-w-none p-4 focus:outline-none min-h-[150px]' },
  },
  onUpdate: () => {
    emit('update:modelValue', editor.value.getHTML());
  },
});

const toggleAiDropdown = () => {
    if (editor.value.state.selection.empty || isAiProcessing.value) return;
    isAiDropdownOpen.value = !isAiDropdownOpen.value;
};



const applyAiAction = async (action) => {
  if (isAiProcessing.value) return;
  
  isAiDropdownOpen.value = false;
  isAiProcessing.value = true;

  // 保存原始选择范围和选中的文本
  const originalSelection = {
    from: editor.value.state.selection.from,
    to: editor.value.state.selection.to
  };
  
  const selectedText = editor.value.state.doc.textBetween(originalSelection.from, originalSelection.to, ' ');

  try {
    if (!selectedText.trim()) {
      alert('请先选择要优化的文本');
      return;
    }

    const promptConfig = AI_PROMPTS[action];
    if (!promptConfig) {
      alert('未找到对应的处理方式');
      return;
    }

    // 检查 API Key
    if (!aiService.hasApiKey()) {
      emit('showApiKeyConfig', 'ai_text_optimize');
      return;
    }

    let hasStartedReplacing = false;
    let currentPosition = originalSelection.from;
    
    // 流式处理，实时更新编辑器内容
    const result = await aiService.streamProcess(
      promptConfig.prompt,
      selectedText,
      // 流式回调，实时插入内容
      (chunk, accumulatedText) => {
        if (chunk && chunk.trim()) {
          // 如果这是第一次接收到内容，先删除选中的文本
          if (!hasStartedReplacing) {
            editor.value.chain()
              .focus()
              .setTextSelection({ from: originalSelection.from, to: originalSelection.to })
              .deleteSelection()
              .run();
            hasStartedReplacing = true;
            currentPosition = originalSelection.from;
          }
          
          // 在当前位置插入新的chunk
          editor.value.chain()
            .focus()
            .setTextSelection({ from: currentPosition, to: currentPosition })
            .insertContent(chunk)
            .run();
          
          // 更新当前位置
          currentPosition += chunk.length;
        }
      }
    );

    // 流式完成后的最终处理
    if (!result || !result.trim()) {
      if (hasStartedReplacing) {
        // 如果已经开始替换但失败了，恢复原始文本
        editor.value.chain()
          .focus()
          .setTextSelection({ from: originalSelection.from, to: currentPosition })
          .insertContent(selectedText)
          .run();
        alert('AI 处理失败，已恢复原文本');
      } else {
        // 如果还没有开始替换就失败了，原文本仍然存在，只需提示
        alert('AI 处理失败，原文本保持不变');
      }
    }

  } catch (error) {
    console.error('AI 处理错误:', error);
    
    // 发生错误时的处理
    try {
      if (hasStartedReplacing) {
        // 如果已经开始替换但出错了，尝试恢复原始文本
        editor.value.chain()
          .focus()
          .setTextSelection({ from: originalSelection.from, to: currentPosition })
          .insertContent(selectedText)
          .run();
        alert(`AI 处理失败，已恢复原文本: ${error.message}`);
      } else {
        // 如果还没开始替换就出错了，只需要确保选择状态正确
        const currentDoc = editor.value.state.doc;
        const docSize = currentDoc.content.size;
        
        if (originalSelection.from < docSize && originalSelection.to <= docSize) {
          editor.value.chain()
            .focus()
            .setTextSelection({ from: originalSelection.from, to: originalSelection.to })
            .run();
        }
        alert(`AI 处理失败: ${error.message}`);
      }
    } catch (e) {
      console.error('错误处理时出错:', e);
      alert(`AI 处理失败: ${error.message}`);
    }
  } finally {
    isAiProcessing.value = false;
  }
};

const setLink = () => {
  const previousUrl = editor.value.getAttributes('link').href;
  const url = window.prompt('请输入链接地址', previousUrl);
  if (url === null) return;
  if (url === '') {
    editor.value.chain().focus().extendMarkRange('link').unsetLink().run();
    return;
  }
  editor.value.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
};

const addImage = () => {
  const url = window.prompt('请输入图片链接');
  if (url) {
    editor.value.chain().focus().setImage({ src: url }).run();
  }
};

watch(() => props.modelValue, (newValue) => {
    const isSame = editor.value.getHTML() === newValue;
    if (isSame) return;
    editor.value.commands.setContent(newValue, false);
});
</script>

<style scoped>
.toolbar-button {
    @apply px-2 py-1 text-sm rounded transition-colors duration-200 m-0.5;
}
.toolbar-button:hover { @apply bg-gray-100; }
.is-active { @apply bg-blue-500 text-white; }
.toolbar-button:disabled { @apply opacity-40 cursor-not-allowed; }

.divider {
    @apply w-px h-5 bg-gray-200 mx-2;
}
.tiptap { min-height: 150px; }

/* .tiptap is the class tiptap adds to the editor content, 
   we can style the placeholder through it */
.tiptap p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: #adb5bd;
  pointer-events: none;
  height: 0;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-8px);
  }
  60% {
    transform: translateY(-4px);
  }
}

.animate-bounce {
  animation: bounce 1s infinite;
}
</style> 