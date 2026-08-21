<script setup lang="ts">
import { ref, onMounted } from 'vue'

const isDark = ref(false)

const toggleDark = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved === 'dark' || (!saved && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  } else {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  }
})
</script>

<template>
  <div class="min-h-screen bg-background text-foreground flex flex-col font-sans transition-colors duration-200 selection:bg-zinc-200 dark:selection:bg-zinc-800">
    <!-- 🏛️ 正统 shadcn-admin 标准顶栏 (贴顶贯穿线、Ghost 导航按钮、官方标准 Icon) -->
    <header class="border-b border-border/80 sticky top-0 z-50 bg-background/95 backdrop-blur-md supports-[backdrop-filter]:bg-background/60 transition-colors">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 h-14 sm:h-16 flex items-center justify-between">
        <!-- 左侧: shadcn 标志性微图标 + Minimal Nav + 标识 -->
        <router-link to="/" class="flex items-center space-x-2.5 group select-none">
          <!-- shadcn 官方矢量 Logo 图标 -->
          <div class="w-6 h-6 rounded-md bg-foreground text-background flex items-center justify-center p-1 shadow-xs group-hover:opacity-90 transition-opacity shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" class="h-full w-full fill-current">
              <rect width="256" height="256" fill="none"></rect>
              <line x1="208" y1="128" x2="128" y2="208" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="28"></line>
              <line x1="192" y1="40" x2="40" y2="192" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="28"></line>
            </svg>
          </div>

          <div class="flex items-baseline space-x-2">
            <span class="font-semibold text-sm sm:text-base tracking-tight text-foreground">
              Minimal Nav
            </span>
            <span class="text-xs text-muted-foreground font-mono hidden sm:inline-block">
              / Internal Portal
            </span>
          </div>
        </router-link>

        <!-- 右侧: shadcn 标准 Ghost 导航与工具栏 -->
        <div class="flex items-center space-x-1 sm:space-x-2 text-sm font-medium">
          <router-link 
            to="/" 
            class="px-3 py-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/80 transition-colors font-medium"
            active-class="text-foreground bg-muted font-semibold shadow-xs"
          >
            首页
          </router-link>
          <router-link 
            to="/admin" 
            class="px-3 py-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/80 transition-colors font-medium"
            active-class="text-foreground bg-muted font-semibold shadow-xs"
          >
            管理后台
          </router-link>

          <!-- 细分割竖线 -->
          <div class="h-4 w-px bg-border/80 mx-1"></div>

          <!-- shadcn 标准外框 Icon 按钮 -->
          <button 
            @click="toggleDark"
            class="w-8 h-8 rounded-md border border-border/80 hover:bg-muted hover:text-foreground flex items-center justify-center text-muted-foreground transition-colors cursor-pointer focus:outline-none"
            :title="isDark ? '切换浅色模式' : '切换深色模式'"
          >
            <svg v-if="isDark" class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
              <circle cx="12" cy="12" r="4" />
              <path stroke-linecap="round" d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41" />
            </svg>
            <svg v-else class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
            </svg>
          </button>
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
  </div>
</template>
