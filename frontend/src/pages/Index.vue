<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { useSiteConfig } from '../utils/useSiteConfig'

interface AnnouncementItem {
  id: number
  content: string
  detail_md?: string
  is_active: boolean
  created_at: string
}

interface LinkItem {
  id: number
  title: string
  url: string
  icon: string
  category: string
  sort_order: number
}

const announcements = ref<AnnouncementItem[]>([])
const links = ref<LinkItem[]>() || ref<LinkItem[]>([])
const loading = ref(true)
const searchQuery = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)

// 响应式健康状态存储 { [linkId: number]: { healthy: boolean, latency_ms: number } }
const healthStatusMap = ref<Record<number, { healthy: boolean; latency_ms: number }>>({})

// 分类管理
const selectedCategory = ref<string>('all')

// 获取公告列表
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

// 纯前端客户端网络连通性与延迟探测 (直接反映当前用户本地网络环境)
const pingClient = async (rawUrl: string, timeoutMs = 3000): Promise<{ healthy: boolean; latency_ms: number }> => {
  let targetUrl = rawUrl.trim()
  if (!/^https?:\/\//i.test(targetUrl)) {
    targetUrl = 'https://' + targetUrl
  }

  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs)
  const startTime = performance.now()

  try {
    // 使用 no-cors 模式向目标地址发起轻量探测
    await fetch(targetUrl, {
      method: 'GET',
      mode: 'no-cors',
      cache: 'no-store',
      signal: controller.signal,
    })
    clearTimeout(timeoutId)
    const latency = Math.round(performance.now() - startTime)
    return { healthy: true, latency_ms: latency }
  } catch {
    clearTimeout(timeoutId)
    try {
      const imgPing = await new Promise<number>((resolve, reject) => {
        const img = new Image()
        const timer = setTimeout(() => reject('timeout'), 1500)
        img.onload = () => {
          clearTimeout(timer)
          resolve(Math.round(performance.now() - startTime))
        }
        img.onerror = () => {
          clearTimeout(timer)
          resolve(Math.round(performance.now() - startTime))
        }
        img.src = `${targetUrl.replace(/\/+$/, '')}/favicon.ico?_t=${Date.now()}`
      })
      return { healthy: true, latency_ms: imgPing }
    } catch {
      return { healthy: false, latency_ms: 0 }
    }
  }
}

// 批量轻量探测当前所有导航链接
const probeLinksHealth = async (items: LinkItem[]) => {
  if (!items || items.length === 0) return
  const tasks = items.map(async (item) => {
    const res = await pingClient(item.url)
    healthStatusMap.value[item.id] = res
  })
  await Promise.allSettled(tasks)
}

// 定时轮询定时器与可见性控制 (每 30 秒静默更新一次)
let probeTimer: ReturnType<typeof setInterval> | null = null
const PROBE_INTERVAL_MS = 30000

const startPeriodicProbe = () => {
  stopPeriodicProbe()
  probeTimer = setInterval(() => {
    if (document.visibilityState === 'visible' && links.value && links.value.length > 0) {
      probeLinksHealth(links.value)
    }
  }, PROBE_INTERVAL_MS)
}

const stopPeriodicProbe = () => {
  if (probeTimer) {
    clearInterval(probeTimer)
    probeTimer = null
  }
}

// 当用户切走标签页时暂停测速，切回前台时立即刷新一次
const handleVisibilityChange = () => {
  if (document.visibilityState === 'visible') {
    if (links.value && links.value.length > 0) {
      probeLinksHealth(links.value)
    }
    startPeriodicProbe()
  } else {
    stopPeriodicProbe()
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

// 🌟 为无 Favicon 的首字母卡片生成 Vercel 风格的高级柔和微晶色彩
const getInitialColorClass = (title: string, category = '') => {
  const str = (title + category).toLowerCase()
  if (str.includes('vue')) return 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/25'
  if (str.includes('tailwind') || str.includes('css')) return 'bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/25'
  if (str.includes('go') || str.includes('dev')) return 'bg-cyan-500/10 text-cyan-600 dark:text-cyan-400 border-cyan-500/25'
  if (str.includes('react') || str.includes('front')) return 'bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/25'
  if (str.includes('linear') || str.includes('design') || str.includes('figma')) return 'bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/25'
  if (str.includes('cloud') || str.includes('ops') || str.includes('deploy')) return 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/25'
  
  const hash = str.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  const palettes = [
    'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/25',
    'bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/25',
    'bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border-indigo-500/25',
    'bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/25',
    'bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/25',
    'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/25',
    'bg-teal-500/10 text-teal-600 dark:text-teal-400 border-teal-500/25',
  ]
  return palettes[hash % palettes.length]
}

// 计算所有分类
const categories = computed(() => {
  const cats = new Set<string>()
  if (links.value) {
    links.value.forEach(link => {
      if (link.category) cats.add(link.category)
    })
  }
  return Array.from(cats)
})

// 根据搜索和分类过滤链接
const filteredLinks = computed(() => {
  if (!links.value) return []
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

const { siteConfig, loadSiteConfig } = useSiteConfig()

// 🌟 公告 3 行多页滑动与自动轮播系统
const ANNOUNCE_PAGE_SIZE = 3
const currentAnnouncePage = ref(0)
const isAnnounceHovered = ref(false)

const announcementPages = computed(() => {
  const pages: AnnouncementItem[][] = []
  for (let i = 0; i < announcements.value.length; i += ANNOUNCE_PAGE_SIZE) {
    pages.push(announcements.value.slice(i, i + ANNOUNCE_PAGE_SIZE))
  }
  return pages
})

const totalAnnouncePages = computed(() => announcementPages.value.length || 1)

const nextAnnouncePage = () => {
  if (totalAnnouncePages.value <= 1) return
  currentAnnouncePage.value = (currentAnnouncePage.value + 1) % totalAnnouncePages.value
}

const prevAnnouncePage = () => {
  if (totalAnnouncePages.value <= 1) return
  currentAnnouncePage.value = (currentAnnouncePage.value - 1 + totalAnnouncePages.value) % totalAnnouncePages.value
}

// 自动翻页定时器 (6 秒自动切页，Hover 时智能暂停)
let announceTimer: ReturnType<typeof setInterval> | null = null
const startAnnounceAutoPlay = () => {
  stopAnnounceAutoPlay()
  announceTimer = setInterval(() => {
    if (!isAnnounceHovered.value && totalAnnouncePages.value > 1) {
      nextAnnouncePage()
    }
  }, 6000)
}
const stopAnnounceAutoPlay = () => {
  if (announceTimer) {
    clearInterval(announceTimer)
    announceTimer = null
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
  loadSiteConfig()
  fetchAnnouncements()
  fetchLinks()
  startPeriodicProbe()
  startAnnounceAutoPlay()
  window.addEventListener('keydown', handleKeyDown)
  document.addEventListener('visibilitychange', handleVisibilityChange)
})

onUnmounted(() => {
  stopPeriodicProbe()
  stopAnnounceAutoPlay()
  window.removeEventListener('keydown', handleKeyDown)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})
</script>

<template>
  <div class="space-y-6 sm:space-y-8">
    <!-- 🌟 顶部 3 行垂直分页轮播公告栏 (无卡片包裹 · 纯净通栏 · 自动轮播与分页) -->
    <section 
      v-if="announcements.length > 0" 
      class="py-1 transition-all duration-200"
      @mouseenter="isAnnounceHovered = true"
      @mouseleave="isAnnounceHovered = false"
    >
      <div class="space-y-2">
        <!-- 3 行垂直列表 (无卡片包裹，纯净通栏，带平滑淡入淡出动效) -->
        <transition 
          mode="out-in" 
          enter-active-class="transition duration-200 ease-out" 
          enter-from-class="opacity-0 translate-y-1" 
          enter-to-class="opacity-100 translate-y-0" 
          leave-active-class="transition duration-150 ease-in" 
          leave-from-class="opacity-100 translate-y-0" 
          leave-to-class="opacity-0 -translate-y-1"
        >
          <div :key="currentAnnouncePage" class="space-y-2.5">
            <div 
              v-for="item in (announcementPages[currentAnnouncePage] || [])" 
              :key="item.id" 
              class="flex items-center justify-between gap-3 group py-0.5"
            >
              <!-- 左侧：加粗公告标签 + 标题 -->
              <div class="flex items-center space-x-3 min-w-0 flex-1">
                <span class="px-2 py-0.5 text-xs font-bold rounded bg-secondary text-foreground shrink-0 select-none">
                  公告
                </span>

                <router-link
                  :to="'/docs/' + item.id"
                  class="tracking-tight text-xs sm:text-[13px] font-medium text-foreground/90 group-hover:text-foreground group-hover:underline transition-colors truncate cursor-pointer font-sans"
                  :title="'查看详情: ' + item.content"
                >
                  {{ item.content }}
                </router-link>
              </div>

              <!-- 右侧详情箭头 -->
              <router-link
                :to="'/docs/' + item.id"
                class="opacity-0 group-hover:opacity-100 transition-opacity text-xs text-muted-foreground group-hover:text-foreground shrink-0 font-mono flex items-center space-x-0.5"
                title="查看详情"
              >
                <span>详情</span>
                <svg class="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                </svg>
              </router-link>
            </div>
          </div>
        </transition>

        <!-- 底部微型分页控制器 (当页数 > 1 时呈现，纯净右对齐) -->
        <div v-if="totalAnnouncePages > 1" class="flex items-center justify-end pt-1 space-x-2 text-xs text-muted-foreground">
          <!-- 页码圆点/指示条 -->
          <div class="flex items-center space-x-1">
            <button
              v-for="(_, pIdx) in announcementPages"
              :key="pIdx"
              @click="currentAnnouncePage = pIdx"
              class="h-1.5 rounded-full transition-all duration-200 cursor-pointer"
              :class="currentAnnouncePage === pIdx ? 'w-4 bg-foreground' : 'w-1.5 bg-border hover:bg-muted-foreground/60'"
              :title="`跳转到第 ${pIdx + 1} 页`"
            ></button>
          </div>

          <!-- 上一页 / 下一页微控件 -->
          <div class="flex items-center space-x-1 font-mono">
            <button 
              @click="prevAnnouncePage" 
              class="w-5 h-5 rounded flex items-center justify-center hover:bg-accent text-foreground/80 hover:text-foreground cursor-pointer text-xs transition-colors"
              title="上一页"
            >
              ‹
            </button>
            <span class="text-[11px] text-muted-foreground/80 select-none">{{ currentAnnouncePage + 1 }}/{{ totalAnnouncePages }}</span>
            <button 
              @click="nextAnnouncePage" 
              class="w-5 h-5 rounded flex items-center justify-center hover:bg-accent text-foreground/80 hover:text-foreground cursor-pointer text-xs transition-colors"
              title="下一页"
            >
              ›
            </button>
          </div>
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
          全部应用 ({{ links ? links.length : 0 }})
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
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-5">
        <div 
          v-for="n in 8" 
          :key="n" 
          class="h-28 rounded-xl border border-border/60 bg-muted/20 animate-pulse"
        ></div>
      </div>

      <!-- 链接大卡片网格 (支持 Favicon 高清渲染 + 柔和晶体色彩回退 + 健康探测微呼吸点) -->
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
          class="group relative flex items-start p-5 sm:p-5.5 rounded-xl border border-border/80 bg-card text-card-foreground hover:border-zinc-300 dark:hover:border-zinc-700 hover:-translate-y-0.5 hover:shadow-md hover:shadow-zinc-200/40 dark:hover:shadow-zinc-950/50 transition-all duration-200 ease-out cursor-pointer"
        >
          <!-- 🌟 左侧 Favicon 真实图标槽 (带智能柔和微晶色彩优雅降级) -->
          <div 
            class="relative w-11 h-11 rounded-lg flex items-center justify-center border shrink-0 group-hover:scale-105 transition-transform duration-200 select-none overflow-hidden p-1.5"
            :class="link.icon ? 'bg-secondary/70 border-border/50' : getInitialColorClass(link.title, link.category)"
          >
            <img 
              v-if="link.icon" 
              :src="link.icon" 
              :alt="link.title" 
              class="w-full h-full object-contain"
              @error="handleIconError"
            />
            <div 
              class="w-full h-full items-center justify-center font-bold text-base"
              :class="link.icon ? 'hidden' : 'flex'"
            >
              {{ getInitial(link.title) }}
            </div>

            <!-- 📶 健康度探测微状态指示点 (带精致微光呼吸辉光) -->
            <span 
              v-if="healthStatusMap[link.id]"
              class="absolute top-1 right-1 w-2 h-2 rounded-full border border-background"
              :class="healthStatusMap[link.id].healthy ? 'bg-emerald-500 shadow-[0_0_6px_rgba(16,185,129,0.7)]' : 'bg-rose-500 shadow-[0_0_6px_rgba(244,63,94,0.7)]'"
              :title="healthStatusMap[link.id].healthy ? `服务正常 · 延迟: ${healthStatusMap[link.id].latency_ms}ms` : '服务暂时不可达'"
            ></span>
          </div>

          <!-- 中间标题与信息 -->
          <div class="ml-4 flex-1 min-w-0 pr-2">
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold tracking-tight text-foreground truncate group-hover:text-foreground transition-colors">
                {{ link.title }}
              </h3>
            </div>
            
            <p class="text-xs text-muted-foreground/80 truncate mt-1 font-mono">
              {{ extractHostname(link.url) }}
            </p>

            <!-- 底部：分类 (左) + 实时延迟健康指示 (右) -->
            <div class="mt-3.5 pt-2 border-t border-border/40 flex items-center justify-between text-xs">
              <span class="text-[11px] text-muted-foreground font-medium group-hover:text-foreground transition-colors truncate">
                {{ link.category || '默认' }}
              </span>

              <div 
                v-if="healthStatusMap[link.id]?.healthy"
                class="inline-flex items-center space-x-1.5 text-[10px] font-mono text-muted-foreground/80 shrink-0"
                :title="`连通正常 · 响应时间 ${healthStatusMap[link.id].latency_ms}ms`"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_5px_rgba(16,185,129,0.7)] shrink-0"></span>
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
