<template>
  <div class="notes-page">
    <div class="page-header">
      <div class="header-left">
        <h2>笔记管理</h2>
        <p class="subtitle">管理所有的笔记内容</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary" @click="openCreateModal">
          <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          创建笔记
        </button>
      </div>
    </div>

    <div class="filters">
      <input
        v-model="filters.owner"
        type="text"
        class="input filter-input"
        placeholder="按所有者筛选..."
        @input="handleFilterChange"
      />
      <button class="btn btn-outline" @click="resetFilters">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        重置
      </button>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="notes.length === 0" class="empty-state">
      <svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
      </svg>
      <h3>暂无数据</h3>
    </div>

    <div v-else class="notes-grid">
      <div v-for="note in notes" :key="note.id" class="note-card" @click="openNote(note)">
        <div class="note-header">
          <h3 class="note-title">{{ note.title || '未命名笔记' }}</h3>
          <div class="note-actions">
            <button class="icon-btn" @click.stop="openEditModal(note)" title="编辑">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
              </svg>
            </button>
            <button class="icon-btn danger" @click.stop="confirmDelete(note)" title="删除">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
              </svg>
            </button>
          </div>
        </div>
        <p class="note-content">{{ getPreviewText(note.content) || '暂无内容' }}</p>
        <div class="note-meta">
          <span class="meta-item">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
            {{ note.owner || '-' }}
          </span>
          <span class="meta-item">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            {{ formatDate(note.modify_ts || note.update_ts) }}
          </span>
        </div>
      </div>
    </div>

    <div v-if="total > pageSize" class="pagination">
      <button
        class="btn btn-outline"
        :disabled="currentPage === 1"
        @click="changePage(currentPage - 1)"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ currentPage }} / {{ totalPages }} 页，共 {{ total }} 条
      </span>
      <button
        class="btn btn-outline"
        :disabled="currentPage === totalPages"
        @click="changePage(currentPage + 1)"
      >
        下一页
      </button>
    </div>

    <div v-if="showViewModal" class="modal-overlay" @click.self="closeViewModal">
      <div class="modal view-modal">
        <div class="modal-header">
          <h3>{{ noteToView?.title || '笔记详情' }}</h3>
          <button class="close-btn" @click="closeViewModal">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        <div class="view-content">
          <div class="view-section">
            <label class="section-label">内容</label>
            <div ref="previewRef" class="content-box vditor-reset"></div>
          </div>
          <div class="view-info">
            <div class="info-item">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
              </svg>
              <span>{{ noteToView?.owner || '-' }}</span>
            </div>
            <div class="info-item">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>创建于 {{ formatFullDate(noteToView?.create_ts) }}</span>
            </div>
            <div class="info-item">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
              <span>更新于 {{ formatFullDate(noteToView?.modify_ts || noteToView?.update_ts) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ modalMode === 'create' ? '创建笔记' : '编辑笔记' }}</h3>
          <button class="close-btn" @click="closeModal">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <form class="modal-form" @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="title">标题 *</label>
            <input
              id="title"
              v-model="formData.title"
              type="text"
              class="input"
              placeholder="请输入笔记标题"
              required
            />
          </div>

          <div class="form-group">
            <div class="editor-header">
              <label for="content">内容 *</label>
            </div>
            <div class="md-editor" :class="{ fullscreen: isFullscreen }">
              <div ref="vditorRef" class="vditor-host"></div>
            </div>
          </div>

          <div class="form-group">
            <label for="cover">封面</label>
            <input
              id="cover"
              v-model="formData.cover"
              type="text"
              class="input"
              placeholder="可选，封面 fid"
            />
          </div>

          <div class="modal-actions">
            <button type="button" class="btn btn-outline" @click="closeModal">
              取消
            </button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              {{ submitting ? '提交中...' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>确认删除</h3>
        </div>
        <p class="modal-content">确定要删除笔记 "{{ noteToDelete?.title || '未命名笔记' }}" 吗？此操作不可恢复。</p>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="showDeleteModal = false">
            取消
          </button>
          <button class="btn btn-danger" :disabled="deleting" @click="handleDelete">
            {{ deleting ? '删除中...' : '确定删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { noteApi } from '@/api/note'
import { useToast } from '@/composables/useToast'
import { useThemeStore } from '@/stores/theme'
import type { Note } from '@/types'

const toast = useToast()
const themeStore = useThemeStore()

const loading = ref(false)
const notes = ref<Note[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const filters = reactive({
  owner: ''
})

const showModal = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const submitting = ref(false)

const showViewModal = ref(false)
const noteToView = ref<Note | null>(null)

const showDeleteModal = ref(false)
const noteToDelete = ref<Note | null>(null)
const deleting = ref(false)

const vditorRef = ref<HTMLDivElement | null>(null)
const previewRef = ref<HTMLDivElement | null>(null)
const vditorInstance = ref<Vditor | null>(null)

const formData = reactive({
  id: '',
  title: '',
  content: '',
  cover: ''
})

const isFullscreen = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const formatDate = (timestamp: number) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`

  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const formatFullDate = (timestamp?: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const fetchNotes = async () => {
  loading.value = true
  try {
    const res = await noteApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      owner: filters.owner || undefined
    })
    notes.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Failed to fetch notes:', error)
    toast.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  currentPage.value = page
  fetchNotes()
}

const handleFilterChange = () => {
  currentPage.value = 1
  fetchNotes()
}

const resetFilters = () => {
  filters.owner = ''
  currentPage.value = 1
  fetchNotes()
}

const openNote = (note: Note) => {
  noteToView.value = note
  showViewModal.value = true
  nextTick(() => {
    renderPreview(note.content || '')
  })
}

const closeViewModal = () => {
  showViewModal.value = false
  noteToView.value = null
  if (previewRef.value) {
    previewRef.value.innerHTML = ''
  }
}

const openCreateModal = () => {
  modalMode.value = 'create'
  formData.id = ''
  formData.title = ''
  formData.content = ''
  formData.cover = ''
  showModal.value = true
  nextTick(() => {
    initVditor('')
  })
}

const openEditModal = (note: Note) => {
  modalMode.value = 'edit'
  formData.id = note.id
  formData.title = note.title || ''
  formData.content = note.content || ''
  formData.cover = note.cover || ''
  showModal.value = true
  nextTick(() => {
    initVditor(formData.content)
  })
}

const closeModal = () => {
  destroyVditor()
  showModal.value = false
  isFullscreen.value = false
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (vditorInstance.value) {
      formData.content = vditorInstance.value.getValue()
    }
    if (!formData.content || !formData.content.trim()) {
      toast.error('请输入笔记内容')
      return
    }
    if (modalMode.value === 'create') {
      await noteApi.create({
        title: formData.title,
        content: formData.content,
        cover: formData.cover || undefined
      })
      toast.success('创建成功')
    } else {
      await noteApi.update({
        id: formData.id,
        title: formData.title,
        content: formData.content,
        cover: formData.cover || undefined
      })
      toast.success('更新成功')
    }
    closeModal()
    fetchNotes()
  } catch (error) {
    console.error('Failed to submit:', error)
    toast.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const confirmDelete = (note: Note) => {
  noteToDelete.value = note
  showDeleteModal.value = true
}

const handleDelete = async () => {
  if (!noteToDelete.value) return

  deleting.value = true
  try {
    await noteApi.delete(noteToDelete.value.id)
    toast.success('删除成功')
    showDeleteModal.value = false
    fetchNotes()
  } catch (error) {
    console.error('Failed to delete:', error)
    toast.error('删除失败')
  } finally {
    deleting.value = false
  }
}

const getPreviewText = (markdown: string) => {
  if (!markdown) return ''
  return markdown
    .replace(/```[\s\S]*?```/g, '')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/!\[.*?\]\(.*?\)/g, '')
    .replace(/\[(.*?)\]\(.*?\)/g, '$1')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/^>\s+/gm, '')
    .replace(/^[-*+]\s+/gm, '')
    .replace(/[*_~]/g, '')
    .trim()
}

const getHljsStyle = () => (themeStore.theme === 'dark' ? 'github-dark' : 'github')

const renderPreview = (markdown: string) => {
  if (!previewRef.value) return
  previewRef.value.innerHTML = ''
  void Vditor.preview(previewRef.value, markdown, {
    mode: themeStore.theme === 'dark' ? 'dark' : 'light',
    hljs: {
      style: getHljsStyle()
    }
  })
}

const initVditor = (value: string) => {
  if (!vditorRef.value) return

  if (vditorInstance.value) {
    vditorInstance.value.setValue(value, true)
    vditorInstance.value.setPreviewMode(isFullscreen.value ? 'both' : 'editor')
    return
  }

  vditorInstance.value = new Vditor(vditorRef.value, {
    height: 520,
    mode: 'sv',
    theme: themeStore.theme === 'dark' ? 'dark' : 'classic',
    cache: {
      enable: false
    },
    value,
    placeholder: '请输入 Markdown 内容',
    toolbar: [
      'emoji',
      'headings',
      'bold',
      'italic',
      'strike',
      '|',
      'line',
      'quote',
      'list',
      'ordered-list',
      'check',
      'outdent',
      'indent',
      '|',
      'code',
      'inline-code',
      'insert-after',
      'insert-before',
      'table',
      'link',
      '|',
      'undo',
      'redo',
      'fullscreen',
      'edit-mode',
      'help'
    ],
    input: (val) => {
      formData.content = val
    },
    esc: () => {
      if (isFullscreen.value) {
        toggleFullscreen()
      }
    },
    ctrlEnter: () => {
      handleSubmit()
    },
    preview: {
      hljs: {
        style: getHljsStyle()
      }
    }
  })

  vditorInstance.value.setPreviewMode(isFullscreen.value ? 'both' : 'editor')
}

const destroyVditor = () => {
  if (!vditorInstance.value) return
  vditorInstance.value.destroy()
  vditorInstance.value = null
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  if (vditorInstance.value) {
    vditorInstance.value.setPreviewMode(isFullscreen.value ? 'both' : 'editor')
  }
}

watch(
  () => themeStore.theme,
  (theme) => {
    if (vditorInstance.value) {
      vditorInstance.value.setTheme(theme === 'dark' ? 'dark' : 'classic')
    }
    if (noteToView.value) {
      renderPreview(noteToView.value.content || '')
    }
  }
)

onMounted(() => {
  fetchNotes()
})
</script>

<style scoped>
.notes-page {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.header-left h2 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 6px;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 14px;
  margin: 0;
}

.filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

.filter-input {
  max-width: 260px;
}

.loading,
.empty-state {
  background-color: var(--bg-primary);
  border-radius: 16px;
  padding: 48px;
  text-align: center;
  border: 1px solid var(--border-color);
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  color: var(--text-secondary);
}

.notes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.note-card {
  background-color: var(--bg-primary);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  transition: all 0.25s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
  overflow: hidden;
}

.note-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary-color), var(--accent-purple), var(--accent-orange));
  transform: scaleX(0);
  transition: transform 0.3s;
}

.note-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary-color);
}

.note-card:hover::before {
  transform: scaleX(1);
}

.note-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.note-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.3;
}

.note-actions {
  display: flex;
  gap: 8px;
}

.note-content {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.6;
  max-height: 4.8em;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
}

.note-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.meta-item .icon {
  width: 14px;
  height: 14px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background-color: transparent;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid transparent;
}

.icon-btn:hover {
  background-color: var(--primary-light);
  color: var(--primary-color);
  border-color: var(--primary-color);
  transform: scale(1.1);
}

.icon-btn.danger:hover {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger-color);
  border-color: var(--danger-color);
}

.icon-btn .icon {
  width: 20px;
  height: 20px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.modal-header h3 {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.close-btn {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background-color: transparent;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  border: 1px solid transparent;
}

.close-btn:hover {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  border-color: var(--border-color);
  transform: scale(1.05);
}

.close-btn .icon {
  width: 22px;
  height: 22px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-group label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.3px;
}

.md-editor {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  background-color: var(--bg-primary);
  position: relative;
  min-height: 520px;
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-xs {
  padding: 6px 10px;
  font-size: 12px;
  border-radius: 8px;
}

.md-editor.fullscreen {
  position: fixed;
  top: 24px;
  left: 24px;
  right: 24px;
  bottom: 24px;
  z-index: 2000;
  border-radius: 16px;
  box-shadow: var(--shadow-xl);
  display: flex;
  flex-direction: column;
}

.vditor-host {
  min-height: 220px;
}

.md-editor.fullscreen .vditor-host {
  flex: 1;
  min-height: 0;
}

.md-editor.fullscreen :deep(.vditor) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.md-editor.fullscreen :deep(.vditor .vditor-content) {
  flex: 1;
  min-height: 0;
}

:root[data-theme='dark'] .md-editor :deep(.vditor) {
  color: #e5e7eb;
}

:root[data-theme='dark'] .md-editor :deep(.vditor .vditor-input),
:root[data-theme='dark'] .md-editor :deep(.vditor-ir__editor),
:root[data-theme='dark'] .md-editor :deep(.vditor-sv__textarea),
:root[data-theme='dark'] .md-editor :deep(.vditor-wysiwyg__editor) {
  color: #e5e7eb;
}

:root[data-theme='dark'] .md-editor :deep(.vditor) {
  background-color: #111827;
  border-color: rgba(255, 255, 255, 0.08);
}

:root[data-theme='dark'] .md-editor :deep(.vditor-toolbar) {
  background-color: #0f172a;
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

:root[data-theme='dark'] .md-editor :deep(.vditor-toolbar .vditor-toolbar__item) {
  color: #cbd5f5;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-toolbar .vditor-toolbar__item:hover) {
  background-color: rgba(255, 255, 255, 0.08);
}

:root[data-theme='dark'] .md-editor :deep(.vditor-panel) {
  background-color: #0f172a;
  border-color: rgba(255, 255, 255, 0.08);
}

:root[data-theme='dark'] .md-editor :deep(.vditor-input::placeholder),
:root[data-theme='dark'] .md-editor :deep(.vditor-sv__textarea::placeholder) {
  color: #94a3b8;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-content) {
  background-color: #111827;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview) {
  background-color: #0b1220;
  color: #e5e7eb;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset) {
  color: #e5e7eb;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset p),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset li),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset blockquote),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset h1),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset h2),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset h3),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset h4),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset h5),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset h6),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset table),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset th),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .vditor-reset td) {
  color: #e5e7eb;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview a) {
  color: #93c5fd;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview pre) {
  background-color: rgba(255, 255, 255, 0.06);
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview code) {
  color: inherit;
  background-color: transparent;
}

:root[data-theme='dark'] .md-editor :deep(.vditor-preview .hljs),
:root[data-theme='dark'] .md-editor :deep(.vditor-preview .hljs *) {
  color: unset;
}

:deep(.modal) {
  max-width: 980px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 12px;
}

.modal-sm {
  max-width: 460px;
}

.modal-content {
  margin: 24px 0;
  color: var(--text-secondary);
  line-height: 1.7;
  font-size: 15px;
}

.view-modal {
  max-width: 700px;
}

.view-content {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.view-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.content-box {
  background-color: var(--bg-secondary);
  border-radius: 12px;
  padding: 16px;
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-primary);
  white-space: pre-wrap;
  word-break: break-word;
}

.view-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-secondary);
  font-size: 13px;
}

.info-item .icon {
  width: 16px;
  height: 16px;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-input {
    max-width: 100%;
  }
}
</style>
