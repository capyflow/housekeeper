<template>
  <div class="checkin-target-page">
    <!-- Hero 区域 -->
    <section class="hero card">
      <div class="hero-text">
        <h2>{{ currentTarget?.title || '打卡任务' }}</h2>
        <p>坚持打卡，让每个日子都闪闪发光 ✨</p>
      </div>
      <div class="hero-actions">
        <button class="btn btn-outline" @click="openCreateModal" :disabled="submitting">新建任务</button>
        <button 
          class="btn btn-primary" 
          @click="handleCheckInToday" 
          :disabled="submitting || !currentTarget"
          :class="{ checked: todayChecked }"
        >
          {{ todayChecked ? '✅ 今日已打卡' : '📅 今日打卡' }}
        </button>
      </div>
    </section>

    <!-- 任务列表 -->
    <section class="target-list-section card">
      <div class="section-header">
        <h3>我的任务</h3>
        <span>{{ targets.length }} 个</span>
      </div>
      
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <div v-else-if="targets.length === 0" class="empty-state">
        <div class="empty-icon">🎯</div>
        <p>还没有打卡任务，创建一个开始吧！</p>
      </div>

      <div v-else class="target-chips">
        <button
          v-for="item in targets"
          :key="item.id"
          class="target-chip"
          :class="{ active: item.id === currentTargetId }"
          @click="currentTargetId = item.id"
        >
          <div class="chip-content">
            <span class="name">{{ item.title || '未命名任务' }}</span>
            <span class="desc">{{ item.description || '暂无描述' }}</span>
          </div>
          <div class="chip-meta">
            <span class="total-days">⏱ {{ item.total_days }}天</span>
            <span class="time-range">{{ formatTimeRange(item.time_range) }}</span>
          </div>
        </button>
      </div>
    </section>

    <!-- 热力图 -->
    <section class="heatmap-section card" v-if="currentTarget">
      <div class="section-header">
        <h3>近一年打卡热力图</h3>
        <div class="legend">
          <span>未打卡</span>
          <i class="lv1"></i>
          <i class="lv2"></i>
          <i class="lv3"></i>
          <i class="lv4"></i>
          <span>天天打</span>
        </div>
      </div>

      <div class="heatmap-wrap">
        <div class="weekday-labels">
          <span>一</span>
          <span>三</span>
          <span>五</span>
        </div>
        <div class="weeks-container">
          <div class="week" v-for="(week, wIdx) in heatmapWeeks" :key="wIdx">
            <button
              v-for="cell in week"
              :key="cell.dateKey"
              class="day"
              :class="[`lv${cell.level}`, { today: cell.isToday, disabled: cell.outOfRange }]"
              :title="cell.tooltip"
              @click="handleCheckInByDate(cell.date)"
              :disabled="cell.disabled"
            ></button>
          </div>
        </div>
      </div>

      <div class="stats-summary">
        <div class="stat-item">
          <strong>{{ currentTarget.total_days }}</strong>
          <span>总打卡天数</span>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <strong>{{ getStreakDays() }}</strong>
          <span>当前连续天数</span>
        </div>
      </div>

      <p class="hint">点击任意日期可补打卡；今天打卡后，今天格子会点亮 🌟</p>
    </section>

    <!-- 创建/编辑弹窗 -->
    <div v-if="modalOpen" class="modal-overlay" @click="modalOpen = false">
      <div class="modal edit-modal" @click.stop>
        <h3>{{ isEditing ? '编辑任务' : '新建打卡任务' }}</h3>
        
        <div class="form-group">
          <label>任务名称 *</label>
          <input 
            class="input" 
            v-model="formData.title" 
            maxlength="50" 
            placeholder="例如：健身 / 阅读 / 学英语" 
          />
        </div>

        <div class="form-group">
          <label>描述</label>
          <textarea 
            class="textarea" 
            v-model="formData.description" 
            rows="3"
            placeholder="添加一些描述，帮助自己记住为什么开始..."
          ></textarea>
        </div>

        <div class="form-group">
          <label>
            <input type="checkbox" v-model="formData.timeRange.enabled" />
            限制打卡时间段 (不勾选则全天可以打卡)
          </label>
          
          <div v-if="formData.timeRange.enabled" class="time-range-inputs">
            <div class="time-input">
              <label>开始时间</label>
              <input 
                type="number" 
                v-model.number="startTimeInput"
                min="0" 
                max="2359" 
                placeholder="HHMM"
              />
            </div>
            <span class="separator">~</span>
            <div class="time-input">
              <label>结束时间</label>
              <input 
                type="number" 
                v-model.number="endTimeInput"
                min="0" 
                max="2359" 
                placeholder="HHMM"
              />
            </div>
          </div>
        </div>

        <div class="modal-actions">
          <button class="btn btn-outline" @click="modalOpen = false">取消</button>
          <button 
            class="btn btn-primary" 
            :disabled="submitting || !formData.title.trim()" 
            @click="handleSubmit"
          >
            {{ isEditing ? '保存' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { checkinApi } from '@/api/checkin'
import type { CheckInTarget, TimeRange } from '@/types'
import { useToast } from '@/composables/useToast'

const { success, error: toastError } = useToast()

// 状态
const loading = ref(false)
const submitting = ref(false)
const targets = ref<CheckInTarget[]>([])
const currentTargetId = ref('')
const modalOpen = ref(false)
const isEditing = ref(false)

// 表单数据
const formData = ref({
  title: '',
  description: '',
  timeRange: {
    enabled: false,
    start: 6 * 60,      // 06:00
    end: 23 * 60 + 59   // 23:59
  } as TimeRange
})

// 时间输入临时变量
const startTimeInput = ref(600)
const endTimeInput = ref(2359)

const DAY_MS = 24 * 60 * 60 * 1000

// 工具函数
const normalizeDayStart = (ts: number) => {
  const d = new Date(ts)
  return new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
}

const dayKey = (ts: number) => {
  const d = new Date(ts)
  const y = d.getFullYear()
  const m = `${d.getMonth() + 1}`.padStart(2, '0')
  const day = `${d.getDate()}`.padStart(2, '0')
  return `${y}-${m}-${day}`
}

const formatTimeRange = (tr: TimeRange) => {
  if (!tr.enabled) return '全天'
  const startHour = Math.floor(tr.start / 60)
  const startMin = tr.start % 60
  const endHour = Math.floor(tr.end / 60)
  const endMin = tr.end % 60
  const format = (h: number, m: number) => `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
  
  if (tr.start > tr.end) {
    return `${format(startHour, startMin)} ~ ${format(endHour, endMin)} (次日)`
  }
  return `${format(startHour, startMin)} ~ ${format(endHour, endMin)}`
}

const getCurrentTarget = computed(() => 
  targets.value.find(t => t.id === currentTargetId.value) || null
)

const countByDay = computed(() => {
  const map: { [key: string]: number } = {}
  const counts = getCurrentTarget.value?.check_in_counts || {}
  Object.entries(counts).forEach(([date, count]) => {
    map[date] = count
  })
  return map
})

const todayChecked = computed(() => {
  const today = dayKey(Date.now())
  return (countByDay.value[today] || 0) > 0
})

// 热力图计算
const heatmapWeeks = computed(() => {
  const now = new Date()
  const todayStart = normalizeDayStart(now.getTime())
  const start = new Date(todayStart - 364 * DAY_MS)
  const startOffset = start.getDay()
  const alignedStart = start.getTime() - startOffset * DAY_MS

  const totalCells = 53 * 7
  const weeks: Array<Array<{ date: number; dateKey: string; level: number; tooltip: string; disabled: boolean; outOfRange: boolean }>> = []

  for (let i = 0; i < totalCells; i++) {
    const ts = alignedStart + i * DAY_MS
    const key = dayKey(ts)
    const count = countByDay.value[key] || 0
    const level = count === 0 ? 0 : Math.min(count, 4) // 最多 4 级
    const inRange = ts >= start.getTime() && ts <= todayStart
    const outOfRange = !inRange
    
    const weekIndex = Math.floor(i / 7)
    if (!weeks[weekIndex]) {
      weeks[weekIndex] = []
    }
    weeks[weekIndex].push({
      date: ts,
      dateKey: key,
      level,
      tooltip: `${key} ${inRange ? `打卡 ${count} 次` : '超出范围'}`,
      disabled: !inRange,
      outOfRange
    })
  }

  return weeks.slice(0, 53)
})

// 连续打卡天数计算
const getStreakDays = () => {
  const today = normalizeDayStart(Date.now())
  let streak = 0
  let current = today
  
  while (true) {
    const key = dayKey(current)
    if ((countByDay.value[key] || 0) > 0) {
      streak++
      current -= DAY_MS
    } else {
      break
    }
  }
  
  return streak
}

// 加载任务列表
const fetchTargets = async () => {
  loading.value = true
  try {
    const res = await checkinApi.list({ page: 1, page_size: 100 })
    targets.value = (res.list || []).sort((a, b) => b.create_time - a.create_time)
    
    if (!currentTargetId.value && targets.value.length > 0) {
      currentTargetId.value = targets.value[0].id
    }
  } catch (e) {
    console.error('fetch targets failed', e)
    toastError('加载任务失败')
  } finally {
    loading.value = false
  }
}

// 打开创建/编辑弹窗
const openCreateModal = () => {
  isEditing.value = false
  formData.value = {
    title: '',
    description: '',
    timeRange: {
      enabled: false,
      start: 6 * 60,
      end: 23 * 60 + 59
    }
  }
  startTimeInput.value = 600
  endTimeInput.value = 2359
  modalOpen.value = true
}

const openEditModal = (target: CheckInTarget) => {
  isEditing.value = true
  formData.value = JSON.parse(JSON.stringify(target))
  startTimeInput.value = formData.value.timeRange.start
  endTimeInput.value = formData.value.timeRange.end
  modalOpen.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formData.value.title.trim()) {
    toastError('请输入任务名称')
    return
  }

  submitting.value = true
  try {
    if (isEditing.value && currentTarget.value) {
      await checkinApi.update({
        id: currentTarget.value.id,
        title: formData.value.title,
        description: formData.value.description,
        time_range: formData.value.timeRange
      })
      success('更新成功')
    } else {
      await checkinApi.create({
        title: formData.value.title,
        description: formData.value.description,
        time_range: formData.value.timeRange
      })
      success('创建成功')
    }
    
    modalOpen.value = false
    await fetchTargets()
  } catch (e) {
    console.error('submit form failed', e)
    toastError('操作失败')
  } finally {
    submitting.value = false
  }
}

// 打卡逻辑
const upsertCheckIn = async (targetTs: number) => {
  if (!getCurrentTarget.value) {
    toastError('请先选择或创建打卡任务')
    return
  }

  const targetKey = dayKey(targetTs)
  const currentCounts = { ...getCurrentTarget.value.check_in_counts }
  
  // 检查是否已经打过卡
  if ((currentCounts[targetKey] || 0) > 0) {
    success('此日已打卡过')
    return
  }

  submitting.value = true
  try {
    const resp = await checkinApi.checkIn({
      id: getCurrentTarget.value.id,
      target_ts: targetTs
    })
    
    if (!resp.success) {
      toastError(resp.message)
      return
    }

    // 更新本地数据
    currentCounts[targetKey] = 1
    const updatedTarget = {
      ...getCurrentTarget.value,
      check_in_counts: currentCounts,
      total_days: getCurrentTarget.value.total_days + 1
    }
    
    const idx = targets.value.findIndex(t => t.id === getCurrentTarget.value?.id)
    if (idx >= 0) {
      targets.value[idx] = updatedTarget
    }

    success(targetKey === dayKey(Date.now()) ? '今日打卡成功！' : `补卡成功：${targetKey}`)
  } catch (e) {
    console.error('check in failed', e)
    toastError('打卡失败')
  } finally {
    submitting.value = false
  }
}

const handleCheckInToday = async () => {
  await upsertCheckIn(Date.now())
}

const handleCheckInByDate = async (ts: number) => {
  const todayStart = normalizeDayStart(Date.now())
  if (normalizeDayStart(ts) > todayStart) {
    toastError('不能打卡未来时间')
    return
  }
  await upsertCheckIn(ts)
}

watch(currentTargetId, () => {
  // 当切换任务时，重新计算
})

onMounted(() => {
  fetchTargets()
})
</script>

<style scoped>
.checkin-target-page {
  display: grid;
  gap: 16px;
}

/* Hero Section */
.hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 24px;
  border-radius: 16px;
}

.hero h2 {
  font-size: 28px;
  margin: 0 0 8px 0;
}

.hero p {
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
}

.hero-actions {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-outline {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.4);
  color: white;
}

.btn-outline:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
}

.btn-primary {
  background: white;
  color: #667eea;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.btn-primary.checked {
  background: #4ade80;
  color: #065f46;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Target List Section */
.target-list-section {
  padding: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h3 {
  font-size: 16px;
  margin: 0;
}

.section-header span {
  font-size: 13px;
  color: var(--text-secondary);
}

.target-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.target-chip {
  background: var(--bg-primary);
  border: 2px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  min-width: 200px;
  flex: 1;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.target-chip:hover {
  border-color: var(--accent-color, #667eea);
  transform: translateY(-2px);
}

.target-chip.active {
  border-color: #667eea;
  background: #f5f3ff;
}

.chip-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name {
  font-weight: 600;
  font-size: 15px;
}

.desc {
  font-size: 12px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chip-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
}

.total-days {
  color: var(--accent-color, #667eea);
  font-weight: 500;
}

.time-range {
  color: var(--text-secondary);
}

/* Heatmap Section */
.heatmap-section {
  padding: 16px;
}

.legend {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-secondary);
}

.legend i {
  width: 12px;
  height: 12px;
  border-radius: 3px;
  display: inline-block;
}

.lv1 { background: #d1fae5; }
.lv2 { background: #6ee7b7; }
.lv3 { background: #34d399; }
.lv4 { background: #10b981; }

.heatmap-wrap {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 8px;
  margin-top: 8px;
}

.weekday-labels {
  display: grid;
  grid-template-rows: repeat(7, 14px);
  gap: 4px;
  margin-top: 3px;
  font-size: 11px;
  color: var(--text-secondary);
}

.weekday-labels span {
  height: 14px;
  line-height: 14px;
}

.weekday-labels span:nth-child(1) { grid-row: 2; }
.weekday-labels span:nth-child(2) { grid-row: 4; }
.weekday-labels span:nth-child(3) { grid-row: 6; }

.weeks-container {
  display: flex;
  gap: 4px;
}

.week {
  display: grid;
  grid-template-rows: repeat(7, 14px);
  gap: 4px;
}

.day {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  border: 1px solid rgba(27, 31, 35, 0.08);
  background: #f0f2f5;
  padding: 0;
  cursor: pointer;
  transition: all 0.2s;
}

.day:hover:not(:disabled) {
  transform: scale(1.3);
  z-index: 1;
}

.day.today {
  outline: 2px solid #667eea;
  outline-offset: 1px;
}

.day.disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.stats-summary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  margin: 20px 0;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.stat-item {
  text-align: center;
}

.stat-item strong {
  font-size: 24px;
  color: var(--accent-color, #667eea);
  display: block;
}

.stat-item span {
  font-size: 12px;
  color: var(--text-secondary);
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-color);
}

.hint {
  font-size: 13px;
  color: var(--text-secondary);
  text-align: center;
  margin-top: 16px;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.5;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.edit-modal {
  background: white;
  border-radius: 16px;
  padding: 24px;
  width: 90%;
  max-width: 480px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.edit-modal h3 {
  margin: 0 0 20px 0;
  font-size: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.input,
.textarea {
  width: 100%;
  padding: 12px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: border-color 0.2s;
}

.input:focus,
.textarea:focus {
  outline: none;
  border-color: var(--accent-color, #667eea);
}

.textarea {
  resize: vertical;
  min-height: 80px;
}

.time-range-inputs {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
  padding: 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.time-input {
  flex: 1;
}

.time-input label {
  display: block;
  font-size: 12px;
  margin-bottom: 4px;
  color: var(--text-secondary);
}

.time-input input {
  width: 100%;
  padding: 8px;
  border: 2px solid var(--border-color);
  border-radius: 6px;
  font-size: 14px;
  text-align: center;
}

.separator {
  font-size: 18px;
  color: var(--text-secondary);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .hero {
    flex-direction: column;
    align-items: stretch;
    padding: 20px;
  }

  .hero-actions {
    justify-content: space-between;
  }

  .target-chip {
    min-width: 100%;
  }

  .stats-summary {
    flex-direction: column;
    gap: 12px;
  }

  .stat-divider {
    width: 100%;
    height: 1px;
  }
}
</style>
