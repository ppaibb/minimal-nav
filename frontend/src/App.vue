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
    <!-- 🏛️ 方案 A: 瑞士出版物报刊标头 (Swiss Editorial Header - 冷峻直线、纯文字排版、零油腻) -->
    <header class="border-b border-border/80 sticky top-0 z-50 bg-background/90 backdrop-blur-md transition-colors">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 h-15 sm:h-16 flex items-center justify-between">
        <!-- 左侧: 冷峻几何点缀 + 报刊大写标头 + 代码等宽注释 -->
        <router-link to="/" class="flex items-center space-x-2 group select-none">
          <span class="w-2 h-2 bg-foreground shrink-0 group-hover:rotate-45 transition-transform duration-200 inline-block"></span>
          <div class="flex items-baseline space-x-2">
            <span class="font-mono font-semibold text-sm sm:text-base tracking-wider uppercase text-foreground group-hover:opacity-80 transition-opacity">
              MINIMAL NAV
            </span>
            <span class="text-xs font-mono text-muted-foreground hidden sm:inline-block">
              // 内部索引目录
            </span>
          </div>
        </router-link>

        <!-- 中间: 极简等宽元数据 (大屏展示，无色块) -->
        <div class="hidden md:flex items-center space-x-2 text-xs font-mono text-muted-foreground select-none">
          <span>[ 状态: 正常运行 ]</span>
        </div>

        <!-- 右侧: 纯文字导航 (无任何背景色块、纯粹字阶与下划线) -->
        <div class="flex items-center space-x-5 sm:space-x-7 text-sm font-medium">
          <router-link 
            to="/" 
            class="text-muted-foreground hover:text-foreground transition-colors cursor-pointer py-1"
            active-class="text-foreground font-semibold underline underline-offset-8 decoration-foreground/60"
          >
            首页
          </router-link>
          <router-link 
            to="/admin" 
            class="text-muted-foreground hover:text-foreground transition-colors cursor-pointer py-1"
            active-class="text-foreground font-semibold underline underline-offset-8 decoration-foreground/60"
          >
            管理后台
          </router-link>

          <!-- 极细垂直分割线 -->
          <div class="h-3.5 w-px bg-border/80"></div>

          <!-- 细线条模式切换按钮 (零背景色块) -->
          <button 
            @click="toggleDark"
            class="text-muted-foreground hover:text-foreground transition-colors cursor-pointer p-1 focus:outline-none"
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

    <!-- 底部 (冷峻直线分割) -->
    <footer class="border-t border-border/80 py-8 text-xs text-muted-foreground">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 flex flex-col sm:flex-row justify-between items-center gap-3">
        <p>© Minimal Nav. Designed for high efficiency & restrained elegance.</p>
        <p class="font-mono text-xs text-muted-foreground/70">Go 1.21 + Vue 3 + Tailwind CSS</p>
      </div>
    </footer>
  </div>
</template>
