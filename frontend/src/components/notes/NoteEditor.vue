<template>
  <div class="note-editor">
    <div ref="editorRef" class="vditor-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import type { OpenFile } from '@/store/notes'
import { useThemeStore } from '@/store/theme'

declare global {
  interface Window {
    mermaid: {
      render(id: string, text: string): Promise<{ svg: string }>
      initialize(config: Record<string, unknown>): void
    }
  }
}

interface Props {
  file: OpenFile
  padding?: number
}

const props = withDefaults(defineProps<Props>(), {
  padding: 40,
})

const emit = defineEmits<{
  (e: 'update:content', content: string): void
  (e: 'blur'): void
}>()

const editorRef = ref<HTMLDivElement>()
const themeStore = useThemeStore()
let editor: Vditor | null = null
let isUpdating = false
let observer: MutationObserver | null = null
let enhanceTimeout: ReturnType<typeof setTimeout> | null = null

function applyPadding(value: number) {
  const container = editorRef.value
  if (!container) return
  container.style.setProperty('--editor-padding-h', value + 'px')
  container.style.setProperty('--editor-padding-v', value > 0 ? '20px' : '4px')
  container.classList.toggle('compact', value === 0)
}

function getTheme() {
  return document.documentElement.getAttribute('data-theme') || 'dark'
}

let mermaidLoading: Promise<void> | null = null

function ensureMermaid(): Promise<boolean> {
  if (window.mermaid) return Promise.resolve(true)
  if (mermaidLoading) return mermaidLoading.then(() => true)
  mermaidLoading = new Promise<void>((resolve, reject) => {
    const script = document.createElement('script')
    script.src = 'https://unpkg.com/vditor@3.11.2/dist/js/mermaid/mermaid.min.js'
    script.onload = () => resolve()
    script.onerror = () => reject()
    document.head.appendChild(script)
  })
  return mermaidLoading.then(() => {
    window.mermaid.initialize({
      theme: getTheme() === 'dark' ? 'dark' : 'default',
      sequence: { showSequenceNumbers: true },
    })
    return true
  }).catch(() => false)
}

function enhanceIRContent() {
  const container = editorRef.value
  if (!container) return
  const irElement = container.querySelector('.vditor-ir')
  if (!irElement) return

  numberHeadings(irElement)
  enhanceCodeBlocks(irElement)
  ensureMermaid().then(loaded => {
    if (loaded) renderMermaidInIR(irElement)
  })
}

function numberHeadings(container: Element) {
  const counters = [0, 0, 0, 0, 0, 0]
  const headings = container.querySelectorAll('h1, h2, h3, h4, h5, h6')

  headings.forEach(h => {
    const level = parseInt(h.tagName[1]) - 1
    counters[level]++
    for (let i = level + 1; i < 6; i++) counters[i] = 0

    const num = counters.slice(0, level + 1).join('.')
    const existing = h.querySelector('.vditor-heading-number')
    if (existing) {
      existing.textContent = num + '. '
      return
    }
    const prefix = document.createElement('span')
    prefix.className = 'vditor-heading-number'
    prefix.textContent = num + '. '
    const marker = h.querySelector('.vditor-ir__marker--heading')
    if (marker) {
      marker.after(prefix)
    } else {
      h.insertBefore(prefix, h.firstChild)
    }
  })
}

function renderMermaidInIR(container: Element) {
  if (!window.mermaid) return
  container.querySelectorAll('.language-mermaid').forEach(el => {
    if (el.closest('[data-mermaid-rendered]')) return
    const pre = el.closest('pre')
    if (!pre) return
    pre.setAttribute('data-mermaid-rendered', 'true')
    const theme = getTheme() === 'dark' ? 'dark' : 'default'
    try {
      window.mermaid.render(
        'mermaid-' + Date.now() + '-' + Math.random().toString(36).slice(2, 8),
        el.textContent || ''
      ).then(({ svg }) => {
        pre.innerHTML = svg
      })
    } catch {
      // ignore render errors
    }
  })
}

function enhanceCodeBlocks(container: Element) {
  container.querySelectorAll('pre > code').forEach(code => {
    const pre = code.parentElement
    if (!pre || pre.closest('[data-code-enhanced]')) return
    if (code.classList.contains('language-mermaid') ||
        code.classList.contains('language-flowchart') ||
        code.classList.contains('language-echarts') ||
        code.classList.contains('language-mindmap') ||
        code.classList.contains('language-plantuml') ||
        code.classList.contains('language-graphviz')) {
      return
    }
    pre.setAttribute('data-code-enhanced', 'true')
    pre.style.position = 'relative'

    const lang = Array.from(code.classList)
      .find(c => c.startsWith('language-'))
      ?.replace('language-', '') || ''

    if (lang) {
      const langLabel = document.createElement('span')
      langLabel.className = 'vditor-code-lang'
      langLabel.textContent = lang
      pre.appendChild(langLabel)
    }

    const copyBtn = document.createElement('button')
    copyBtn.className = 'vditor-code-copy'
    copyBtn.textContent = '复制'
    copyBtn.addEventListener('click', (e) => {
      e.stopPropagation()
      const text = code.textContent || ''
      navigator.clipboard.writeText(text).then(() => {
        copyBtn.textContent = '已复制'
        setTimeout(() => { copyBtn.textContent = '复制' }, 1500)
      })
    })
    pre.appendChild(copyBtn)
  })
}

function setupObserver() {
  if (!editorRef.value) return
  const irElement = editorRef.value.querySelector('.vditor-ir')
  if (!irElement) return

  observer = new MutationObserver(() => {
    if (enhanceTimeout) clearTimeout(enhanceTimeout)
    enhanceTimeout = setTimeout(enhanceIRContent, 50)
  })
  observer.observe(irElement, {
    childList: true,
    subtree: true,
    characterData: true,
  })
}

function createEditor() {
  if (!editorRef.value) return

  const vditorTheme = themeStore.isDark ? 'dark' : 'classic'
  const contentTheme = themeStore.isDark ? 'dark' : 'light'

  editor = new Vditor(editorRef.value, {
    mode: 'ir',
    theme: vditorTheme,
    toolbar: [
      'emoji', 'headings', 'bold', 'italic', 'strike', '|',
      'line', 'quote', 'list', 'ordered-list', 'check', '|',
      'code', 'inline-code', 'table', '|',
      'undo', 'redo', '|',
      'outline', 'edit-preview',
    ],
    outline: {
      enable: true,
      position: 'right',
    },
    tab: '    ',
    cache: { enable: false },
    width: '100%',
    height: '100%',
    preview: {
      theme: {
        current: contentTheme,
        path: 'https://unpkg.com/vditor@3.11.2/dist/css/content-theme',
      },
      markdown: {
        // @ts-ignore - Vditor types may be outdated
        codeBlockPreview: true,
      },
      hljs: {
        enable: true,
        lineNumber: true,
        style: themeStore.isDark ? 'github-dark' : 'github',
      },
    },
    input: (value: string) => {
      if (!isUpdating) {
        emit('update:content', value)
      }
    },
    blur: () => {
      emit('blur')
    },
    after: () => {
      if (editor && props.file) {
        isUpdating = true
        editor.setValue(props.file.content)
        isUpdating = false
        applyPadding(props.padding)
      }
      setupObserver()
      setTimeout(enhanceIRContent, 100)
    },
  })
}

function destroyEditor() {
  if (observer) {
    observer.disconnect()
    observer = null
  }
  if (enhanceTimeout) {
    clearTimeout(enhanceTimeout)
    enhanceTimeout = null
  }
  if (editor) {
    editor.destroy()
    editor = null
  }
}

function toggleOutline() {
  if (!editorRef.value) return
  const btn = editorRef.value.querySelector('[data-type="outline"]') as HTMLElement
  if (btn) btn.click()
}

defineExpose({ toggleOutline })

watch(
  () => props.file?.path,
  (newPath, oldPath) => {
    if (newPath !== oldPath && editor) {
      isUpdating = true
      editor.setValue(props.file?.content || '')
      isUpdating = false
    }
  }
)

watch(
  () => props.padding,
  (val) => applyPadding(val)
)

let unwatchTheme: (() => void) | null = null

onMounted(() => {
  createEditor()

  unwatchTheme = watch(
    () => themeStore.isDark,
    (isDark) => {
      if (!editor) return
      const vditorTheme = isDark ? 'dark' : 'classic'
      const contentTheme = isDark ? 'dark' : 'light'
      const codeTheme = isDark ? 'github-dark' : 'github'
      editor.setTheme(vditorTheme, contentTheme, codeTheme)
      if (window.mermaid) {
        window.mermaid.initialize({
          theme: isDark ? 'dark' : 'default',
          sequence: { showSequenceNumbers: true },
        })
      }
      applyPadding(props.padding)
      setTimeout(enhanceIRContent, 150)
    }
  )
})

onBeforeUnmount(() => {
  if (unwatchTheme) {
    unwatchTheme()
    unwatchTheme = null
  }
  destroyEditor()
})
</script>

<style scoped>
.note-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.vditor-container {
  flex: 1;
  overflow: hidden;
}

:deep(.vditor) {
  border: none !important;
  border-radius: 0 !important;
}

:deep(.vditor-toolbar) {
  border-bottom: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
  padding: 4px 8px !important;
  min-height: 36px !important;
}

:deep(.vditor-toolbar__item) {
  background: transparent !important;
}

:deep(.vditor-toolbar__item:hover) {
  background: var(--bg-hover) !important;
}

:deep(.vditor-toolbar__item) svg {
  fill: var(--text-secondary) !important;
}

:deep(.vditor-toolbar__item--current) svg {
  fill: var(--accent) !important;
}

:deep(.vditor-content) {
  background: var(--bg-primary) !important;
}

:deep(.vditor-ir) {
  background: var(--bg-primary) !important;
  color: var(--text-primary) !important;
  font-size: 15px !important;
  line-height: 1.8 !important;
  padding: var(--editor-padding-v, 20px) var(--editor-padding-h, 40px) !important;
}

:deep(.vditor-ir:focus) {
  outline: none !important;
}

:deep(.vditor-ir pre.vditor-reset) {
  padding: var(--editor-padding-v, 20px) var(--editor-padding-h, 40px) !important;
}

.compact :deep(.vditor-ir p),
.compact :deep(.vditor-ir h1),
.compact :deep(.vditor-ir h2),
.compact :deep(.vditor-ir h3),
.compact :deep(.vditor-ir h4),
.compact :deep(.vditor-ir h5),
.compact :deep(.vditor-ir h6),
.compact :deep(.vditor-ir li),
.compact :deep(.vditor-ir blockquote),
.compact :deep(.vditor-ir pre),
.compact :deep(.vditor-ir table),
.compact :deep(.vditor-ir hr) {
  margin-left: 0 !important;
  margin-right: 0 !important;
}

:deep(.vditor-outline) {
  border-left: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
  font-size: 12px !important;
}

:deep(.vditor-outline__item) {
  color: var(--text-secondary) !important;
}

:deep(.vditor-outline__item:hover) {
  color: var(--text-primary) !important;
  background: var(--bg-hover) !important;
}

:deep(.vditor-outline__item--current) {
  color: var(--accent) !important;
}

:deep(*::-webkit-scrollbar) {
  width: 6px !important;
  height: 6px !important;
}

:deep(*::-webkit-scrollbar-track) {
  background: transparent !important;
}

:deep(*::-webkit-scrollbar-thumb) {
  background: rgba(128, 128, 128, 0.3) !important;
  border-radius: 3px !important;
}

:deep(*::-webkit-scrollbar-thumb:hover) {
  background: rgba(128, 128, 128, 0.5) !important;
}

/* ===== Vditor 暗黑模式全覆盖 ===== */

/* 编辑区 tooltip */
:deep(.vditor-tip) {
  background: var(--bg-tertiary) !important;
  border-color: var(--border-color) !important;
  color: var(--text-primary) !important;
}

/* 下拉菜单 / 面板 */
:deep(.vditor-panel),
:deep(.vditor-menu--panel) {
  background: var(--bg-secondary) !important;
  border-color: var(--border-color) !important;
  box-shadow: var(--shadow) !important;
}
:deep(.vditor-panel--item) {
  color: var(--text-primary) !important;
}
:deep(.vditor-panel--item:hover) {
  background: var(--bg-hover) !important;
}

/* 提示 */
:deep(.vditor-hint) {
  background: var(--bg-secondary) !important;
  border-color: var(--border-color) !important;
}
:deep(.vditor-hint--item) {
  color: var(--text-primary) !important;
}
:deep(.vditor-hint--item:hover) {
  background: var(--bg-hover) !important;
}
:deep(.vditor-hint--item--current) {
  background: var(--bg-selected) !important;
}

/* 预览区 */
:deep(.vditor-preview) {
  background: var(--bg-primary) !important;
}
:deep(.vditor-reset) {
  color: var(--text-primary) !important;
  background: var(--bg-primary) !important;
}

/* 预览区代码块背景 */
:deep(.vditor-reset pre) {
  background: var(--bg-tertiary) !important;
}
:deep(.vditor-reset code) {
  background: transparent !important;
}

/* 预览表格 */
:deep(.vditor-reset table tr) {
  background: var(--bg-primary) !important;
  border-color: var(--border-color) !important;
}
:deep(.vditor-reset table td),
:deep(.vditor-reset table th) {
  border-color: var(--border-color) !important;
}
:deep(.vditor-reset table tr:nth-child(2n)) {
  background: var(--bg-secondary) !important;
}

/* 预览区链接 */
:deep(.vditor-reset a) {
  color: var(--accent) !important;
}

/* 预览区 blockquote */
:deep(.vditor-reset blockquote) {
  color: var(--text-secondary) !important;
  border-left-color: var(--border-color) !important;
}

/* 预览区代码块行号 */
:deep(.vditor-linenumber__rows) {
  border-right: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
}
:deep(.vditor-linenumber__rows > span::before) {
  color: var(--text-muted) !important;
  opacity: 0.5 !important;
}

/* IR 模式代码块 */
:deep(.vditor-ir pre) {
  background: var(--bg-tertiary) !important;
  border: 1px solid var(--border-color) !important;
}

/* IR 模式行内代码 */
:deep(.vditor-ir code) {
  background: var(--bg-tertiary) !important;
  color: var(--text-primary) !important;
}

/* IR 模式标题标记颜色 */
:deep(.vditor-ir__marker--heading) {
  color: var(--accent) !important;
  opacity: 0.6 !important;
}

/* IR 模式粗体/斜体标记 */
:deep(.vditor-ir__marker--bi) {
  color: var(--accent) !important;
  opacity: 0.5 !important;
}

/* IR 模式链接 */
:deep(.vditor-ir__marker--link),
:deep(.vditor-ir__link) {
  color: var(--accent) !important;
}

/* IR 模式引用块 */
:deep(.vditor-ir blockquote) {
  border-left-color: var(--accent) !important;
  opacity: 0.9 !important;
}

/* IR 模式分隔线 */
:deep(.vditor-ir hr) {
  border-color: var(--border-color) !important;
}

/* IR 模式列表标记 */
:deep(.vditor-ir li) {
  color: var(--text-primary) !important;
}

/* IR 模式表格 */
:deep(.vditor-ir table tr) {
  background: var(--bg-primary) !important;
}
:deep(.vditor-ir table td),
:deep(.vditor-ir table th) {
  border-color: var(--border-color) !important;
}

/* 预览区 Mermaid 背景 */
:deep(.mermaid svg) {
  background: transparent !important;
}

/* 预览区复制按钮 */
:deep(.vditor-copy span) {
  opacity: 0.6 !important;
  transition: opacity 0.2s !important;
}
:deep(pre:hover .vditor-copy span) {
  opacity: 1 !important;
}

/* 预览区复制按钮 SVG 图标颜色 */
:deep(.vditor-copy svg) {
  fill: var(--text-muted) !important;
}

/* 大纲面板 dark 适配 */
:deep(.vditor-outline__title) {
  color: var(--text-secondary) !important;
  border-bottom-color: var(--border-color) !important;
}

/* 大纲搜索 */
:deep(.vditor-outline__input) {
  background: var(--bg-primary) !important;
  border-color: var(--border-input) !important;
  color: var(--text-primary) !important;
}

/* 大纲 empty 状态 */
:deep(.vditor-outline__empty) {
  color: var(--text-muted) !important;
}
</style>
