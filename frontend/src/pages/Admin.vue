<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface LinkItem {
  id: number
  title: string
  url: string
  category: string
  icon?: string
}

interface AnnouncementItem {
  id: number
  content: string
  is_active: boolean
  created_at: string
}

const isAuthenticated = ref(false)
const inputPassword = ref('')
const authError = ref('')
const links = ref<LinkItem[]>([])
const announcements = ref<AnnouncementItem[]>([])

const newLink = ref({ title: '', url: '', category: '开发协作', icon: '' })
const newAnnouncement = ref({ content: '', is_active: true })
const activeTab = ref<'links' | 'announcements' | 'backup'>('links')
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const faviconLoading = ref(false)
const pingLoading = ref(false)
const pingResults = ref<Record<number, { healthy: boolean; latency_ms: number; status_code: number; error: string }>>({})

const fileInputRef = ref<HTMLInputElement | null>(null)
const jsonInputRef = ref<HTMLInputElement | null>(null)

const showMessage = (text: string, type: 'success' | 'error' = 'success') => {
  message.value = { type, text }
  setTimeout(() => {
    message.value = null
  }, 3500)
}

// 获取 Auth Token
const getToken = () => localStorage.getItem('admin_token') || ''

// 检查登录状态
const checkAuth = async () => {
  const token = getToken()
  if (!token) {
    isAuthenticated.value = false
    return
  }
  try {
    const res = await fetch('/api/auth/check', {
      headers: { 'X-Admin-Token': token }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0) {
        isAuthenticated.value = true
        loadData()
      } else {
        isAuthenticated.value = false
      }
    } else {
      isAuthenticated.value = false
    }
  } catch {
    isAuthenticated.value = false
  }
}

// 登录口令解锁
const handleLogin = async () => {
  if (!inputPassword.value) {
    authError.value = '请输入口令'
    return
  }
  authError.value = ''
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: inputPassword.value }),
    })
    const data = await res.json()
    if (data.code === 0 && data.data?.token) {
      localStorage.setItem('admin_token', data.data.token)
      isAuthenticated.value = true
      inputPassword.value = ''
      showMessage('管理权限已解锁')
      loadData()
    } else {
      authError.value = data.msg || '口令错误，请重试'
    }
  } catch {
    authError.value = '网络错误，请稍后重试'
  }
}

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('admin_token')
  isAuthenticated.value = false
  showMessage('已退出管理后台')
}

// 刷新数据
const loadData = async () => {
  try {
    const [linksRes, annRes] = await Promise.all([
      fetch('/api/links'),
      fetch('/api/announcements')
    ])
    if (linksRes.ok) {
      const data = await linksRes.json()
      if (data.code === 0) links.value = data.data || []
    }
    if (annRes.ok) {
      const data = await annRes.json()
      if (data.code === 0) announcements.value = data.data || []
    }
  } catch {
    showMessage('获取数据失败', 'error')
  }
}

// 自动提取目标网址 Favicon
const handleAutoFavicon = async () => {
  if (!newLink.value.url) {
    showMessage('请先输入目标网址', 'error')
    return
  }
  faviconLoading.value = true
  try {
    const res = await fetch(`/api/tools/favicon?url=${encodeURIComponent(newLink.value.url)}`)
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0 && data.data?.favicon) {
        newLink.value.icon = data.data.favicon
        showMessage('已自动获取并匹配高清 Favicon 图标')
      }
    }
  } catch {
    showMessage('自动获取图标失败', 'error')
  } finally {
    faviconLoading.value = false
  }
}

// 添加链接
const handleAddLink = async () => {
  if (!newLink.value.title || !newLink.value.url) {
    showMessage('请填写完整链接信息', 'error')
    return
  }
  try {
    const res = await fetch('/api/links', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'X-Admin-Token': getToken()
      },
      body: JSON.stringify(newLink.value),
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('链接添加成功')
      newLink.value = { title: '', url: '', category: '开发协作', icon: '' }
      loadData()
    } else {
      showMessage(data.msg || '添加失败', 'error')
    }
  } catch {
    showMessage('网络错误', 'error')
  }
}

// 删除链接
const handleDeleteLink = async (id: number) => {
  if (!confirm('确认删除该链接？')) return
  try {
    const res = await fetch(`/api/links/${id}`, { 
      method: 'DELETE',
      headers: { 'X-Admin-Token': getToken() }
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('链接已删除')
      loadData()
    } else {
      showMessage(data.msg || '删除失败', 'error')
    }
  } catch {
    showMessage('删除失败', 'error')
  }
}

// 添加公告
const handleAddAnnouncement = async () => {
  if (!newAnnouncement.value.content) {
    showMessage('请输入公告内容', 'error')
    return
  }
  try {
    const res = await fetch('/api/announcements', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'X-Admin-Token': getToken()
      },
      body: JSON.stringify(newAnnouncement.value),
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('公告发布成功')
      newAnnouncement.value = { content: '', is_active: true }
      loadData()
    } else {
      showMessage(data.msg || '发布失败', 'error')
    }
  } catch {
    showMessage('网络错误', 'error')
  }
}

// 切换公告状态
const handleToggleAnnouncement = async (id: number) => {
  try {
    const res = await fetch(`/api/announcements/${id}/toggle`, { 
      method: 'PUT',
      headers: { 'X-Admin-Token': getToken() }
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('状态已更新')
      loadData()
    }
  } catch {
    showMessage('操作失败', 'error')
  }
}

// 删除公告
const handleDeleteAnnouncement = async (id: number) => {
  if (!confirm('确认删除此公告？')) return
  try {
    const res = await fetch(`/api/announcements/${id}`, { 
      method: 'DELETE',
      headers: { 'X-Admin-Token': getToken() }
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('公告已删除')
      loadData()
    } else {
      showMessage(data.msg || '删除失败', 'error')
    }
  } catch {
    showMessage('删除失败', 'error')
  }
}

// 全量连通性 Ping 检测
const handleTestAllPing = async () => {
  pingLoading.value = true
  pingResults.value = {}
  for (const item of links.value) {
    try {
      const res = await fetch('/api/tools/ping', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: item.url }),
      })
      if (res.ok) {
        const data = await res.json()
        if (data.code === 0 && data.data) {
          pingResults.value[item.id] = data.data
        }
      }
    } catch {
      pingResults.value[item.id] = { healthy: false, latency_ms: 0, status_code: 0, error: '网络中断' }
    }
  }
  pingLoading.value = false
  showMessage('全部链接连通性探测完毕')
}

// 导出 JSON 备份
const handleExportBackup = () => {
  window.open(`/api/backup/export?token=${getToken()}`, '_blank')
}

// 上传 Chrome 书签文件
const handleBookmarkFileChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return
  const file = target.files[0]

  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await fetch('/api/backup/import-bookmarks', {
      method: 'POST',
      headers: { 'X-Admin-Token': getToken() },
      body: formData,
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage(`书签导入成功，共批量录入 ${data.data.imported_count} 条导航链接！`)
      loadData()
    } else {
      showMessage(data.msg || '书签解析导入失败', 'error')
    }
  } catch {
    showMessage('上传书签失败', 'error')
  } finally {
    if (fileInputRef.value) fileInputRef.value.value = ''
  }
}

// 上传 JSON 备份恢复
const handleJsonBackupChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return
  const file = target.files[0]

  const reader = new FileReader()
  reader.onload = async (event) => {
    try {
      const content = event.target?.result as string
      const parsed = JSON.parse(content)
      const res = await fetch('/api/backup/import', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Admin-Token': getToken()
        },
        body: JSON.stringify({
          mode: 'merge',
          links: parsed.links || [],
          announcements: parsed.announcements || []
        }),
      })
      const data = await res.json()
      if (data.code === 0) {
        showMessage(`JSON 备份还原成功！恢复了 ${data.data.imported_links} 条链接`)
        loadData()
      } else {
        showMessage(data.msg || '恢复失败', 'error')
      }
    } catch {
      showMessage('JSON 文件解析失败，格式不正确', 'error')
    } finally {
      if (jsonInputRef.value) jsonInputRef.value.value = ''
    }
  }
  reader.readAsText(file)
}

onMounted(() => {
  checkAuth()
})
</script>

<template>
  <div class="space-y-10 max-w-5xl mx-auto">
    <!-- 🔒 未登录状态：口令解锁面板 -->
    <div v-if="!isAuthenticated" class="py-12 flex items-center justify-center">
      <div class="w-full max-w-md p-8 rounded-xl border border-border bg-card shadow-lg space-y-6 text-center">
        <!-- 锁图标 -->
        <div class="w-12 h-12 rounded-xl bg-secondary flex items-center justify-center mx-auto text-foreground border border-border/60">
          <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
            <path stroke-linecap="round" stroke-linejoin="round" d="M7 11V7a5 5 0 0110 0v4" />
          </svg>
        </div>

        <div class="space-y-1.5">
          <h2 class="text-xl font-bold tracking-tight text-foreground">管理后台口令验证</h2>
          <p class="text-xs text-muted-foreground">请输入管理员口令解锁导航链接与公告维护权限 (默认: admin123)</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4 text-left">
          <div class="space-y-1.5">
            <input 
              v-model="inputPassword"
              type="password"
              placeholder="输入管理员口令..." 
              class="w-full h-11 bg-background border border-border text-sm rounded-lg px-3.5 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
              autofocus
              required
            />
            <p v-if="authError" class="text-xs text-rose-500 font-medium">{{ authError }}</p>
          </div>

          <button 
            type="submit" 
            class="w-full h-11 bg-foreground text-background hover:opacity-90 rounded-lg text-sm font-semibold transition-opacity cursor-pointer shadow-sm"
          >
            解锁控制台
          </button>
        </form>
      </div>
    </div>

    <!-- 🔓 已解锁管理状态：完整控制台 -->
    <div v-else class="space-y-10">
      <!-- 顶部标题与控制栏 -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <h2 class="text-2xl sm:text-3xl font-bold tracking-tight text-foreground">后台资源管理</h2>
          <p class="text-sm text-muted-foreground mt-1">配置团队导航链接、公告发布、健康探测与数据备份</p>
        </div>

        <div class="flex items-center space-x-3 self-start sm:self-auto">
          <!-- 选项卡切换 -->
          <div class="flex items-center space-x-1 bg-secondary/60 p-1 rounded-lg border border-border/50">
            <button
              @click="activeTab = 'links'"
              :class="[
                'text-xs sm:text-sm px-3.5 py-1.5 rounded-md font-medium transition-all cursor-pointer',
                activeTab === 'links' 
                  ? 'bg-card text-foreground shadow-xs font-semibold' 
                  : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              导航链接 ({{ links.length }})
            </button>
            <button
              @click="activeTab = 'announcements'"
              :class="[
                'text-xs sm:text-sm px-3.5 py-1.5 rounded-md font-medium transition-all cursor-pointer',
                activeTab === 'announcements' 
                  ? 'bg-card text-foreground shadow-xs font-semibold' 
                  : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              系统公告 ({{ announcements.length }})
            </button>
            <button
              @click="activeTab = 'backup'"
              :class="[
                'text-xs sm:text-sm px-3.5 py-1.5 rounded-md font-medium transition-all cursor-pointer',
                activeTab === 'backup' 
                  ? 'bg-card text-foreground shadow-xs font-semibold' 
                  : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              备份与导入
            </button>
          </div>

          <!-- 退出登录按钮 -->
          <button 
            @click="handleLogout"
            class="text-xs text-muted-foreground hover:text-rose-500 border border-border/80 px-2.5 py-1.5 rounded-md hover:bg-rose-500/10 transition-colors cursor-pointer"
            title="锁定管理后台"
          >
            退出
          </button>
        </div>
      </div>

      <!-- 全局提示条 -->
      <div 
        v-if="message" 
        :class="[
          'px-4 py-2.5 text-sm rounded-lg border transition-all duration-200 shadow-sm flex items-center space-x-2',
          message.type === 'success' 
            ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-600 dark:text-emerald-400'
            : 'bg-rose-500/10 border-rose-500/30 text-rose-500'
        ]"
      >
        <span>{{ message.text }}</span>
      </div>

      <!-- 1. 链接管理面板 -->
      <div v-if="activeTab === 'links'" class="space-y-8">
        <!-- 新增链接表单 -->
        <form @submit.prevent="handleAddLink" class="p-6 rounded-xl border border-border bg-card shadow-sm space-y-5">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-semibold text-foreground tracking-tight">添加新导航链接</h3>
            <span class="text-xs text-muted-foreground">支持自动抓取 Favicon 网站图标</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-muted-foreground">链接名称</label>
              <input 
                v-model="newLink.title" 
                placeholder="例如: Figma, Linear, GitHub" 
                class="w-full h-10 bg-background border border-border text-sm rounded-lg px-3 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
                required
              />
            </div>
            <div class="space-y-1.5">
              <div class="flex items-center justify-between">
                <label class="text-xs font-medium text-muted-foreground">目标网址</label>
                <button
                  type="button"
                  @click="handleAutoFavicon"
                  class="text-[11px] font-mono text-blue-600 dark:text-blue-400 hover:underline cursor-pointer"
                >
                  {{ faviconLoading ? '抓取中...' : '自动抓取图标' }}
                </button>
              </div>
              <input 
                v-model="newLink.url" 
                placeholder="https://example.com" 
                class="w-full h-10 bg-background border border-border text-sm rounded-lg px-3 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
                @blur="!newLink.icon && newLink.url && handleAutoFavicon()"
                required
              />
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-muted-foreground">所属分类</label>
              <input 
                v-model="newLink.category" 
                placeholder="例如: 开发协作, 部署运维, 设计资源" 
                class="w-full h-10 bg-background border border-border text-sm rounded-lg px-3 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
              />
            </div>
          </div>

          <!-- 图标预览槽 -->
          <div v-if="newLink.icon" class="flex items-center space-x-3 p-2.5 rounded-lg bg-secondary/40 border border-border/50 text-xs">
            <span class="text-muted-foreground">图标预览:</span>
            <img :src="newLink.icon" alt="icon" class="w-5 h-5 object-contain" />
            <span class="font-mono text-muted-foreground truncate max-w-xs">{{ newLink.icon }}</span>
          </div>

          <div class="flex justify-end pt-1">
            <button 
              type="submit" 
              class="bg-foreground text-background hover:opacity-90 text-sm px-6 py-2 rounded-lg font-medium transition-opacity cursor-pointer shadow-sm"
            >
              保存并添加
            </button>
          </div>
        </form>

        <!-- 链接列表与连通性检测 -->
        <div class="rounded-xl border border-border bg-card overflow-hidden shadow-sm">
          <div class="px-6 py-4 border-b border-border flex items-center justify-between bg-muted/20">
            <h3 class="text-sm font-semibold text-foreground">全部链接列表 ({{ links.length }})</h3>
            <button
              @click="handleTestAllPing"
              :disabled="pingLoading"
              class="text-xs border border-border bg-background hover:bg-accent px-3 py-1.5 rounded-md font-mono flex items-center space-x-1.5 transition-colors cursor-pointer"
            >
              <span class="w-2 h-2 rounded-full" :class="pingLoading ? 'bg-amber-500 animate-spin' : 'bg-emerald-500'"></span>
              <span>{{ pingLoading ? '正在全量探测...' : '一键探测连通性' }}</span>
            </button>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left text-sm">
              <thead class="bg-muted/40 text-xs font-semibold text-muted-foreground border-b border-border">
                <tr>
                  <th class="px-6 py-3">图标 / 标题</th>
                  <th class="px-6 py-3">目标网址</th>
                  <th class="px-6 py-3">分类</th>
                  <th class="px-6 py-3">连通健康度</th>
                  <th class="px-6 py-3 text-right">操作</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr v-for="link in links" :key="link.id" class="hover:bg-muted/30 transition-colors">
                  <td class="px-6 py-3.5 font-medium text-foreground flex items-center space-x-3">
                    <div class="w-7 h-7 rounded bg-secondary flex items-center justify-center border border-border/40 shrink-0 overflow-hidden">
                      <img v-if="link.icon" :src="link.icon" class="w-4 h-4 object-contain" />
                      <span v-else class="text-xs font-bold text-foreground">{{ link.title.charAt(0).toUpperCase() }}</span>
                    </div>
                    <span>{{ link.title }}</span>
                  </td>
                  <td class="px-6 py-3.5 font-mono text-xs text-muted-foreground truncate max-w-xs">
                    <a :href="link.url" target="_blank" class="hover:underline hover:text-foreground">
                      {{ link.url }}
                    </a>
                  </td>
                  <td class="px-6 py-3.5">
                    <span class="px-2 py-0.5 rounded text-[11px] font-medium bg-secondary text-muted-foreground">
                      {{ link.category || '默认' }}
                    </span>
                  </td>
                  <td class="px-6 py-3.5 font-mono text-xs">
                    <span v-if="pingResults[link.id]" class="inline-flex items-center space-x-1.5">
                      <span class="w-2 h-2 rounded-full" :class="pingResults[link.id].healthy ? 'bg-emerald-500' : 'bg-rose-500'"></span>
                      <span :class="pingResults[link.id].healthy ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600'">
                        {{ pingResults[link.id].healthy ? `${pingResults[link.id].latency_ms}ms (HTTP ${pingResults[link.id].status_code})` : '不可达' }}
                      </span>
                    </span>
                    <span v-else class="text-muted-foreground/60">未探测</span>
                  </td>
                  <td class="px-6 py-3.5 text-right">
                    <button 
                      @click="handleDeleteLink(link.id)"
                      class="text-xs text-rose-500 hover:text-rose-700 hover:underline cursor-pointer font-medium"
                    >
                      删除
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 2. 系统公告管理面板 -->
      <div v-if="activeTab === 'announcements'" class="space-y-8">
        <!-- 发布公告表单 -->
        <form @submit.prevent="handleAddAnnouncement" class="p-6 rounded-xl border border-border bg-card shadow-sm space-y-4">
          <h3 class="text-base font-semibold text-foreground tracking-tight">发布全站系统公告</h3>

          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">公告内容</label>
            <input 
              v-model="newAnnouncement.content" 
              placeholder="例如: 本周六凌晨将进行服务器升级，期间可能出现短暂访问延迟。" 
              class="w-full h-10 bg-background border border-border text-sm rounded-lg px-3 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
              required
            />
          </div>

          <div class="flex items-center justify-between pt-1">
            <label class="flex items-center space-x-2 text-sm text-muted-foreground cursor-pointer">
              <input 
                v-model="newAnnouncement.is_active" 
                type="checkbox" 
                class="rounded border-border w-4 h-4 text-foreground focus:ring-foreground"
              />
              <span>立即在前台生效展示</span>
            </label>

            <button 
              type="submit" 
              class="bg-foreground text-background hover:opacity-90 text-sm px-6 py-2 rounded-lg font-medium transition-opacity cursor-pointer shadow-sm"
            >
              发布公告
            </button>
          </div>
        </form>

        <!-- 公告列表 -->
        <div class="rounded-xl border border-border bg-card overflow-hidden shadow-sm">
          <div class="px-6 py-4 border-b border-border bg-muted/20">
            <h3 class="text-sm font-semibold text-foreground">全部公告 ({{ announcements.length }})</h3>
          </div>

          <div class="divide-y divide-border">
            <div 
              v-for="item in announcements" 
              :key="item.id" 
              class="p-5 flex items-center justify-between gap-4 hover:bg-muted/20 transition-colors"
            >
              <div class="space-y-1">
                <div class="flex items-center space-x-2">
                  <span 
                    :class="[
                      'w-2 h-2 rounded-full',
                      item.is_active ? 'bg-emerald-500' : 'bg-zinc-400'
                    ]"
                  ></span>
                  <span class="text-xs font-mono text-muted-foreground">
                    {{ item.is_active ? '正在生效' : '已停用' }}
                  </span>
                </div>
                <p class="text-sm text-foreground font-medium">{{ item.content }}</p>
              </div>

              <div class="flex items-center space-x-3 shrink-0">
                <button 
                  @click="handleToggleAnnouncement(item.id)"
                  class="text-xs text-muted-foreground hover:text-foreground border border-border px-3 py-1.5 rounded-md hover:bg-accent transition-colors cursor-pointer"
                >
                  {{ item.is_active ? '停用' : '启用' }}
                </button>
                <button 
                  @click="handleDeleteAnnouncement(item.id)"
                  class="text-xs text-rose-500 hover:text-rose-700 hover:underline cursor-pointer"
                >
                  删除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. 数据备份与 Chrome 书签导入 -->
      <div v-if="activeTab === 'backup'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Chrome 书签导入 -->
        <div class="p-6 rounded-xl border border-border bg-card shadow-sm space-y-4">
          <div class="flex items-center space-x-2.5">
            <div class="w-8 h-8 rounded-lg bg-secondary flex items-center justify-center text-foreground">
              📑
            </div>
            <div>
              <h3 class="text-base font-semibold text-foreground">Chrome / Edge 书签导入</h3>
              <p class="text-xs text-muted-foreground">自动将书签文件夹解析为导航分类并批量入库</p>
            </div>
          </div>

          <p class="text-xs text-muted-foreground leading-relaxed">
            在浏览器书签管理器中点击「导出书签」，上传生成的 <code class="font-mono bg-secondary px-1 py-0.5 rounded">.html</code> 文件即可一键批量导入。
          </p>

          <input 
            ref="fileInputRef"
            type="file" 
            accept=".html,.htm" 
            class="hidden" 
            @change="handleBookmarkFileChange"
          />

          <button 
            @click="fileInputRef?.click()"
            class="w-full py-2.5 border border-border bg-background hover:bg-accent text-foreground text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center space-x-2 shadow-xs"
          >
            <span>选择 HTML 书签文件上传</span>
          </button>
        </div>

        <!-- JSON 备份与还原 -->
        <div class="p-6 rounded-xl border border-border bg-card shadow-sm space-y-4">
          <div class="flex items-center space-x-2.5">
            <div class="w-8 h-8 rounded-lg bg-secondary flex items-center justify-center text-foreground">
              💾
            </div>
            <div>
              <h3 class="text-base font-semibold text-foreground">全站数据 JSON 备份与恢复</h3>
              <p class="text-xs text-muted-foreground">导出完整的链接与公告数据文件，随时跨设备还原</p>
            </div>
          </div>

          <p class="text-xs text-muted-foreground leading-relaxed">
            支持一键下载完整备份 JSON，或上传历史备份文件进行合并恢复。
          </p>

          <input 
            ref="jsonInputRef"
            type="file" 
            accept=".json" 
            class="hidden" 
            @change="handleJsonBackupChange"
          />

          <div class="flex items-center space-x-3">
            <button 
              @click="handleExportBackup"
              class="flex-1 py-2.5 bg-foreground text-background hover:opacity-90 text-sm font-medium rounded-lg transition-opacity cursor-pointer shadow-xs text-center"
            >
              导出 JSON 备份
            </button>
            <button 
              @click="jsonInputRef?.click()"
              class="flex-1 py-2.5 border border-border bg-background hover:bg-accent text-foreground text-sm font-medium rounded-lg transition-colors cursor-pointer shadow-xs text-center"
            >
              导入 JSON 恢复
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
