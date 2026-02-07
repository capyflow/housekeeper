<template>
  <div class="layout">
    <!-- 移动端遮罩层 -->
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false"></div>

    <aside class="sidebar" :class="{ 'sidebar-open': sidebarOpen }">
      <div class="sidebar-header">
        <h2 class="logo">Housekeeper</h2>
        <button class="sidebar-close" @click="sidebarOpen = false">
          <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/" class="nav-item" exact-active-class="active" @click="closeSidebarOnMobile">
          <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
          </svg>
          <span>仪表盘</span>
        </router-link>
        <button class="nav-item nav-group" :class="{ active: isNotesSection }" type="button" @click="toggleNotesMenu">
          <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
          <span>笔记</span>
          <svg class="icon caret" :class="{ open: notesMenuOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
          </svg>
        </button>
        <div v-show="notesMenuOpen" class="nav-submenu">
          <router-link to="/shareboard" class="nav-subitem" exact-active-class="active" @click="closeSidebarOnMobile">
            共享看板
          </router-link>
          <router-link to="/notes" class="nav-subitem" exact-active-class="active" @click="closeSidebarOnMobile">
            笔记
          </router-link>
        </div>
      </nav>
    </aside>

    <div class="main-content">
      <header class="header">
        <div class="header-left">
          <button class="menu-toggle" @click="sidebarOpen = !sidebarOpen">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </button>
          <h1 class="page-title">{{ currentPageTitle }}</h1>
        </div>
        <div class="header-right">
          <button class="theme-toggle" @click="themeStore.toggleTheme()" :title="themeStore.theme === 'light' ? '切换到深色模式' : '切换到浅色模式'">
            <svg v-if="themeStore.theme === 'light'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
            </svg>
            <svg v-else class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
            </svg>
          </button>
          <div class="user-info">
            <span class="username">{{ userStore.user?.username || 'Admin' }}</span>
            <button class="btn-logout" @click="handleLogout">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
              </svg>
              退出
            </button>
          </div>
        </div>
      </header>

      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const themeStore = useThemeStore()

const sidebarOpen = ref(false)
const notesMenuOpen = ref(false)

const isNotesSection = computed(() => {
  const name = route.name as string
  return name === 'ShareBoard' || name === 'Notes'
})

const currentPageTitle = computed(() => {
  const titleMap: Record<string, string> = {
    Dashboard: '仪表盘',
    ShareBoard: '共享看板',
    Notes: '笔记'
  }
  return titleMap[route.name as string] || '首页'
})

const closeSidebarOnMobile = () => {
  // 在移动端点击导航项后关闭侧边栏
  if (window.innerWidth <= 768) {
    sidebarOpen.value = false
  }
}

const toggleNotesMenu = () => {
  notesMenuOpen.value = !notesMenuOpen.value
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

watch(
  () => isNotesSection.value,
  (isActive) => {
    if (isActive) {
      notesMenuOpen.value = true
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: 280px;
  background: var(--sidebar-bg);
  color: white;
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 16px rgba(0, 0, 0, 0.1);
  position: relative;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 1px;
  background: linear-gradient(to bottom, transparent, rgba(255,255,255,0.1), transparent);
}

.sidebar-header {
  padding: 32px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: flex-start;
  min-height: 88px;
}

.logo {
  font-size: 24px;
  font-weight: 700;
  color: white;
  letter-spacing: -0.5px;
  line-height: 1;
  flex: 1;
}

.sidebar-nav {
  flex: 1;
  padding: 24px 12px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  margin-bottom: 8px;
  color: var(--sidebar-text);
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  font-weight: 500;
  position: relative;
}

.nav-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
  transform: translateX(4px);
}

.nav-item.active {
  background-color: rgba(255, 255, 255, 0.15);
  color: var(--sidebar-active);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 24px;
  background-color: white;
  border-radius: 0 4px 4px 0;
}

.nav-item .icon {
  width: 22px;
  height: 22px;
}

.nav-group {
  background: transparent;
  border: none;
  width: 100%;
  text-align: left;
}

.nav-group .caret {
  margin-left: auto;
  transition: transform 0.2s ease;
}

.nav-group .caret.open {
  transform: rotate(180deg);
}

.nav-submenu {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 4px 0 12px 12px;
}

.nav-subitem {
  color: var(--sidebar-text);
  padding: 10px 16px;
  border-radius: 10px;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.nav-subitem:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
  transform: translateX(4px);
}

.nav-subitem.active {
  background-color: rgba(255, 255, 255, 0.15);
  color: var(--sidebar-active);
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--bg-secondary);
}

.header {
  height: 72px;
  background-color: var(--bg-primary);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 36px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  position: relative;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0;
  flex: 1;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.5px;
  line-height: 72px;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.theme-toggle {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
}

.theme-toggle:hover {
  background-color: var(--primary-light);
  color: var(--primary-color);
  border-color: var(--primary-color);
  transform: rotate(15deg) scale(1.05);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
}

.theme-toggle .icon {
  width: 22px;
  height: 22px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 14px;
  padding-left: 16px;
  border-left: 1px solid var(--border-color);
}

.username {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-secondary);
}

.btn-logout {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background-color: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-logout:hover {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger-color);
  border-color: var(--danger-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.2);
}

.btn-logout .icon {
  width: 18px;
  height: 18px;
}

.content {
  flex: 1;
  padding: 32px 36px;
  overflow-y: auto;
  background-color: var(--bg-secondary);
}

/* 移动端样式 */
.sidebar-overlay {
  display: none;
}

.menu-toggle {
  display: none;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background-color: transparent;
  color: var(--text-primary);
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.menu-toggle:hover {
  background-color: var(--bg-secondary);
  border-color: var(--border-color);
}

.menu-toggle .icon {
  width: 24px;
  height: 24px;
}

.sidebar-close {
  display: none;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background-color: transparent;
  color: white;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.sidebar-close:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.sidebar-close .icon {
  width: 20px;
  height: 20px;
}

@media (max-width: 768px) {
  .layout {
    flex-direction: column;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 280px;
    transform: translateX(-100%);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    z-index: 1000;
  }

  .sidebar-open {
    transform: translateX(0);
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 999;
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .sidebar-header {
    padding: 20px 20px;
    min-height: 60px;
    justify-content: space-between;
  }

  .logo {
    font-size: 22px;
    flex: 0;
  }

  .sidebar-close {
    display: flex;
    flex-shrink: 0;
  }

  .menu-toggle {
    display: flex;
    margin-right: 12px;
  }

  .header {
    padding: 0 16px;
    height: 60px;
  }

  .header-left {
    gap: 0;
    display: flex;
    align-items: center;
  }

  .page-title {
    font-size: 18px;
    line-height: 60px;
    margin: 0;
  }

  .header-right {
    gap: 8px;
  }

  .theme-toggle {
    width: 36px;
    height: 36px;
  }

  .username {
    display: none;
  }

  .btn-logout {
    padding: 8px 12px;
    font-size: 13px;
  }

  .btn-logout span {
    display: none;
  }

  .btn-logout .icon {
    margin: 0;
  }

  .content {
    padding: 20px 16px;
  }

  .main-content {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 260px;
  }

  .sidebar-header {
    padding: 16px;
    min-height: 56px;
  }

  .logo {
    font-size: 20px;
  }

  .menu-toggle {
    width: 36px;
    height: 36px;
    margin-right: 8px;
  }

  .menu-toggle .icon {
    width: 22px;
    height: 22px;
  }

  .header {
    height: 56px;
  }

  .page-title {
    font-size: 16px;
    line-height: 56px;
  }

  .theme-toggle {
    width: 32px;
    height: 32px;
  }

  .theme-toggle .icon {
    width: 18px;
    height: 18px;
  }

  .btn-logout {
    padding: 6px 10px;
  }

  .btn-logout .icon {
    width: 16px;
    height: 16px;
  }

  .sidebar-close {
    width: 28px;
    height: 28px;
  }

  .sidebar-close .icon {
    width: 18px;
    height: 18px;
  }

  .content {
    padding: 16px 12px;
  }
}
</style>
