<template>
  <div class="dashboard">
    <div class="content-grid">
      <div class="card recent-notes">
        <div class="section-header">
          <h3 class="card-title">最近笔记</h3>
          <button class="link-btn" @click="$router.push('/notes')">全部笔记</button>
        </div>
        <div v-if="loadingNotes" class="loading">加载中...</div>
        <div v-else-if="recentNotes.length === 0" class="empty">暂无数据</div>
        <div v-else class="note-list">
          <div
            v-for="note in recentNotes"
            :key="note.id"
            class="note-item"
            @click="viewNote(note)"
          >
            <div class="note-info">
              <h4 class="note-title">{{ note.title || '未命名笔记' }}</h4>
              <p class="note-meta">
                <span>{{ note.owner || '-' }}</span>
                <span class="separator">•</span>
                <span>{{ formatNoteDate(note.modify_ts || note.update_ts) }}</span>
              </p>
              <p class="note-preview">{{ getPreviewText(note.content) || '暂无内容' }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="card recent-boards">
        <div class="section-header">
          <h3 class="card-title">看板</h3>
          <button class="link-btn" @click="$router.push('/shareboard')">全部看板</button>
        </div>
        <div v-if="loadingBoards" class="loading">加载中...</div>
        <div v-else-if="recentBoards.length === 0" class="empty">暂无数据</div>
        <div v-else class="board-list">
          <div
            v-for="board in recentBoards"
            :key="board.id"
            class="board-item"
            @click="viewBoard(board)"
          >
            <div class="board-info">
              <h4 class="board-name">{{ board.board_name }}</h4>
              <p class="board-meta">
                <span>{{ board.device_name }}</span>
                <span class="separator">•</span>
                <span>{{ formatBoardDate(board.create_time) }}</span>
              </p>
            </div>
            <div class="board-owner">{{ board.owner }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { noteApi } from '@/api/note'
import { shareBoardApi } from '@/api/shareboard'
import type { Note, ShareBoard } from '@/types'

const router = useRouter()

const loadingNotes = ref(false)
const loadingBoards = ref(false)
const recentNotes = ref<Note[]>([])
const recentBoards = ref<ShareBoard[]>([])

const formatBoardDate = (timestamp: number) => {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatNoteDate = (timestamp: number) => {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
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

const viewNote = (note: Note) => {
  router.push('/notes')
}

const viewBoard = (board: ShareBoard) => {
  router.push('/shareboard')
}

const fetchRecentNotes = async () => {
  loadingNotes.value = true
  try {
    const res = await noteApi.list({ page: 1, page_size: 6 })
    recentNotes.value = res.list || []
  } catch (error) {
    console.error('Failed to fetch notes:', error)
  } finally {
    loadingNotes.value = false
  }
}

const fetchRecentBoards = async () => {
  loadingBoards.value = true
  try {
    const res = await shareBoardApi.list({ page: 1, page_size: 5 })
    recentBoards.value = res.list || []
  } catch (error) {
    console.error('Failed to fetch boards:', error)
  } finally {
    loadingBoards.value = false
  }
}

onMounted(() => {
  fetchRecentNotes()
  fetchRecentBoards()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.content-grid {
  display: grid;
  grid-template-columns: 3fr 1fr;
  gap: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.link-btn {
  font-size: 13px;
  color: var(--text-secondary);
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  transition: all 0.2s ease;
}

.link-btn:hover {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.card-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.3px;
}

.loading,
.empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
  font-size: 15px;
}

.recent-notes {
  min-height: 400px;
}

.note-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.note-item {
  padding: 18px;
  background-color: var(--bg-secondary);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
}

.note-item:hover {
  background-color: var(--bg-primary);
  border-color: var(--primary-color);
  transform: translateX(6px);
  box-shadow: var(--shadow-md);
}

.note-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.note-meta {
  font-size: 13px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.note-preview {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.board-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.board-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background-color: var(--bg-secondary);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
}

.board-item:hover {
  background-color: var(--bg-primary);
  border-color: var(--primary-color);
  transform: translateX(6px);
  box-shadow: var(--shadow-md);
}

.board-info {
  flex: 1;
}

.board-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.board-meta {
  font-size: 13px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 10px;
}

.separator {
  color: var(--border-color);
}

.board-owner {
  font-size: 13px;
  font-weight: 500;
  color: var(--primary-color);
  padding: 6px 14px;
  background-color: var(--primary-light);
  border-radius: 8px;
}

.action-btn span {
  font-size: 15px;
  font-weight: 600;
}

@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard {
    gap: 20px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .stat-card {
    padding: 20px 16px;
    gap: 12px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
  }

  .stat-icon .icon {
    width: 24px;
    height: 24px;
  }

  .stat-label {
    font-size: 12px;
  }

  .stat-value {
    font-size: 24px;
  }

  .content-grid {
    gap: 16px;
  }

  .card-title {
    font-size: 18px;
    margin-bottom: 16px;
  }

  .board-item {
    padding: 14px;
  }

  .board-name {
    font-size: 15px;
  }

  .board-meta {
    font-size: 12px;
  }

  .action-btn {
    padding: 16px 14px;
  }

  .action-btn span {
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .stat-card {
    padding: 16px;
  }

  .stat-card::before {
    height: 3px;
  }

  .stat-value {
    font-size: 22px;
  }

  .card {
    padding: 16px;
  }

  .board-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .board-owner {
    align-self: flex-end;
  }
}
</style>
