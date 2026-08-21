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
    <!-- 大气通透的顶部导航栏 -->
    <header class="border-b border-border/50 backdrop-blur-md supports-[backdrop-filter]:bg-background/80 sticky top-0 z-50 transition-colors">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 h-18 flex items-center justify-between">
        <router-link to="/" class="flex items-center space-x-3 group">
          <span class="w-2.5 h-2.5 rounded-full bg-foreground group-hover:scale-125 transition-transform duration-200 inline-block"></span>
          <div class="flex flex-col">
            <span class="font-semibold text-base tracking-tight text-foreground group-hover:opacity-85 transition-opacity">Minimal Nav</span>
            <span class="text-[11px] text-muted-foreground tracking-wider uppercase font-mono hidden sm:inline-block">Internal Portal</span>
          </div>
        </router-link>

        <div class="flex items-center space-x-4 sm:space-x-6 text-sm">
          <router-link 
            to="/" 
            class="px-3.5 py-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-secondary/60 transition-all font-medium"
            active-class="text-foreground bg-secondary/80 font-semibold"
          >
            首页
          </router-link>
          <router-link 
            to="/admin" 
            class="px-3.5 py-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-secondary/60 transition-all font-medium"
            active-class="text-foreground bg-secondary/80 font-semibold"
          >
            管理后台
          </router-link>

          <div class="h-4 w-px bg-border/80"></div>

          <button 
            @click="toggleDark"
            class="p-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-secondary transition-all cursor-pointer border border-border/60 hover:border-border shadow-subtle"
            title="切换浅色/深色模式"
          >
            <svg v-if="isDark" class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- 宽阔的主体区域 -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-6 sm:px-8 lg:px-12 py-10 sm:py-14">
      <router-view />
    </main>

    <!-- 底部 -->
    <footer class="border-t border-border/50 py-8 text-center text-xs text-muted-foreground">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 flex flex-col sm:flex-row justify-between items-center gap-3">
        <p>© Minimal Nav. Designed for high efficiency & restrained elegance.</p>
        <p class="font-mono text-xs text-muted-foreground/70">Go 1.21 + Vue 3 + Tailwind CSS</p>
      </div>
    </footer>
  </div>
</template>
