import { createApp } from 'vue'
import { createPinia } from 'pinia'
import i18n from '@/i18n'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import App from '@/App.vue'
import '@/style.css'

// ============================================================
// API 就绪检测 (Wails / 浏览器兼容)
// ============================================================
const waitForApi = (): Promise<void> => {
  return new Promise((resolve) => {
    if (window.pywebview && window.pywebview.api) {
      resolve()
      return
    }

    window.addEventListener('pywebviewready', () => {
      resolve()
    }, { once: true })

    // 超时兜底，防止挂起
    setTimeout(resolve, 3000)
  })
}

// ============================================================
// 启动 Vue 应用
// ============================================================
waitForApi().then(() => {
  // 设置前端日志重定向至后端
  const api = window.pywebview?.api
  if (api && typeof api.log_from_js === 'function') {
    const levels = ['log', 'info', 'warn', 'error', 'debug'] as const

    levels.forEach(level => {
      const original = console[level]
      ;(console as any)[level] = (...args: any[]) => {
        original.apply(console, args)
        const message = args.map(arg =>
          typeof arg === 'object' ? JSON.stringify(arg, null, 2) : String(arg)
        ).join(' ')
        const mappedLevel = level === 'log' ? 'info' : level
        api.log_from_js(mappedLevel, message).catch(() => {})
      }
    })
    console.log('🚀 前端日志已连接至 Wails 后端')
  }

  const app = createApp(App)

  // 导出进度回调
  ;(window as any).__onExportProgress = (payload: { stage: string; done: number; total: number }) => {
    const { exportProgress } = useVpkManagerStore()
    exportProgress.stage = payload.stage
    exportProgress.done = payload.done
    exportProgress.total = payload.total
  }

  app.use(createPinia())
  app.use(i18n)

  app.mount('#app')
})
