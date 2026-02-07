<template>
  <div class="dashboard">
    <div class="content-grid">
      <div class="card recent-boards">
        <h3 class="card-title">最近的看板</h3>
        <div v-if="loading" class="loading">加载中...</div>
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
                <span>{{ formatDate(board.create_time) }}</span>
              </p>
            </div>
            <div class="board-owner">{{ board.owner }}</div>
          </div>
        </div>
      </div>

      <div class="card quick-actions">
        <h3 class="card-title">快捷操作</h3>
        <div class="action-list">
          <button class="action-btn" @click="$router.push('/shareboard')">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            <span>创建新看板</span>
          </button>
          <button class="action-btn" @click="$router.push('/shareboard')">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <span>查看所有看板</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { shareBoardApi } from '@/api/shareboard'
import type { ShareBoard } from '@/types'

const router = useRouter()

const loading = ref(false)
const recentBoards = ref<ShareBoard[]>([])

const formatDate = (timestamp: number) => {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const viewBoard = (board: ShareBoard) => {
  router.push('/shareboard')
}

const fetchRecentBoards = async () => {
  loading.value = true
  try {
    const res = await shareBoardApi.list({ page: 1, page_size: 5 })
    recentBoards.value = res.list || []
  } catch (error) {
    console.error('Failed to fetch boards:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
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
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.card-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 24px;
  letter-spacing: -0.3px;
}

.recent-boards {
  min-height: 400px;
}

.loading,
.empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
  font-size: 15px;
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

.quick-actions {
  min-height: 400px;
}

.action-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px;
  background-color: var(--bg-secondary);
  border-radius: 12px;
  color: var(--text-primary);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
  border: 1px solid var(--border-color);
  position: relative;
  overflow: hidden;
}

.action-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(to bottom, var(--primary-color), var(--accent-purple));
  transform: scaleY(0);
  transition: transform 0.3s;
}

.action-btn:hover {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--accent-purple) 100%);
  color: white;
  border-color: transparent;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.action-btn:hover::before {
  transform: scaleY(1);
}

.action-btn .icon {
  width: 22px;
  height: 22px;
  transition: transform 0.3s;
}

.action-btn:hover .icon {
  transform: scale(1.1);
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
