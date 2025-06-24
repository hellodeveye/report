<template>
  <div class="border border-gray-300 rounded-lg">
    <div v-if="editor" class="flex items-center flex-wrap p-2 border-b bg-gray-50 rounded-t-lg">
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
        <button @click="toggleAiDropdown" :disabled="editor.state.selection.empty" class="toolbar-button flex items-center">
          AI 优化
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
        </button>
        <div v-if="isAiDropdownOpen" class="absolute z-10 mt-1 bg-white rounded-md shadow-lg border w-48 max-h-60 overflow-y-auto">
          <ul class="py-1">
            <li v-for="option in aiOptions" :key="option">
              <a href="#" @click.prevent="applyAiAction(option)" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">{{ option }}</a>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <EditorContent :editor="editor" />
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

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '请输入...' },
});

const emit = defineEmits(['update:modelValue']);

const isAiDropdownOpen = ref(false);
const aiOptions = [
  '重构', '博客化', '提取要点', '改写', '缩短', '扩写',
  '总结', '简化', '修正拼写', '继续写作', '使用激动语气',
  '添加表情 🙂', '去除表情', '翻译成瑞典语', '翻译成德语', '一句话总结'
];

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
    if (editor.value.state.selection.empty) return;
    isAiDropdownOpen.value = !isAiDropdownOpen.value;
};

const applyAiAction = (action) => {
  isAiDropdownOpen.value = false;
  const { from, to } = editor.value.state.selection;
  const text = editor.value.state.doc.textBetween(from, to, ' ');

  if (text) {
    const optimizedText = `[${action}] ${text} ✨`;
    editor.value.chain().focus().deleteRange({ from, to }).insertContent(optimizedText).run();
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
.toolbar-button:hover { @apply bg-gray-200; }
.is-active { @apply bg-gray-300 text-black; }
.toolbar-button:disabled { @apply opacity-40 cursor-not-allowed; }

.divider {
    @apply w-px h-5 bg-gray-300 mx-2;
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
</style> 