<template>
  <div class="record-page">
    <section class="hero card">
      <div class="hero-text">
        <h2>{{ activeRecord?.title || '打卡记录' }}</h2>
        <p>用热力图追踪习惯，坚持打卡让今天亮起来。</p>
      </div>
      <div class="hero-actions">
        <button class="btn btn-outline" @click="openCreate" :disabled="submitting">新建记录</button>
        <button class="btn btn-primary" @click="checkInToday" :disabled="submitting || !activeRecord">
          {{ checkedToday ? '今天已打卡' : '今天打卡' }}
        </button>
      </div>
    </section>

    <section class="record-grid card">
      <div class="record-head">
        <h3>我的记录</h3>
        <span>{{ records.length }} 条</span>
      </div>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="records.length === 0" class="empty">还没有打卡记录，先创建一个。</div>
      <div v-else class="record-list">
        <button
          v-for="item in records"
          :key="item.id"
          class="record-chip"
          :class="{ active: item.id === activeRecordId }"
          @click="activeRecordId = item.id"
        >
          <span class="name">{{ item.title || '未命名记录' }}</span>
          <span class="meta">{{ formatDate(item.create_ts) }}</span>
        </button>
      </div>
    </section>

    <section class="heatmap card" v-if="activeRecord">
      <div class="heatmap-head">
        <h3>近一年打卡热力图</h3>
        <div class="legend">
          <span>少</span>
          <i class="lv0"></i>
          <i class="lv1"></i>
          <i class="lv2"></i>
          <i class="lv3"></i>
          <span>多</span>
        </div>
      </div>

      <div class="heatmap-wrap">
        <div class="weekday-col">
          <span>一</span>
          <span>三</span>
          <span>五</span>
        </div>
        <div class="weeks">
          <div class="week" v-for="(week, wIdx) in heatmapWeeks" :key="wIdx">
            <button
              v-for="cell in week"
              :key="cell.dateKey"
              class="day"
              :class="[`lv${cell.level}`, { today: cell.isToday }]"
              :title="cell.tooltip"
              @click="checkInByDate(cell.date)"
            ></button>
          </div>
        </div>
      </div>

      <p class="hint">点击任意日期可补打卡；今天打卡后，今天格子会点亮。</p>
    </section>

    <div v-if="createModalOpen" class="modal-overlay" @click="createModalOpen = false">
      <div class="modal create-modal" @click.stop>
        <h3>新建打卡记录</h3>
        <div class="form-item">
          <label>记录名称</label>
          <input class="input" v-model="createTitle" maxlength="50" placeholder="例如：健身 / 阅读 / 学英语" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="createModalOpen = false">取消</button>
          <button class="btn btn-primary" :disabled="submitting || !createTitle.trim()" @click="createRecord">创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { recordApi } from '@/api/record'
import type { Record as CheckRecord } from '@/types'
import { useToast } from '@/composables/useToast'

const { success, error: toastError } = useToast()

const loading = ref(false)
const submitting = ref(false)
const records = ref<CheckRecord[]>([])
const activeRecordId = ref('')

const createModalOpen = ref(false)
const createTitle = ref('')

const DAY_MS = 24 * 60 * 60 * 1000

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

const activeRecord = computed(() => records.value.find((v) => v.id === activeRecordId.value) || null)

const countByDay = computed(() => {
  const map: { [key: string]: number } = {}
  const timestamps = activeRecord.value?.record_ts || []
  timestamps.forEach((ts) => {
    const key = dayKey(ts)
    map[key] = (map[key] || 0) + 1
  })
  return map
})

const checkedToday = computed(() => {
  const today = dayKey(Date.now())
  return (countByDay.value[today] || 0) > 0
})

const heatmapWeeks = computed(() => {
  const now = new Date()
  const todayStart = normalizeDayStart(now.getTime())
  const start = new Date(todayStart - 364 * DAY_MS)
  const startOffset = start.getDay()
  const alignedStart = start.getTime() - startOffset * DAY_MS

  const totalDays = 53 * 7
  const weeks: Array<Array<{ date: number; dateKey: string; level: 0 | 1 | 2 | 3; tooltip: string; isToday: boolean }>> = []

  for (let i = 0; i < totalDays; i += 1) {
    const ts = alignedStart + i * DAY_MS
    const key = dayKey(ts)
    const count = countByDay.value[key] || 0
    const level: 0 | 1 | 2 | 3 = count === 0 ? 0 : count === 1 ? 1 : count <= 3 ? 2 : 3
    const inRange = ts >= start.getTime() && ts <= todayStart
    const dayCell = {
      date: ts,
      dateKey: key,
      level: inRange ? level : 0,
      tooltip: `${key} ${inRange ? `打卡 ${count} 次` : '超出范围'}`,
      isToday: ts === todayStart
    }
    const weekIndex = Math.floor(i / 7)
    if (!weeks[weekIndex]) {
      weeks[weekIndex] = []
    }
    weeks[weekIndex].push(dayCell)
  }

  return weeks
})

const formatDate = (ts: number) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleDateString('zh-CN')
}

const fetchRecords = async () => {
  loading.value = true
  try {
    const res = await recordApi.list({ page: 1, page_size: 100 })
    records.value = (res.list || []).sort((a, b) => b.create_ts - a.create_ts)
    if (!activeRecordId.value && records.value.length > 0) {
      activeRecordId.value = records.value[0].id
    }
    if (activeRecordId.value && !records.value.find((v) => v.id === activeRecordId.value) && records.value.length > 0) {
      activeRecordId.value = records.value[0].id
    }
  } catch (e) {
    console.error('fetch record list failed', e)
    toastError('加载打卡记录失败')
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  createTitle.value = ''
  createModalOpen.value = true
}

const createRecord = async () => {
  if (!createTitle.value.trim()) return
  submitting.value = true
  try {
    const res = await recordApi.create({
      title: createTitle.value.trim(),
      record_ts: []
    })
    createModalOpen.value = false
    await fetchRecords()
    if (res?.info?.id) {
      activeRecordId.value = res.info.id
    }
    success('创建成功')
  } catch (e) {
    console.error('create record failed', e)
    toastError('创建失败')
  } finally {
    submitting.value = false
  }
}

const upsertRecordDate = async (targetTs: number) => {
  if (!activeRecord.value) {
    toastError('请先创建打卡记录')
    return
  }

  const targetKey = dayKey(targetTs)
  const nextSet = new Set((activeRecord.value.record_ts || []).map((v) => normalizeDayStart(v)))
  nextSet.add(normalizeDayStart(targetTs))

  const merged = Array.from(nextSet)
    .map((v) => normalizeDayStart(v))
    .sort((a, b) => a - b)

  if ((countByDay.value[targetKey] || 0) > 0) {
    return
  }

  submitting.value = true
  try {
    await recordApi.update({
      id: activeRecord.value.id,
      title: activeRecord.value.title,
      record_ts: merged
    })
    const idx = records.value.findIndex((v) => v.id === activeRecord.value?.id)
    if (idx >= 0) {
      records.value[idx] = { ...records.value[idx], record_ts: merged }
    }
    success(targetKey === dayKey(Date.now()) ? '今日已打卡' : `已补卡 ${targetKey}`)
  } catch (e) {
    console.error('update record failed', e)
    toastError('打卡失败')
  } finally {
    submitting.value = false
  }
}

const checkInToday = async () => {
  await upsertRecordDate(Date.now())
}

const checkInByDate = async (ts: number) => {
  const todayStart = normalizeDayStart(Date.now())
  if (normalizeDayStart(ts) > todayStart) return
  await upsertRecordDate(ts)
}

onMounted(() => {
  fetchRecords()
})
</script>

<style scoped>
.record-page {
  display: grid;
  gap: 16px;
}

.hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #0f172a 0%, #14532d 100%);
  color: #f8fafc;
}

.hero-text h2 {
  font-size: 26px;
  line-height: 1.2;
  margin-bottom: 8px;
}

.hero-text p {
  color: #d1fae5;
}

.hero-actions {
  display: flex;
  gap: 12px;
}

.record-grid {
  padding: 16px;
}

.record-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.record-head h3 {
  font-size: 16px;
}

.record-head span {
  font-size: 13px;
  color: var(--text-secondary);
}

.record-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.record-chip {
  border: 1px solid var(--border-color);
  background: var(--bg-primary);
  border-radius: 10px;
  padding: 8px 12px;
  display: grid;
  text-align: left;
  gap: 4px;
  min-width: 170px;
}

.record-chip.active {
  border-color: #16a34a;
  background: #f0fdf4;
}

.name {
  font-weight: 600;
}

.meta {
  font-size: 12px;
  color: var(--text-secondary);
}

.heatmap-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
  gap: 10px;
  flex-wrap: wrap;
}

.legend {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}

.legend i {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  display: inline-block;
}

.lv0 { background: #ebedf0; }
.lv1 { background: #9be9a8; }
.lv2 { background: #40c463; }
.lv3 { background: #30a14e; }

.heatmap-wrap {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 6px;
}

.weekday-col {
  display: grid;
  grid-template-rows: repeat(7, 12px);
  gap: 4px;
  margin-top: 2px;
  font-size: 11px;
  color: var(--text-secondary);
}

.weekday-col span {
  height: 12px;
  line-height: 12px;
}

.weekday-col span:nth-child(1) { grid-row: 2; }
.weekday-col span:nth-child(2) { grid-row: 4; }
.weekday-col span:nth-child(3) { grid-row: 6; }

.weeks {
  display: flex;
  gap: 4px;
}

.week {
  display: grid;
  grid-template-rows: repeat(7, 12px);
  gap: 4px;
}

.day {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  border: 1px solid rgba(27, 31, 35, 0.06);
  background: #ebedf0;
  padding: 0;
}

.day.today {
  outline: 2px solid #0ea5e9;
  outline-offset: 1px;
}

.day:hover {
  transform: scale(1.15);
  border-color: rgba(0, 0, 0, 0.2);
}

.hint {
  margin-top: 12px;
  font-size: 13px;
  color: var(--text-secondary);
}

.loading,
.empty {
  color: var(--text-secondary);
  font-size: 14px;
}

.create-modal h3 {
  margin-bottom: 16px;
}

.form-item label {
  display: block;
  font-size: 14px;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 16px;
}

@media (max-width: 768px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .hero-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
