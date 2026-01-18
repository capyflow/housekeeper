<template>
  <div class="shareboard-page">
    <div class="page-header">
      <div class="header-left">
        <h2>共享看板管理</h2>
        <p class="subtitle">管理所有的共享剪切板内容</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary" @click="openCreateModal">
          <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          创建看板
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

    <div v-else-if="boards.length === 0" class="empty-state">
      <svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
      </svg>
      <h3>暂无看板</h3>
      <p>点击上方"创建看板"按钮开始创建</p>
    </div>

    <div v-else class="boards-grid">
      <div v-for="board in boards" :key="board.id" class="board-card">
        <div class="board-header">
          <h3 class="board-title">{{ board.board_name }}</h3>
          <div class="board-actions">
            <button class="icon-btn" @click="openEditModal(board)" title="编辑">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
              </svg>
            </button>
            <button class="icon-btn danger" @click="confirmDelete(board)" title="删除">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
              </svg>
            </button>
          </div>
        </div>

        <div class="board-content">
          <p class="content-text">{{ board.content }}</p>
        </div>

        <div class="board-footer">
          <div class="board-meta">
            <span class="meta-item">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
              </svg>
              {{ board.device_name }}
            </span>
            <span class="meta-item">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
              </svg>
              {{ board.owner }}
            </span>
          </div>
          <div class="board-time">
            {{ formatDate(board.update_time) }}
          </div>
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

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ modalMode === 'create' ? '创建看板' : '编辑看板' }}</h3>
          <button class="close-btn" @click="closeModal">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <form class="modal-form" @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="board_name">看板名称 *</label>
            <input
              id="board_name"
              v-model="formData.board_name"
              type="text"
              class="input"
              placeholder="请输入看板名称"
              required
            />
          </div>

          <div class="form-group">
            <label for="content">内容 *</label>
            <textarea
              id="content"
              v-model="formData.content"
              class="input textarea"
              placeholder="请输入内容"
              rows="6"
              required
            ></textarea>
          </div>

          <div class="form-group">
            <label for="device_name">设备名称</label>
            <input
              id="device_name"
              v-model="formData.device_name"
              type="text"
              class="input"
              placeholder="请输入设备名称"
            />
          </div>

          <div class="form-group">
            <label for="owner">所有者</label>
            <input
              id="owner"
              v-model="formData.owner"
              type="text"
              class="input"
              placeholder="请输入所有者"
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

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>确认删除</h3>
        </div>
        <p class="modal-content">确定要删除看板 "{{ boardToDelete?.board_name }}" 吗？此操作不可恢复。</p>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { shareBoardApi } from '@/api/shareboard'
import type { ShareBoard } from '@/types'

const loading = ref(false)
const boards = ref<ShareBoard[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(9)

const filters = reactive({
  owner: ''
})

const showModal = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const submitting = ref(false)

const showDeleteModal = ref(false)
const boardToDelete = ref<ShareBoard | null>(null)
const deleting = ref(false)

const formData = reactive({
  id: '',
  board_name: '',
  content: '',
  device_name: '',
  owner: ''
})

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

const fetchBoards = async () => {
  loading.value = true
  try {
    const res = await shareBoardApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      owner: filters.owner || undefined
    })
    boards.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('Failed to fetch boards:', error)
    alert('加载数据失败')
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  currentPage.value = page
  fetchBoards()
}

const handleFilterChange = () => {
  currentPage.value = 1
  fetchBoards()
}

const resetFilters = () => {
  filters.owner = ''
  currentPage.value = 1
  fetchBoards()
}

const openCreateModal = () => {
  modalMode.value = 'create'
  formData.id = ''
  formData.board_name = ''
  formData.content = ''
  formData.device_name = ''
  formData.owner = ''
  showModal.value = true
}

const openEditModal = (board: ShareBoard) => {
  modalMode.value = 'edit'
  formData.id = board.id
  formData.board_name = board.board_name
  formData.content = board.content
  formData.device_name = board.device_name
  formData.owner = board.owner
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (modalMode.value === 'create') {
      await shareBoardApi.create({
        board_name: formData.board_name,
        content: formData.content,
        device_name: formData.device_name,
        owner: formData.owner
      })
      alert('创建成功')
    } else {
      await shareBoardApi.update({
        id: formData.id,
        board_name: formData.board_name,
        content: formData.content,
        device_name: formData.device_name
      })
      alert('更新成功')
    }
    closeModal()
    fetchBoards()
  } catch (error) {
    console.error('Failed to submit:', error)
    alert('操作失败')
  } finally {
    submitting.value = false
  }
}

const confirmDelete = (board: ShareBoard) => {
  boardToDelete.value = board
  showDeleteModal.value = true
}

const handleDelete = async () => {
  if (!boardToDelete.value) return

  deleting.value = true
  try {
    await shareBoardApi.delete(boardToDelete.value.id)
    alert('删除成功')
    showDeleteModal.value = false
    fetchBoards()
  } catch (error) {
    console.error('Failed to delete:', error)
    alert('删除失败')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchBoards()
})
</script>

<style scoped>
.shareboard-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.filters {
  display: flex;
  gap: 12px;
}

.filter-input {
  max-width: 300px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: var(--text-secondary);
  margin: 0 auto 16px;
}

.empty-state h3 {
  font-size: 18px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-state p {
  color: var(--text-secondary);
}

.boards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.board-card {
  background-color: var(--bg-primary);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.board-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.board-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.board-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.board-actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background-color: transparent;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: var(--bg-secondary);
  color: var(--primary-color);
}

.icon-btn.danger:hover {
  color: var(--danger-color);
}

.icon-btn .icon {
  width: 18px;
  height: 18px;
}

.board-content {
  flex: 1;
  min-height: 80px;
}

.content-text {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  word-break: break-word;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.board-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.board-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-secondary);
}

.meta-item .icon {
  width: 14px;
  height: 14px;
}

.board-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
}

.page-info {
  font-size: 14px;
  color: var(--text-secondary);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.modal-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background-color: transparent;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
}

.close-btn .icon {
  width: 20px;
  height: 20px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.textarea {
  resize: vertical;
  min-height: 120px;
  font-family: inherit;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}

.modal-sm {
  max-width: 400px;
}

.modal-content {
  margin: 20px 0;
  color: var(--text-secondary);
  line-height: 1.6;
}
</style>
