<template>
  <div class="border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50/50 dark:bg-gray-800/40 relative">
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
import { watch } from 'vue';

const props = defineProps({
  content: { type: String, default: '' },
});

const editor = useEditor({
  content: props.content,
  editable: false, // 设为只读
  extensions: [
    StarterKit,
    Underline,
    Highlight,
    Link.configure({ openOnClick: true }), // 允许点击链接
    Subscript,
    Superscript,
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Image,
  ],
  editorProps: {
    attributes: { 
      class: 'prose max-w-none p-3 focus:outline-none min-h-[50px] text-sm text-gray-800 dark:text-gray-100' 
    },
  },
});

// 监听 content 变化并更新编辑器
watch(() => props.content, (newContent) => {
  if (editor.value && newContent !== editor.value.getHTML()) {
    editor.value.commands.setContent(newContent);
  }
});
</script>

<style scoped>
/* 为只读编辑器调整样式 */
:deep(.ProseMirror) {
  background-color: transparent;
}

:deep(.ProseMirror p) {
  margin: 0.5em 0;
}

:deep(.ProseMirror h1, .ProseMirror h2, .ProseMirror h3) {
  margin: 0.8em 0 0.4em 0;
}

:deep(.ProseMirror ul, .ProseMirror ol) {
  margin: 0.5em 0;
  padding-left: 1.2em;
}

:deep(.ProseMirror li) {
  margin: 0.2em 0;
}

:deep(.ProseMirror img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.375rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

:deep(.ProseMirror a) {
  color: #4f46e5;
  text-decoration: underline;
}

:deep(.ProseMirror a:hover) {
  color: #3730a3;
}

:deep(.ProseMirror code) {
  background-color: rgba(156, 163, 175, 0.2);
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
}

:deep(.ProseMirror mark) {
  background-color: #fef08a;
  padding: 0.1em 0.2em;
  border-radius: 0.25rem;
}
</style> 