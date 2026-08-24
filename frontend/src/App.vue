<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import CommandPalette from './components/CommandPalette.vue'
import { useSiteConfig } from './utils/useSiteConfig'

const { siteConfig, loadSiteConfig } = useSiteConfig()

interface LinkItem {
  id: number
  title: string
  url: string
  category: string
  icon?: string
}

type ThemeMode = 'light' | 'dark' | 'system'

const themeMode = ref<ThemeMode>('system')
const isDropdownOpen = ref(false)
const dropdownRef = ref<HTMLDivElement | null>(null)
const globalLinks = ref<LinkItem[]>([])

// 获取全局链接列表供 ⌘K 搜索面板使用
const fetchGlobalLinks = async () => {
  try {
    const res = await fetch('/api/links')
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0 && Array.isArray(data.data)) {
        globalLinks.value = data.data
      }
    }
  } catch (err) {
    console.error('Failed to fetch global links:', err)
  }
}

// 应用主题 (支持平滑视觉过渡)
const applyTheme = (mode: ThemeMode, event?: MouseEvent) => {
  themeMode.value = mode
  localStorage.setItem('theme_mode', mode)

  const updateDOM = () => {
    if (mode === 'system') {
      const isSystemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      if (isSystemDark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    } else if (mode === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // 如果浏览器支持 View Transitions API 且用户没有设置减少动画
  if (
    event &&
    'startViewTransition' in document &&
    !window.matchMedia('(prefers-reduced-motion: reduce)').matches
  ) {
    const x = event.clientX
    const y = event.clientY
    const endRadius = Math.hypot(
      Math.max(x, window.innerWidth - x),
      Math.max(y, window.innerHeight - y)
    )

    // @ts-ignore
    const transition = document.startViewTransition(() => {
      updateDOM()
    })

    transition.ready.then(() => {
      const clipPath = [
        `circle(0px at ${x}px ${y}px)`,
        `circle(${endRadius}px at ${x}px ${y}px)`
      ]
      document.documentElement.animate(
        {
          clipPath: mode === 'dark' ? clipPath : [...clipPath].reverse()
        },
        {
          duration: 350,
          easing: 'ease-in-out',
          pseudoElement: mode === 'dark' ? '::view-transition-new(root)' : '::view-transition-old(root)'
        }
      )
    })
  } else {
    updateDOM()
  }

  isDropdownOpen.value = false
}

// 监听系统主题变化
const handleSystemThemeChange = (e: MediaQueryListEvent) => {
  if (themeMode.value === 'system') {
    if (e.matches) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }
}

// 点击外部关闭下拉菜单
const handleClickOutside = (e: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    isDropdownOpen.value = false
  }
}

onMounted(() => {
  loadSiteConfig()
  const saved = (localStorage.getItem('theme_mode') as ThemeMode) || 'system'
  applyTheme(saved)
  fetchGlobalLinks()

  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', handleSystemThemeChange)
  window.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.removeEventListener('change', handleSystemThemeChange)
  window.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="min-h-screen bg-background text-foreground flex flex-col font-sans transition-colors duration-200 selection:bg-zinc-200 dark:selection:bg-zinc-800">
    <!-- 🏛️ 极致极简顶栏 (左侧 Logo + 自定义站点名称，右侧 太阳皮肤切换 + 齿轮设置) -->
    <header class="border-b border-border/80 sticky top-0 z-40 bg-background/95 backdrop-blur-md supports-[backdrop-filter]:bg-background/60 transition-colors">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-14 sm:h-16 flex items-center justify-between">
        <!-- 左侧: 仅保留 shadcn 官方矢量 Logo 图标 + 站点名称 -->
        <router-link to="/" class="flex items-center space-x-2.5 group select-none">
          <!-- shadcn 官方矢量 Logo 图标 -->
          <div class="w-6 h-6 rounded-md bg-foreground text-background flex items-center justify-center p-1 shadow-xs group-hover:opacity-90 transition-opacity shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" class="h-full w-full fill-current">
              <rect width="256" height="256" fill="none"></rect>
              <line x1="208" y1="128" x2="128" y2="208" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="28"></line>
              <line x1="192" y1="40" x2="40" y2="192" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="28"></line>
            </svg>
          </div>

          <span class="font-semibold text-sm sm:text-base tracking-tight text-foreground">
            {{ siteConfig.site_name }}
          </span>
        </router-link>

        <!-- 右侧: 仅保留 皮肤切换 (默认 Light 太阳图标) 与 齿轮设置 (与 Light 按钮样式完全一致，不做加深) -->
        <div class="flex items-center space-x-1.5 sm:space-x-2 text-sm">
          <!-- 🌟 皮肤切换下拉菜单 -->
          <div ref="dropdownRef" class="relative">
            <button
              @click.stop="isDropdownOpen = !isDropdownOpen"
              class="w-8 h-8 sm:w-9 sm:h-9 rounded-md flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-accent transition-colors cursor-pointer focus:outline-none"
              title="切换显示模式"
            >
              <!-- 默认常态展示 Light 太阳图标 -->
              <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
                <circle cx="12" cy="12" r="4" />
                <path stroke-linecap="round" d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41" />
              </svg>
            </button>

            <!-- 下拉浮层卡片 (纯正 shadcn/ui DropdownMenu 规范) -->
            <transition
              enter-active-class="transition duration-150 ease-out"
              enter-from-class="transform scale-95 opacity-0 -translate-y-1"
              enter-to-class="transform scale-100 opacity-100 translate-y-0"
              leave-active-class="transition duration-100 ease-in"
              leave-from-class="transform scale-100 opacity-100 translate-y-0"
              leave-to-class="transform scale-95 opacity-0 -translate-y-1"
            >
              <div
                v-if="isDropdownOpen"
                class="absolute right-0 mt-2 w-36 rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-lg z-50 focus:outline-none"
              >
                <!-- Light -->
                <button
                  @click="applyTheme('light', $event)"
                  class="w-full flex items-center justify-between px-2.5 py-1.5 text-xs sm:text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors cursor-pointer text-left"
                >
                  <div class="flex items-center space-x-2">
                    <svg class="w-3.5 h-3.5 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
                      <circle cx="12" cy="12" r="4" />
                      <path stroke-linecap="round" d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41" />
                    </svg>
                    <span>Light</span>
                  </div>
                  <svg v-if="themeMode === 'light'" class="w-3.5 h-3.5 text-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </button>

                <!-- Dark -->
                <button
                  @click="applyTheme('dark', $event)"
                  class="w-full flex items-center justify-between px-2.5 py-1.5 text-xs sm:text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors cursor-pointer text-left"
                >
                  <div class="flex items-center space-x-2">
                    <svg class="w-3.5 h-3.5 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
                    </svg>
                    <span>Dark</span>
                  </div>
                  <svg v-if="themeMode === 'dark'" class="w-3.5 h-3.5 text-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </button>

                <!-- System -->
                <button
                  @click="applyTheme('system', $event)"
                  class="w-full flex items-center justify-between px-2.5 py-1.5 text-xs sm:text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors cursor-pointer text-left"
                >
                  <div class="flex items-center space-x-2">
                    <svg class="w-3.5 h-3.5 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
                      <rect width="20" height="14" x="2" y="3" rx="2" />
                      <line x1="8" x2="16" y1="21" y2="21" />
                      <line x1="12" x2="12" y1="17" y2="21" />
                    </svg>
                    <span>System</span>
                  </div>
                  <svg v-if="themeMode === 'system'" class="w-3.5 h-3.5 text-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                </button>
              </div>
            </transition>
          </div>

          <!-- ⚙️ 管理后台 (齿轮设置按钮 - 与 Light 按钮完全一致，无多余加深背景) -->
          <router-link 
            to="/admin" 
            class="w-8 h-8 sm:w-9 sm:h-9 rounded-md flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            title="管理后台配置"
          >
            <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 010 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.6 6.6 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 010-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28z" />
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </router-link>
        </div>
      </div>
    </header>

    <!-- 宽阔的主体区域 -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-6 sm:px-8 lg:px-12 py-10 sm:py-14">
      <router-view />
    </main>

    <!-- 底部 (贯穿直线) -->
    <footer class="border-t border-border/80 py-8 text-xs text-muted-foreground">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 flex flex-col sm:flex-row justify-between items-center gap-3">
        <p>© Minimal Nav. Designed for high efficiency & restrained elegance.</p>
        <p class="font-mono text-xs text-muted-foreground/70">Go 1.21 + Vue 3 + Tailwind CSS</p>
      </div>
    </footer>

    <!-- 🌟 全局 Command ⌘K 搜索面板 -->
    <CommandPalette ref="commandPaletteRef" :links="globalLinks" />
  </div>
</template>
