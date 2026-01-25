<template>
  <router-view />
  <ToastContainer />
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'
import ToastContainer from '@/components/ToastContainer.vue'

// Initialize theme store
const themeStore = useThemeStore()
</script>

<style>
/* CSS Reset & Base Styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Light Theme Variables */
:root[data-theme="light"] {
  /* Primary Colors */
  --primary-color: #3b82f6;
  --primary-light: #dbeafe;
  --primary-dark: #1e40af;

  /* Secondary Colors */
  --secondary-color: #10b981;
  --secondary-light: #d1fae5;

  /* Accent Colors */
  --accent-orange: #f59e0b;
  --accent-orange-light: #fef3c7;
  --accent-purple: #8b5cf6;
  --accent-purple-light: #ede9fe;

  /* Background Colors */
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;

  /* Text Colors */
  --text-primary: #1f2937;
  --text-secondary: #6b7280;
  --text-tertiary: #9ca3af;

  /* Border Colors */
  --border-color: #e5e7eb;
  --border-light: #f3f4f6;

  /* Sidebar */
  --sidebar-bg: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  --sidebar-text: rgba(255, 255, 255, 0.8);
  --sidebar-active: #ffffff;

  /* Status Colors */
  --success-color: #10b981;
  --warning-color: #f59e0b;
  --danger-color: #ef4444;
  --info-color: #3b82f6;

  /* Shadows */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

/* Dark Theme Variables */
:root[data-theme="dark"] {
  /* Primary Colors */
  --primary-color: #60a5fa;
  --primary-light: #1e3a5f;
  --primary-dark: #93c5fd;

  /* Secondary Colors */
  --secondary-color: #34d399;
  --secondary-light: #1e4d3c;

  /* Accent Colors */
  --accent-orange: #fbbf24;
  --accent-orange-light: #4d3a1a;
  --accent-purple: #a78bfa;
  --accent-purple-light: #3730a3;

  /* Background Colors */
  --bg-primary: #1f2937;
  --bg-secondary: #111827;
  --bg-tertiary: #0f172a;

  /* Text Colors */
  --text-primary: #f9fafb;
  --text-secondary: #d1d5db;
  --text-tertiary: #9ca3af;

  /* Border Colors */
  --border-color: #374151;
  --border-light: #4b5563;

  /* Sidebar */
  --sidebar-bg: linear-gradient(135deg, #4c51bf 0%, #5b21b6 100%);
  --sidebar-text: rgba(255, 255, 255, 0.7);
  --sidebar-active: #ffffff;

  /* Status Colors */
  --success-color: #34d399;
  --warning-color: #fbbf24;
  --danger-color: #f87171;
  --info-color: #60a5fa;

  /* Shadows */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.6);
}

body {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  transition: background-color 0.3s ease, color 0.3s ease;
}

/* Common Styles */
.input {
  padding: 10px 14px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  transition: all 0.2s;
  outline: none;
}

.input:hover {
  border-color: var(--primary-color);
}

.input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.input::placeholder {
  color: var(--text-tertiary);
}

button {
  cursor: pointer;
  border: none;
  outline: none;
  font-family: inherit;
  transition: all 0.2s;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn {
  padding: 10px 18px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.btn-outline {
  background-color: transparent;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-outline:hover:not(:disabled) {
  background-color: var(--bg-secondary);
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.btn-danger {
  background-color: var(--danger-color);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background-color: #dc2626;
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

/* Card Styles */
.card {
  background-color: var(--bg-primary);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal {
  background-color: var(--bg-primary);
  border-radius: 12px;
  padding: 28px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-xl);
  animation: modalSlideIn 0.3s ease-out;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Scrollbar Styles */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--bg-secondary);
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-tertiary);
}

/* Link Styles */
a {
  text-decoration: none;
  color: inherit;
}
</style>
