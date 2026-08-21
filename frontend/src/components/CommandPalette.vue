<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'

interface LinkItem {
  id: number
  title: string
  url: string
  category: string
  icon?: string
}

const props = defineProps<{
  links: LinkItem[]
}>()

const isOpen = ref(false)
const search = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)
const listRef = ref<HTMLDivElement | null>(null)

// 提取域名
const extractHostname = (url: string) => {
  try {
    const u = new URL(url.startsWith('http') ? url : `https://${url}`)
    return u.hostname.replace(/^www\./, '')
  } catch {
    return url
  }
}

// 格式化 URL
const formatUrl = (url: string) => {
  if (!url) return '#'
  if (/^https?:\/\//i.test(url)) return url
  return `https://${url}`
}

// 过滤搜索
const filteredList = computed(() => {
  const query = search.value.trim().toLowerCase()
  if (!query) {
    return props.links.slice(0, 12)
  }
  return props.links.filter(item => {
    return (
      item.title.toLowerCase().includes(query) ||
      item.url.toLowerCase().includes(query) ||
      (item.category && item.category.toLowerCase().includes(query))
    )
  })
})

// 监听搜索词变化重置选中项
watch(search, () => {
  selectedIndex.value = 0
})

// 打开面板
const open = () => {
  isOpen.value = true
  search.value = ''
  selectedIndex.value = 0
  nextTick(() => {
    inputRef.value?.focus()
  })
}

// 关闭面板
const close = () => {
  isOpen.value = false
}

// 打开链接
const selectItem = (item: LinkItem) => {
  window.open(formatUrl(item.url), '_blank', 'noopener,noreferrer')
  close()
}

// 键盘控制
const handleKeyDown = (e: KeyboardEvent) => {
  // 快捷键呼出
  if (((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') || (e.key === '/' && !isOpen.value && document.activeElement?.tagName !== 'INPUT' && document.activeElement?.tagName !== 'TEXTAREA')) {
    e.preventDefault()
    if (isOpen.value) {
      close()
    } else {
      open()
    }
    return
  }

  if (!isOpen.value) return

  if (e.key === 'Escape') {
    e.preventDefault()
    close()
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (filteredList.value.length > 0) {
      selectedIndex.value = (selectedIndex.value + 1) % filteredList.value.length
      scrollToSelected()
    }
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (filteredList.value.length > 0) {
      selectedIndex.value = (selectedIndex.value - 1 + filteredList.value.length) % filteredList.value.length
      scrollToSelected()
    }
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (filteredList.value.length > 0 && filteredList.value[selectedIndex.value]) {
      selectItem(filteredList.value[selectedIndex.value])
    }
  }
}

const scrollToSelected = () => {
  nextTick(() => {
    const el = listRef.value?.querySelector(`[data-index="${selectedIndex.value}"]`)
    if (el) {
      el.scrollIntoView({ block: 'nearest' })
    }
  })
}

// 处理图片加载错误降级
const handleImgError = (e: Event) => {
  const target = e.target as HTMLElement
  target.style.display = 'none'
  if (target.nextElementSibling) {
    (target.nextElementSibling as HTMLElement).style.display = 'flex'
  }
}

defineExpose({
  open,
  close,
})

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
  <!-- 全局 Command ⌘K 搜索浮层 (Raycast / shadcn 风格) -->
  <teleport to="body">
    <transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div 
        v-if="isOpen" 
        class="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] px-4 bg-black/50 backdrop-blur-sm"
        @click.self="close"
      >
        <div 
          class="w-full max-w-xl bg-popover text-popover-foreground rounded-xl border border-border shadow-2xl overflow-hidden animate-in fade-in-0 zoom-in-95 duration-150"
        >
          <!-- 搜索输入头 -->
          <div class="flex items-center px-3.5 border-b border-border">
            <svg class="w-4 h-4 text-muted-foreground mr-2.5 shrink-0" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>

            <input
              ref="inputRef"
              v-model="search"
              type="text"
              placeholder="搜索应用、文档或分类 (输入并按 Enter 直达)..."
              class="w-full h-12 bg-transparent text-sm placeholder:text-muted-foreground focus:outline-none"
            />

            <kbd class="text-[10px] font-mono font-medium text-muted-foreground bg-muted border border-border px-1.5 py-0.5 rounded select-none shrink-0">
              ESC
            </kbd>
          </div>

          <!-- 搜索结果列表 -->
          <div ref="listRef" class="max-h-80 overflow-y-auto p-1.5 divide-y divide-border/40">
            <div v-if="filteredList.length > 0" class="space-y-0.5">
              <div
                v-for="(item, index) in filteredList"
                :key="item.id"
                :data-index="index"
                @click="selectItem(item)"
                @mouseenter="selectedIndex = index"
                :class="[
                  'group flex items-center justify-between px-3 py-2.5 rounded-lg cursor-pointer transition-colors text-sm select-none',
                  selectedIndex === index 
                    ? 'bg-accent text-accent-foreground font-medium' 
                    : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'
                ]"
              >
                <!-- 左侧: Favicon / 首字母 + 标题 + 分类 -->
                <div class="flex items-center space-x-3 min-w-0">
                  <div class="w-6 h-6 rounded-md bg-secondary flex items-center justify-center overflow-hidden shrink-0 border border-border/40">
                    <img 
                      v-if="item.icon" 
                      :src="item.icon" 
                      :alt="item.title"
                      class="w-4 h-4 object-contain"
                      @error="handleImgError"
                    />
                    <span class="text-[11px] font-bold text-foreground hidden" :class="{ '!flex': !item.icon }">
                      {{ item.title ? item.title.charAt(0).toUpperCase() : 'N' }}
                    </span>
                  </div>

                  <span class="text-sm font-semibold truncate text-foreground">
                    {{ item.title }}
                  </span>

                  <span class="text-xs font-mono text-muted-foreground truncate hidden sm:inline-block">
                    {{ extractHostname(item.url) }}
                  </span>
                </div>

                <!-- 右侧: 分类小徽标 + 回车直达提示 -->
                <div class="flex items-center space-x-2 shrink-0 ml-3">
                  <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-secondary text-muted-foreground">
                    {{ item.category || '默认' }}
                  </span>

                  <span 
                    v-if="selectedIndex === index"
                    class="text-[10px] font-mono text-muted-foreground bg-background border border-border px-1.5 py-0.5 rounded shadow-2xs"
                  >
                    ↵ 直达
                  </span>
                </div>
              </div>
            </div>

            <!-- 空结果 -->
            <div v-else class="py-10 text-center text-xs text-muted-foreground">
              未找到匹配“{{ search }}”的导航项
            </div>
          </div>

          <!-- 底部快捷键操作提示 -->
          <div class="flex items-center justify-between px-3.5 py-2 bg-muted/30 border-t border-border text-[11px] font-mono text-muted-foreground">
            <div class="flex items-center space-x-3">
              <span><kbd class="bg-muted px-1 py-0.5 rounded border border-border text-[10px]">↑</kbd> <kbd class="bg-muted px-1 py-0.5 rounded border border-border text-[10px]">↓</kbd> 切换导航</span>
              <span><kbd class="bg-muted px-1 py-0.5 rounded border border-border text-[10px]">↵</kbd> 新标签打开</span>
            </div>
            <span>共 {{ links.length }} 项</span>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>
