<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSiteConfig } from '../utils/useSiteConfig'

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
  detail_md?: string
  is_active: boolean
  created_at: string
}

const { siteConfig, loadSiteConfig, updateSiteConfig } = useSiteConfig()
const editSettingsForm = ref({ site_name: '', site_desc: '', icp_beian: '' })
const saveSettingsLoading = ref(false)

const isAuthenticated = ref(false)
const inputPassword = ref('')
const authError = ref('')
const links = ref<LinkItem[]>([])
const announcements = ref<AnnouncementItem[]>([])

const newLink = ref({ title: '', url: '', category: '开发协作', icon: '' })
const newAnnouncement = ref({ content: '', detail_md: '', is_active: true })
const activeTab = ref<'links' | 'announcements' | 'backup' | 'settings'>('links')
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const faviconLoading = ref(false)
const pingLoading = ref(false)
const pingResults = ref<Record<number, { healthy: boolean; latency_ms: number; status_code: number; error: string }>>({})

// 保存系统设置
const handleSaveSettings = async () => {
  if (!editSettingsForm.value.site_name.trim()) {
    showMessage('网站名称不能为空', 'error')
    return
  }
  saveSettingsLoading.value = true
  const res = await updateSiteConfig(editSettingsForm.value, getToken())
  saveSettingsLoading.value = false
  if (res.success) {
    showMessage('网站基本信息已保存并即时生效')
  } else {
    showMessage(res.msg || '保存失败', 'error')
  }
}

// ✏️ 编辑模态框状态
const isEditLinkModalOpen = ref(false)
const editingLink = ref<{ id: number; title: string; url: string; category: string; icon: string }>({
  id: 0,
  title: '',
  url: '',
  category: '',
  icon: ''
})
const editFaviconLoading = ref(false)

const isEditAnnModalOpen = ref(false)
const editingAnnouncement = ref<{ id: number; content: string; detail_md: string; is_active: boolean }>({
  id: 0,
  content: '',
  detail_md: '',
  is_active: true
})

const fileInputRef = ref<HTMLInputElement | null>(null)
const jsonInputRef = ref<HTMLInputElement | null>(null)
const mdFileInputRef = ref<HTMLInputElement | null>(null)
const editMdFileInputRef = ref<HTMLInputElement | null>(null)

// 导入 .md 文件
const handleImportMd = (event: Event, isEdit = false) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const text = (e.target?.result as string) || ''
    if (isEdit) {
      editingAnnouncement.value.detail_md = text
      if (!editingAnnouncement.value.content) {
        const firstLine = text.split('\n')[0].replace(/^#+\s*/, '').trim()
        editingAnnouncement.value.content = firstLine || file.name.replace(/\.md$/i, '')
      }
    } else {
      newAnnouncement.value.detail_md = text
      if (!newAnnouncement.value.content) {
        const firstLine = text.split('\n')[0].replace(/^#+\s*/, '').trim()
        newAnnouncement.value.content = firstLine || file.name.replace(/\.md$/i, '')
      }
    }
    showMessage(`已成功导入 ${file.name}`)
    target.value = ''
  }
  reader.readAsText(file, 'UTF-8')
}

const showMessage = (text: string, type: 'success' | 'error' = 'success') => {
  message.value = { type, text }
  setTimeout(() => {
    message.value = null
  }, 3500)
}

const getToken = () => localStorage.getItem('admin_token') || ''

// 检查登录
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

// 登录
const handleLogin = async () => {
  if (!inputPassword.value) {
    authError.value = '请输入管理员口令'
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
    authError.value = '网络连接错误'
  }
}

// 退出
const handleLogout = () => {
  localStorage.removeItem('admin_token')
  isAuthenticated.value = false
  showMessage('已退出管理权限')
}

// 加载数据
const loadData = async () => {
  try {
    const [linksRes, annRes] = await Promise.all([
      fetch('/api/links'),
      fetch('/api/announcements'),
      loadSiteConfig()
    ])
    editSettingsForm.value = {
      site_name: siteConfig.value.site_name,
      site_desc: siteConfig.value.site_desc,
      icp_beian: siteConfig.value.icp_beian || ''
    }
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

// 自动提取 Favicon
const handleAutoFavicon = async (isEditing = false) => {
  const targetUrl = isEditing ? editingLink.value.url : newLink.value.url
  if (!targetUrl) {
    showMessage('请先输入目标网址', 'error')
    return
  }
  if (isEditing) editFaviconLoading.value = true
  else faviconLoading.value = true

  try {
    const res = await fetch(`/api/tools/favicon?url=${encodeURIComponent(targetUrl)}`)
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0 && data.data?.favicon) {
        if (isEditing) {
          editingLink.value.icon = data.data.favicon
        } else {
          newLink.value.icon = data.data.favicon
        }
        showMessage('已自动获取并匹配高清 Favicon')
      }
    }
  } catch {
    showMessage('获取图标失败', 'error')
  } finally {
    if (isEditing) editFaviconLoading.value = false
    else faviconLoading.value = false
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
      showMessage('导航链接添加成功')
      newLink.value = { title: '', url: '', category: '开发协作', icon: '' }
      loadData()
    } else {
      showMessage(data.msg || '添加失败', 'error')
    }
  } catch {
    showMessage('网络错误', 'error')
  }
}

// 打开编辑链接 Dialog
const openEditLinkModal = (link: LinkItem) => {
  editingLink.value = {
    id: link.id,
    title: link.title,
    url: link.url,
    category: link.category || '默认',
    icon: link.icon || ''
  }
  isEditLinkModalOpen.value = true
}

// 保存编辑链接
const handleSaveEditLink = async () => {
  if (!editingLink.value.title || !editingLink.value.url) {
    showMessage('标题与网址不能为空', 'error')
    return
  }
  try {
    const res = await fetch(`/api/links/${editingLink.value.id}`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'X-Admin-Token': getToken()
      },
      body: JSON.stringify(editingLink.value),
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('导航链接已更新')
      isEditLinkModalOpen.value = false
      loadData()
    } else {
      showMessage(data.msg || '更新失败', 'error')
    }
  } catch {
    showMessage('保存失败', 'error')
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
      newAnnouncement.value = { content: '', detail_md: '', is_active: true }
      loadData()
    } else {
      showMessage(data.msg || '发布失败', 'error')
    }
  } catch {
    showMessage('网络错误', 'error')
  }
}

// 打开编辑公告 Dialog
const openEditAnnModal = (item: AnnouncementItem) => {
  editingAnnouncement.value = {
    id: item.id,
    content: item.content,
    detail_md: item.detail_md || '',
    is_active: item.is_active
  }
  isEditAnnModalOpen.value = true
}

// 保存编辑公告
const handleSaveEditAnnouncement = async () => {
  if (!editingAnnouncement.value.content) {
    showMessage('公告内容不能为空', 'error')
    return
  }
  try {
    const res = await fetch(`/api/announcements/${editingAnnouncement.value.id}`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'X-Admin-Token': getToken()
      },
      body: JSON.stringify(editingAnnouncement.value),
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('公告内容已更新')
      isEditAnnModalOpen.value = false
      loadData()
    } else {
      showMessage(data.msg || '更新失败', 'error')
    }
  } catch {
    showMessage('保存失败', 'error')
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

// 上传 Chrome 书签
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
      showMessage(`书签导入成功，已批量录入 ${data.data.imported_count} 条导航`)
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

// 上传 JSON 备份
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
        showMessage(`JSON 恢复成功！还原了 ${data.data.imported_links} 条链接`)
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
  <div class="space-y-8 max-w-6xl mx-auto">
    <!-- 🔒 认证解锁面板 (正统 shadcn Card 规范) -->
    <div v-if="!isAuthenticated" class="py-16 flex items-center justify-center">
      <div class="w-full max-w-sm rounded-lg border border-border bg-card p-6 shadow-sm space-y-6">
        <div class="space-y-1.5 text-center">
          <h2 class="text-xl font-semibold tracking-tight text-foreground">管理员验证</h2>
          <p class="text-xs text-muted-foreground">请输入管理员口令解锁后台配置权限 (默认: admin123)</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">管理口令</label>
            <input 
              v-model="inputPassword"
              type="password"
              placeholder="••••••••" 
              class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              autofocus
              required
            />
            <p v-if="authError" class="text-xs text-destructive font-medium">{{ authError }}</p>
          </div>

          <button 
            type="submit" 
            class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring h-9 w-full bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 cursor-pointer"
          >
            解锁控制台
          </button>
        </form>
      </div>
    </div>

    <!-- 🔓 已解锁管理状态 (完全符合 shadcn-admin 直线排版) -->
    <div v-else class="space-y-6">
      <!-- 页面 Header: 标题 + 退出按钮 -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <h2 class="text-2xl font-bold tracking-tight text-foreground">后台资源管理</h2>
          <p class="text-xs text-muted-foreground mt-0.5">维护团队内部索引、公告通知、服务健康检测及数据备份</p>
        </div>

        <!-- 退出按钮 -->
        <button 
          @click="handleLogout"
          class="inline-flex items-center justify-center rounded-md text-xs font-medium border border-input bg-background hover:bg-accent hover:text-accent-foreground h-8 px-3 transition-colors cursor-pointer"
          title="锁定管理后台"
        >
          退出管理
        </button>
      </div>

      <!-- 🌟 纯净直线下划线导航 (彻底去除外层包裹卡片，经典瑞士排版) -->
      <div class="flex items-center space-x-6 border-b border-border text-sm">
        <button
          @click="activeTab = 'links'"
          :class="[
            'pb-2.5 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeTab === 'links'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>导航链接</span>
          <span class="text-[11px] font-mono px-1.5 py-0.5 rounded bg-muted text-muted-foreground">{{ links.length }}</span>
        </button>

        <button
          @click="activeTab = 'announcements'"
          :class="[
            'pb-2.5 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeTab === 'announcements'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>公告通知</span>
          <span class="text-[11px] font-mono px-1.5 py-0.5 rounded bg-muted text-muted-foreground">{{ announcements.length }}</span>
        </button>

        <button
          @click="activeTab = 'backup'"
          :class="[
            'pb-2.5 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeTab === 'backup'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>数据备份与导入</span>
        </button>

        <button
          @click="activeTab = 'settings'"
          :class="[
            'pb-2.5 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeTab === 'settings'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>系统设置</span>
        </button>
      </div>

      <!-- 全局消息提示条 -->
      <div 
        v-if="message" 
        :class="[
          'px-4 py-2 text-xs font-medium rounded-md border transition-all flex items-center justify-between',
          message.type === 'success' 
            ? 'bg-zinc-100 dark:bg-zinc-900 border-border text-foreground'
            : 'bg-destructive/10 border-destructive/30 text-destructive'
        ]"
      >
        <span>{{ message.text }}</span>
        <button @click="message = null" class="text-muted-foreground hover:text-foreground">✕</button>
      </div>

      <!-- 1. 导航链接管理 -->
      <div v-if="activeTab === 'links'" class="space-y-6">
        <!-- 录入新链接表单 (shadcn 规范直线表单) -->
        <div class="rounded-lg border border-border bg-card p-5 space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-foreground">录入新导航链接</h3>
            <span class="text-xs text-muted-foreground">输入网址将自动匹配网站高清图标</span>
          </div>

          <form @submit.prevent="handleAddLink" class="space-y-4">
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div class="space-y-1">
                <label class="text-xs font-medium text-muted-foreground">链接标题</label>
                <input 
                  v-model="newLink.title" 
                  placeholder="例如: GitHub, Figma" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  required
                />
              </div>

              <div class="space-y-1">
                <div class="flex items-center justify-between">
                  <label class="text-xs font-medium text-muted-foreground">目标网址</label>
                  <button
                    type="button"
                    @click="handleAutoFavicon(false)"
                    class="text-xs text-muted-foreground hover:text-foreground font-medium transition-colors cursor-pointer flex items-center space-x-1"
                  >
                    <svg class="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                    </svg>
                    <span>{{ faviconLoading ? '获取中...' : '自动获取图标' }}</span>
                  </button>
                </div>
                <input 
                  v-model="newLink.url" 
                  placeholder="https://example.com" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  @blur="!newLink.icon && newLink.url && handleAutoFavicon(false)"
                  required
                />
              </div>

              <div class="space-y-1">
                <label class="text-xs font-medium text-muted-foreground">所属分类</label>
                <input 
                  v-model="newLink.category" 
                  placeholder="例如: 开发协作, 运维部署" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                />
              </div>
            </div>

            <!-- 图标与提交 -->
            <div class="flex items-center justify-between pt-1">
              <div v-if="newLink.icon" class="flex items-center space-x-2 text-xs text-muted-foreground font-mono">
                <span>图标:</span>
                <img :src="newLink.icon" class="w-4 h-4 object-contain" />
                <span class="truncate max-w-xs">{{ newLink.icon }}</span>
              </div>
              <div v-else></div>

              <button 
                type="submit" 
                class="inline-flex items-center justify-center rounded-md text-xs sm:text-sm font-medium transition-colors h-9 px-4 bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 cursor-pointer"
              >
                + 添加链接
              </button>
            </div>
          </form>
        </div>

        <!-- 链接数据表 (正统 shadcn Table 规范) -->
        <div class="rounded-lg border border-border bg-card overflow-hidden">
          <div class="px-4 py-3 border-b border-border flex items-center justify-between bg-muted/40">
            <span class="text-xs font-semibold text-foreground">导航链接列表 ({{ links.length }})</span>
            
            <button
              @click="handleTestAllPing"
              :disabled="pingLoading"
              class="inline-flex items-center space-x-1.5 text-xs font-medium border border-input bg-background hover:bg-accent hover:text-accent-foreground px-3 py-1 rounded-md transition-colors cursor-pointer disabled:opacity-50"
            >
              <svg v-if="!pingLoading" class="w-3.5 h-3.5 text-muted-foreground" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
              </svg>
              <span v-else class="w-2 h-2 rounded-full bg-amber-500 animate-ping"></span>
              <span>{{ pingLoading ? '正在检测服务...' : '检测服务状态' }}</span>
            </button>
          </div>

          <div class="relative w-full overflow-auto">
            <table class="w-full caption-bottom text-xs">
              <thead class="[&_tr]:border-b border-border bg-muted/20">
                <tr class="border-b transition-colors">
                  <th class="h-10 px-4 text-left align-middle font-medium text-muted-foreground">图标 / 标题</th>
                  <th class="h-10 px-4 text-left align-middle font-medium text-muted-foreground">目标网址</th>
                  <th class="h-10 px-4 text-left align-middle font-medium text-muted-foreground">分类</th>
                  <th class="h-10 px-4 text-left align-middle font-medium text-muted-foreground">服务状态</th>
                  <th class="h-10 px-4 text-right align-middle font-medium text-muted-foreground">操作</th>
                </tr>
              </thead>
              <tbody class="[&_tr:last-child]:border-0 divide-y divide-border">
                <tr v-for="link in links" :key="link.id" class="border-b transition-colors hover:bg-muted/50">
                  <td class="p-4 align-middle font-medium text-foreground flex items-center space-x-2.5">
                    <div class="w-6 h-6 rounded bg-secondary flex items-center justify-center border border-border/40 shrink-0 overflow-hidden">
                      <img v-if="link.icon" :src="link.icon" class="w-3.5 h-3.5 object-contain" />
                      <span v-else class="text-[10px] font-bold text-foreground">{{ link.title.charAt(0).toUpperCase() }}</span>
                    </div>
                    <span>{{ link.title }}</span>
                  </td>

                  <td class="p-4 align-middle font-mono text-muted-foreground truncate max-w-xs">
                    <a :href="link.url" target="_blank" class="hover:underline hover:text-foreground">
                      {{ link.url }}
                    </a>
                  </td>

                  <td class="p-4 align-middle">
                    <span class="inline-flex items-center rounded-md border border-border px-2 py-0.5 text-[10px] font-medium text-muted-foreground">
                      {{ link.category || '默认' }}
                    </span>
                  </td>

                  <td class="p-4 align-middle font-mono">
                    <span v-if="pingResults[link.id]" class="inline-flex items-center space-x-1.5">
                      <span class="w-1.5 h-1.5 rounded-full" :class="pingResults[link.id].healthy ? 'bg-emerald-500' : 'bg-rose-500'"></span>
                      <span :class="pingResults[link.id].healthy ? 'text-emerald-600 dark:text-emerald-400' : 'text-destructive'">
                        {{ pingResults[link.id].healthy ? `${pingResults[link.id].latency_ms}ms (HTTP ${pingResults[link.id].status_code})` : '超时不可达' }}
                      </span>
                    </span>
                    <span v-else class="text-muted-foreground/60">未检测</span>
                  </td>

                  <!-- ✏️ 操作区：编辑 + 删除 -->
                  <td class="p-4 align-middle text-right space-x-2">
                    <button 
                      @click="openEditLinkModal(link)"
                      class="text-xs font-medium text-foreground hover:underline cursor-pointer"
                    >
                      编辑
                    </button>
                    <button 
                      @click="handleDeleteLink(link.id)"
                      class="text-xs font-medium text-destructive hover:underline cursor-pointer"
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

      <!-- 2. 系统公告通知管理 -->
      <div v-if="activeTab === 'announcements'" class="space-y-6">
        <!-- 发布公告表单 -->
        <div class="rounded-lg border border-border bg-card p-5 space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-foreground">发布系统公告</h3>
            <span class="text-xs text-muted-foreground">支持纯文本公告或导入长篇 Markdown 文档</span>
          </div>

          <form @submit.prevent="handleAddAnnouncement" class="space-y-4">
            <div class="space-y-1">
              <label class="text-xs font-medium text-muted-foreground">公告标题 / 摘要</label>
              <input 
                v-model="newAnnouncement.content" 
                placeholder="例如: 🔥 AI 编程助手接入指南已上线（支持 Claude Code / Codex / WorkBuddy）" 
                class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                required
              />
            </div>

            <!-- Markdown 详情与导入 -->
            <div class="space-y-1.5">
              <div class="flex items-center justify-between">
                <label class="text-xs font-medium text-muted-foreground">详细内容 (Markdown / 可选)</label>
                <div class="flex items-center space-x-2">
                  <input 
                    ref="mdFileInputRef"
                    type="file" 
                    accept=".md,.markdown,.txt" 
                    class="hidden" 
                    @change="(e) => handleImportMd(e, false)"
                  />
                  <button
                    type="button"
                    @click="mdFileInputRef?.click()"
                    class="text-xs text-muted-foreground hover:text-foreground font-medium transition-colors cursor-pointer flex items-center space-x-1 border border-input bg-background hover:bg-accent px-2 py-0.5 rounded shadow-xs"
                  >
                    <svg class="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                    </svg>
                    <span>导入 .md 文档</span>
                  </button>
                </div>
              </div>
              <textarea 
                v-model="newAnnouncement.detail_md" 
                placeholder="在此粘贴 Markdown 内容，或点击右上角导入本地 .md 文件。内容将动态展示在前台内容页中..."
                rows="6"
                class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-xs font-mono shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-y leading-relaxed"
              ></textarea>
            </div>

            <div class="flex items-center justify-between pt-1">
              <label class="flex items-center space-x-2 text-xs text-muted-foreground cursor-pointer">
                <input 
                  v-model="newAnnouncement.is_active" 
                  type="checkbox" 
                  class="rounded border-input text-foreground focus:ring-ring"
                />
                <span>立即在前台首页生效展示</span>
              </label>

              <button 
                type="submit" 
                class="inline-flex items-center justify-center rounded-md text-xs sm:text-sm font-medium transition-colors h-9 px-4 bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 cursor-pointer"
              >
                + 发布公告
              </button>
            </div>
          </form>
        </div>

        <!-- 公告列表 (Table 规范) -->
        <div class="rounded-lg border border-border bg-card overflow-hidden">
          <div class="px-4 py-3 border-b border-border bg-muted/40">
            <span class="text-xs font-semibold text-foreground">公告记录 ({{ announcements.length }})</span>
          </div>

          <div class="divide-y divide-border text-xs">
            <div 
              v-for="item in announcements" 
              :key="item.id" 
              class="p-4 flex items-center justify-between gap-4 hover:bg-muted/40 transition-colors"
            >
              <div class="space-y-1">
                <div class="flex items-center space-x-2">
                  <span 
                    :class="[
                      'w-2 h-2 rounded-full',
                      item.is_active ? 'bg-emerald-500' : 'bg-zinc-400'
                    ]"
                  ></span>
                  <span class="text-[11px] font-mono text-muted-foreground">
                    {{ item.is_active ? '正在生效' : '已停用' }}
                  </span>
                </div>
                <p class="text-sm text-foreground font-medium">{{ item.content }}</p>
              </div>

              <!-- ✏️ 操作：修改 + 启停 + 删除 -->
              <div class="flex items-center space-x-3 shrink-0">
                <button 
                  @click="openEditAnnModal(item)"
                  class="text-xs font-medium text-foreground hover:underline cursor-pointer"
                >
                  编辑
                </button>
                <button 
                  @click="handleToggleAnnouncement(item.id)"
                  class="text-xs border border-input rounded px-2 py-1 hover:bg-accent cursor-pointer"
                >
                  {{ item.is_active ? '停用' : '启用' }}
                </button>
                <button 
                  @click="handleDeleteAnnouncement(item.id)"
                  class="text-xs font-medium text-destructive hover:underline cursor-pointer"
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
        <!-- Chrome 书签导入 (shadcn 规范 Card) -->
        <div class="rounded-lg border border-border bg-card p-5 space-y-3">
          <div class="space-y-1">
            <h3 class="text-sm font-semibold text-foreground">Chrome / Edge 书签导入</h3>
            <p class="text-xs text-muted-foreground">从浏览器书签管理器中导出 HTML 文件，一键批量解析并建立分类索引</p>
          </div>

          <input 
            ref="fileInputRef"
            type="file" 
            accept=".html,.htm" 
            class="hidden" 
            @change="handleBookmarkFileChange"
          />

          <button 
            @click="fileInputRef?.click()"
            class="w-full inline-flex items-center justify-center rounded-md text-xs font-medium border border-input bg-background hover:bg-accent h-9 transition-colors cursor-pointer shadow-xs"
          >
            选择 .html 书签文件并导入
          </button>
        </div>

        <!-- JSON 备份与还原 (shadcn 规范 Card) -->
        <div class="rounded-lg border border-border bg-card p-5 space-y-3">
          <div class="space-y-1">
            <h3 class="text-sm font-semibold text-foreground">全站 JSON 备份与恢复</h3>
            <p class="text-xs text-muted-foreground">导出完整的链接与公告数据集，支持在任意环境一键导入合并</p>
          </div>

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
              class="flex-1 inline-flex items-center justify-center rounded-md text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-9 transition-colors cursor-pointer shadow-xs"
            >
              导出 JSON 备份
            </button>
            <button 
              @click="jsonInputRef?.click()"
              class="flex-1 inline-flex items-center justify-center rounded-md text-xs font-medium border border-input bg-background hover:bg-accent h-9 transition-colors cursor-pointer shadow-xs"
            >
              导入 JSON 恢复
            </button>
          </div>
        </div>
      </div>

      <!-- 4. 系统站点设置 -->
      <div v-if="activeTab === 'settings'" class="max-w-2xl space-y-6">
        <div class="rounded-lg border border-border bg-card p-5 space-y-5">
          <div class="space-y-1">
            <h3 class="text-base font-semibold text-foreground">网站基本信息设置</h3>
            <p class="text-xs text-muted-foreground">自定义前台首页大标题、面包屑导航名称及全站副标题描述</p>
          </div>

          <form @submit.prevent="handleSaveSettings" class="space-y-4">
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-foreground">网站名称 / 首页大标题</label>
              <input 
                v-model="editSettingsForm.site_name" 
                placeholder="例如: Minimal Nav, 研发团队内部工作台" 
                class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                required
              />
              <p class="text-[11px] text-muted-foreground">展示在浏览器标签页、前台首页大标题、内容页面包屑及底栏。</p>
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-medium text-foreground">网站描述 / 副标题 (可选)</label>
              <textarea 
                v-model="editSettingsForm.site_desc" 
                placeholder="例如: 统一汇聚团队核心工具、部署控制台、设计协作及文档中心，即时检索快速直达。" 
                rows="3"
                class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-y leading-relaxed"
              ></textarea>
              <p class="text-[11px] text-muted-foreground">展示在首页大标题下方，留空则不显示。</p>
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-medium text-foreground">ICP 备案号 (可选)</label>
              <input 
                v-model="editSettingsForm.icp_beian" 
                placeholder="例如: 京ICP备XXXXXXXX号-1" 
                class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              />
              <p class="text-[11px] text-muted-foreground">展示在全站底部页脚居中/左侧，点击可直达工信部备案查询系统。留空则不展示。</p>
            </div>

            <div class="pt-2 border-t border-border flex justify-end">
              <button 
                type="submit" 
                :disabled="saveSettingsLoading"
                class="inline-flex items-center justify-center rounded-md text-xs sm:text-sm font-medium transition-colors h-9 px-5 bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 cursor-pointer disabled:opacity-50"
              >
                {{ saveSettingsLoading ? '正在保存...' : '保存系统设置' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 🌟 编辑导航链接 Dialog (纯正 shadcn/ui Dialog 规范) -->
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
          v-if="isEditLinkModalOpen" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
          @click.self="isEditLinkModalOpen = false"
        >
          <div class="w-full max-w-lg rounded-lg border border-border bg-popover text-popover-foreground p-6 shadow-xl space-y-5 animate-in fade-in-0 zoom-in-95 duration-150">
            <div class="space-y-1.5">
              <h3 class="text-base font-semibold tracking-tight text-foreground">编辑导航链接</h3>
              <p class="text-xs text-muted-foreground">修改该链接的标题、目标地址、分类或图标</p>
            </div>

            <form @submit.prevent="handleSaveEditLink" class="space-y-4">
              <div class="space-y-1">
                <label class="text-xs font-medium text-muted-foreground">链接标题</label>
                <input 
                  v-model="editingLink.title" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  required
                />
              </div>

              <div class="space-y-1">
                <div class="flex items-center justify-between">
                  <label class="text-xs font-medium text-muted-foreground">目标网址</label>
                  <button
                    type="button"
                    @click="handleAutoFavicon(true)"
                    class="text-xs text-muted-foreground hover:text-foreground font-medium transition-colors cursor-pointer flex items-center space-x-1"
                  >
                    <svg class="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                    </svg>
                    <span>{{ editFaviconLoading ? '获取中...' : '重新获取图标' }}</span>
                  </button>
                </div>
                <input 
                  v-model="editingLink.url" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  required
                />
              </div>

              <div class="space-y-1">
                <label class="text-xs font-medium text-muted-foreground">所属分类</label>
                <input 
                  v-model="editingLink.category" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                />
              </div>

              <div class="space-y-1">
                <label class="text-xs font-medium text-muted-foreground">Favicon 图标地址</label>
                <div class="flex items-center space-x-2">
                  <input 
                    v-model="editingLink.icon" 
                    placeholder="https://..." 
                    class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-xs font-mono shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  />
                  <div class="w-9 h-9 rounded bg-secondary flex items-center justify-center border border-border shrink-0 overflow-hidden">
                    <img v-if="editingLink.icon" :src="editingLink.icon" class="w-4 h-4 object-contain" />
                    <span v-else class="text-xs font-bold text-foreground">?</span>
                  </div>
                </div>
              </div>

              <div class="flex justify-end space-x-2 pt-2 border-t border-border">
                <button 
                  type="button" 
                  @click="isEditLinkModalOpen = false"
                  class="inline-flex items-center justify-center rounded-md text-xs font-medium border border-input bg-background hover:bg-accent h-9 px-4 transition-colors cursor-pointer"
                >
                  取消
                </button>
                <button 
                  type="submit" 
                  class="inline-flex items-center justify-center rounded-md text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 transition-colors cursor-pointer shadow-xs"
                >
                  保存修改
                </button>
              </div>
            </form>
          </div>
        </div>
      </transition>
    </teleport>

    <!-- 🌟 编辑公告 Dialog (纯正 shadcn/ui Dialog 规范) -->
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
          v-if="isEditAnnModalOpen" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
          @click.self="isEditAnnModalOpen = false"
        >
          <div class="w-full max-w-lg rounded-lg border border-border bg-popover text-popover-foreground p-6 shadow-xl space-y-5 animate-in fade-in-0 zoom-in-95 duration-150">
            <div class="space-y-1.5">
              <h3 class="text-base font-semibold tracking-tight text-foreground">编辑公告内容</h3>
              <p class="text-xs text-muted-foreground">修改公告文本并调整前台展示状态</p>
            </div>

            <form @submit.prevent="handleSaveEditAnnouncement" class="space-y-4">
              <div class="space-y-1">
                <label class="text-xs font-medium text-muted-foreground">公告标题 / 摘要</label>
                <input 
                  v-model="editingAnnouncement.content" 
                  class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  required
                />
              </div>

              <!-- 编辑 Markdown 详情 -->
              <div class="space-y-1.5">
                <div class="flex items-center justify-between">
                  <label class="text-xs font-medium text-muted-foreground">详细内容 (Markdown / 支持文档渲染)</label>
                  <div class="flex items-center space-x-2">
                    <input 
                      ref="editMdFileInputRef"
                      type="file" 
                      accept=".md,.markdown,.txt" 
                      class="hidden" 
                      @change="(e) => handleImportMd(e, true)"
                    />
                    <button
                      type="button"
                      @click="editMdFileInputRef?.click()"
                      class="text-xs text-muted-foreground hover:text-foreground font-medium transition-colors cursor-pointer flex items-center space-x-1 border border-input bg-background hover:bg-accent px-2 py-0.5 rounded shadow-xs"
                    >
                      <span>替换导入 .md</span>
                    </button>
                  </div>
                </div>
                <textarea 
                  v-model="editingAnnouncement.detail_md" 
                  placeholder="在此编辑长篇 Markdown 详情、接入步骤或代码示例..."
                  rows="7"
                  class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-xs font-mono shadow-xs transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring resize-y leading-relaxed"
                ></textarea>
              </div>

              <div class="flex items-center space-x-2 pt-1">
                <input 
                  v-model="editingAnnouncement.is_active" 
                  type="checkbox" 
                  id="edit-active"
                  class="rounded border-input text-foreground focus:ring-ring"
                />
                <label for="edit-active" class="text-xs text-muted-foreground cursor-pointer">
                  在前台首页生效展示
                </label>
              </div>

              <div class="flex justify-end space-x-2 pt-2 border-t border-border">
                <button 
                  type="button" 
                  @click="isEditAnnModalOpen = false"
                  class="inline-flex items-center justify-center rounded-md text-xs font-medium border border-input bg-background hover:bg-accent h-9 px-4 transition-colors cursor-pointer"
                >
                  取消
                </button>
                <button 
                  type="submit" 
                  class="inline-flex items-center justify-center rounded-md text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 transition-colors cursor-pointer shadow-xs"
                >
                  保存修改
                </button>
              </div>
            </form>
          </div>
        </div>
      </transition>
    </teleport>
  </div>
</template>
