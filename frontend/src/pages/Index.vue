<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted } from 'vue'

interface LinkItem {
  id: number
  title: string
  url: string
  category: string
  icon?: string
  latency_ms?: number
  healthy?: boolean
}

interface AnnouncementItem {
  id: number
  content: string
  is_active: boolean
  created_at: string
}

const announcements = ref<AnnouncementItem[]>([])
const links = ref<LinkItem[]>([])
const loading = ref(true)
const selectedCategory = ref<string>('all')
const searchQuery = ref<string>('')
const searchInputRef = ref<HTMLInputElement | null>(null)
const healthStatusMap = ref<Record<number, { healthy: boolean; latency_ms: number }>>({})

// 获取生效公告
const fetchAnnouncements = async () => {
  try {
    const res = await fetch('/api/announcements/active')
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0 && Array.isArray(data.data)) {
        announcements.value = data.data
      }
    }
  } catch (err) {
    console.error('Failed to fetch announcements:', err)
  }
}

// 获取链接列表
const fetchLinks = async () => {
  try {
    const res = await fetch('/api/links')
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0 && Array.isArray(data.data)) {
        links.value = data.data
        // 自动探测前台健康状态 (非阻塞)
        probeLinksHealth(data.data)
      }
    }
  } catch (err) {
    console.error('Failed to fetch links:', err)
  } finally {
    loading.value = false
  }
}

// 异步轻量探测应用健康度
const probeLinksHealth = async (items: LinkItem[]) => {
  // 逐个轻量探测
  for (const item of items.slice(0, 10)) {
    try {
      const res = await fetch('/api/tools/ping', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: item.url }),
      })
      if (res.ok) {
        const json = await res.json()
        if (json.code === 0 && json.data) {
          healthStatusMap.value[item.id] = {
            healthy: json.data.healthy,
            latency_ms: json.data.latency_ms,
          }
        }
      }
    } catch {
      // 忽略单个报错
    }
  }
}

// 提取 URL 中的简要域名
const extractHostname = (url: string) => {
  try {
    const u = new URL(url.startsWith('http') ? url : `https://${url}`)
    return u.hostname.replace(/^www\./, '')
  } catch {
    return url
  }
}

// 提取首字母用于图标展示
const getInitial = (title: string) => {
  if (!title) return 'N'
  return title.trim().charAt(0).toUpperCase()
}

// 计算所有分类
const categories = computed(() => {
  const cats = new Set<string>()
  links.value.forEach(link => {
    if (link.category) cats.add(link.category)
  })
  return Array.from(cats)
})

// 根据搜索和分类过滤链接
const filteredLinks = computed(() => {
  return links.value.filter(link => {
    const matchCat = selectedCategory.value === 'all' || link.category === selectedCategory.value
    const matchSearch = !searchQuery.value.trim() || 
      link.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      link.url.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      link.category.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchCat && matchSearch
  })
})

// 图标加载失败优雅降级
const handleIconError = (e: Event) => {
  const target = e.target as HTMLElement
  target.style.display = 'none'
  if (target.nextElementSibling) {
    (target.nextElementSibling as HTMLElement).style.display = 'flex'
  }
}

// 快捷键 Ctrl+K / Cmd+K 快速聚焦搜索框
const handleKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchInputRef.value?.focus()
  }
}

onMounted(() => {
  fetchAnnouncements()
  fetchLinks()
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
  <div class="space-y-12 sm:space-y-16">
    <!-- 顶部大气问候与状态概览 (Hero Section) -->
    <section class="pb-2">
      <div>
        <h1 class="text-3xl sm:text-4xl font-bold tracking-tight text-foreground">
          企业效率与工具导航
        </h1>
        <p class="text-sm sm:text-base text-muted-foreground mt-2 max-w-2xl leading-relaxed">
          统一汇聚团队核心工具、部署控制台、设计协作及文档中心，即时检索快速直达。
        </p>
      </div>
    </section>

    <!-- 顶部公告区 (Announcements) -->
    <section v-if="announcements.length > 0" class="border-y border-border/50 py-3.5">
      <div class="space-y-2">
        <div 
          v-for="item in announcements" 
          :key="item.id" 
          class="text-xs sm:text-sm text-muted-foreground leading-relaxed flex items-center space-x-3"
        >
          <span class="px-2 py-0.5 text-[11px] font-semibold rounded bg-secondary text-foreground shrink-0">
            公告
          </span>
          
          <!-- 直接点击公告标题跳转至对应指南或页面 -->
          <router-link
            v-if="item.content.includes('AI') || item.content.includes('指南') || item.content.includes('Claude') || item.content.includes('Codex')"
            to="/docs/ai"
            class="tracking-tight text-foreground/90 hover:text-primary hover:underline transition-colors truncate cursor-pointer"
            title="点击查看 AI 编程助手接入指南"
          >
            {{ item.content }}
          </router-link>
          <span v-else class="tracking-tight text-foreground/90 truncate">
            {{ item.content }}
          </span>
        </div>
      </div>
    </section>

    <!-- 检索与分类过滤器控制台 -->
    <section class="flex flex-col lg:flex-row lg:items-center justify-between gap-5">
      <!-- 分类标签 -->
      <div class="flex flex-wrap items-center gap-2">
        <button
          @click="selectedCategory = 'all'"
          :class="[
            'px-4 py-2 rounded-lg transition-all cursor-pointer border text-sm font-medium',
            selectedCategory === 'all'
              ? 'bg-foreground text-background border-foreground shadow-sm'
              : 'border-border/60 text-muted-foreground hover:text-foreground hover:bg-secondary/80'
          ]"
        >
          全部应用 ({{ links.length }})
        </button>
        <button
          v-for="cat in categories"
          :key="cat"
          @click="selectedCategory = cat"
          :class="[
            'px-4 py-2 rounded-lg transition-all cursor-pointer border text-sm font-medium',
            selectedCategory === cat
              ? 'bg-foreground text-background border-foreground shadow-sm'
              : 'border-border/60 text-muted-foreground hover:text-foreground hover:bg-secondary/80'
          ]"
        >
          {{ cat }}
        </button>
      </div>

      <!-- 强化的大气搜索框 -->
      <div class="relative w-full lg:w-80">
        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-muted-foreground">
          <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input 
          ref="searchInputRef"
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索快捷方式或地址..." 
          class="w-full h-11 bg-background border border-border text-sm rounded-lg pl-10 pr-16 placeholder:text-muted-foreground/60 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground transition-all shadow-subtle"
        />
        <div class="absolute inset-y-0 right-0 pr-3 flex items-center space-x-1.5 pointer-events-none">
          <span v-if="!searchQuery" class="text-[11px] font-mono bg-secondary px-1.5 py-0.5 rounded border border-border/80 text-muted-foreground">
            ⌘K
          </span>
          <button 
            v-else 
            @click="searchQuery = ''" 
            class="pointer-events-auto text-xs text-muted-foreground hover:text-foreground p-1"
          >
            ✕
          </button>
        </div>
      </div>
    </section>

    <!-- 导航网格区 (Navigation Grid) -->
    <section>
      <!-- 加载中骨架 -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
        <div 
          v-for="n in 8" 
          :key="n" 
          class="h-28 rounded-xl border border-border/60 bg-muted/20 animate-pulse"
        ></div>
      </div>

      <!-- 链接大卡片网格 (支持 Favicon 高清渲染 + 优雅回退 + 健康探测微呼吸点) -->
      <div 
        v-else-if="filteredLinks.length > 0" 
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-5"
      >
        <a 
          v-for="link in filteredLinks" 
          :key="link.id" 
          :href="link.url" 
          target="_blank" 
          rel="noopener noreferrer"
          class="group relative flex items-start p-5 sm:p-5.5 rounded-xl border border-border/80 bg-card text-card-foreground hover:border-zinc-400 dark:hover:border-zinc-600 hover:-translate-y-1 hover:shadow-lg hover:shadow-zinc-200/50 dark:hover:shadow-zinc-950/60 transition-all duration-200 ease-out cursor-pointer"
        >
          <!-- 🌟 左侧 Favicon 真实图标槽 (带优雅降级) -->
          <div class="relative w-11 h-11 rounded-lg bg-secondary flex items-center justify-center border border-border/50 shrink-0 group-hover:scale-105 transition-transform duration-200 select-none overflow-hidden p-1.5">
            <img 
              v-if="link.icon" 
              :src="link.icon" 
              :alt="link.title"
              class="w-full h-full object-contain"
              @error="handleIconError"
            />
            <div 
              class="w-full h-full items-center justify-center font-bold text-base text-foreground"
              :class="link.icon ? 'hidden' : 'flex'"
            >
              {{ getInitial(link.title) }}
            </div>

            <!-- 📶 健康度探测微状态指示点 -->
            <span 
              v-if="healthStatusMap[link.id]"
              class="absolute top-1 right-1 w-2 h-2 rounded-full border border-background"
              :class="healthStatusMap[link.id].healthy ? 'bg-emerald-500' : 'bg-rose-500'"
              :title="healthStatusMap[link.id].healthy ? `服务正常 · 延迟: ${healthStatusMap[link.id].latency_ms}ms` : '服务暂时不可达'"
            ></span>
          </div>

          <!-- 中间标题与信息 -->
          <div class="ml-4 flex-1 min-w-0 pr-2">
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold tracking-tight text-foreground truncate group-hover:text-primary transition-colors">
                {{ link.title }}
              </h3>
            </div>
            
            <p class="text-xs text-muted-foreground mt-1 font-mono truncate">
              {{ extractHostname(link.url) }}
            </p>

            <!-- 底部：分类 (左) + 实时延迟健康指示 (右) -->
            <div class="mt-3.5 pt-2 border-t border-border/40 flex items-center justify-between text-xs">
              <span class="text-[11px] text-muted-foreground font-medium group-hover:text-foreground transition-colors truncate">
                {{ link.category || '默认' }}
              </span>

              <div 
                v-if="healthStatusMap[link.id]?.healthy"
                class="inline-flex items-center space-x-1.5 text-[10px] font-mono text-muted-foreground/70 shrink-0"
                :title="`连通正常 · 响应时间 ${healthStatusMap[link.id].latency_ms}ms`"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0"></span>
                <span>{{ healthStatusMap[link.id].latency_ms }}ms</span>
              </div>
            </div>
          </div>

          <!-- 右上角跳转微动效箭头 -->
          <svg 
            class="w-4 h-4 text-muted-foreground/30 group-hover:text-foreground group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-all duration-200 shrink-0 ml-1 mt-0.5" 
            xmlns="http://www.w3.org/2000/svg" 
            fill="none" 
            viewBox="0 0 24 24" 
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 17L17 7M17 7H7M17 7V17" />
          </svg>
        </a>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-20 border border-dashed border-border/80 rounded-2xl bg-secondary/20">
        <div class="w-12 h-12 rounded-full bg-secondary flex items-center justify-center mx-auto mb-3 text-muted-foreground">
          <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <p class="text-sm font-medium text-foreground">未找到相关导航链接</p>
        <p class="text-xs text-muted-foreground mt-1">您可以尝试更换关键词或清除筛选条件</p>
      </div>
    </section>
  </div>
</template>
