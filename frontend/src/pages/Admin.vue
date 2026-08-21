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

const links = ref<LinkItem[]>([])
const announcements = ref<AnnouncementItem[]>([])

const newLink = ref({ title: '', url: '', category: '开发协作' })
const newAnnouncement = ref({ content: '', is_active: true })
const activeTab = ref<'links' | 'announcements'>('links')
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)

const showMessage = (text: string, type: 'success' | 'error' = 'success') => {
  message.value = { type, text }
  setTimeout(() => {
    message.value = null
  }, 3000)
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
  } catch (err) {
    showMessage('获取数据失败', 'error')
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
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newLink.value),
    })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('链接添加成功')
      newLink.value = { title: '', url: '', category: '开发协作' }
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
    const res = await fetch(`/api/links/${id}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('链接已删除')
      loadData()
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
      headers: { 'Content-Type': 'application/json' },
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
    const res = await fetch(`/api/announcements/${id}/toggle`, { method: 'PUT' })
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
    const res = await fetch(`/api/announcements/${id}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.code === 0) {
      showMessage('公告已删除')
      loadData()
    }
  } catch {
    showMessage('删除失败', 'error')
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="space-y-10 max-w-5xl mx-auto">
    <!-- 顶部标题 -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border/50 pb-6">
      <div>
        <h2 class="text-2xl sm:text-3xl font-bold tracking-tight text-foreground">后台资源管理</h2>
        <p class="text-sm text-muted-foreground mt-1">配置团队导航链接与全站通知公告</p>
      </div>

      <!-- 选项卡切换 -->
      <div class="flex items-center space-x-2 bg-secondary/60 p-1 rounded-lg border border-border/50 self-start sm:self-auto">
        <button
          @click="activeTab = 'links'"
          :class="[
            'text-sm px-4 py-2 rounded-md font-medium transition-all cursor-pointer',
            activeTab === 'links' 
              ? 'bg-card text-foreground shadow-sm' 
              : 'text-muted-foreground hover:text-foreground'
          ]"
        >
          导航链接 ({{ links.length }})
        </button>
        <button
          @click="activeTab = 'announcements'"
          :class="[
            'text-sm px-4 py-2 rounded-md font-medium transition-all cursor-pointer',
            activeTab === 'announcements' 
              ? 'bg-card text-foreground shadow-sm' 
              : 'text-muted-foreground hover:text-foreground'
          ]"
        >
          系统公告 ({{ announcements.length }})
        </button>
      </div>
    </div>

    <!-- 全局提示条 -->
    <div 
      v-if="message" 
      :class="[
        'px-5 py-3 text-sm rounded-lg border transition-all duration-200 shadow-sm flex items-center space-x-2',
        message.type === 'success' 
          ? 'bg-zinc-100 dark:bg-zinc-900 border-zinc-300 dark:border-zinc-700 text-foreground'
          : 'bg-red-500/10 border-red-500/30 text-red-500'
      ]"
    >
      <span>{{ message.text }}</span>
    </div>

    <!-- 链接管理面板 -->
    <div v-if="activeTab === 'links'" class="space-y-8">
      <!-- 新增链接表单 -->
      <form @submit.prevent="handleAddLink" class="p-6 rounded-xl border border-border/80 bg-card shadow-sm space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-base font-semibold text-foreground tracking-tight">添加新导航链接</h3>
          <span class="text-xs text-muted-foreground">将自动在公共导航页同步展示</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">链接名称</label>
            <input 
              v-model="newLink.title" 
              placeholder="例如: Figma, Linear" 
              class="w-full h-11 bg-background border border-border text-sm rounded-lg px-3.5 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
              required
            />
          </div>
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">目标网址</label>
            <input 
              v-model="newLink.url" 
              placeholder="https://example.com" 
              class="w-full h-11 bg-background border border-border text-sm rounded-lg px-3.5 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
              required
            />
          </div>
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-muted-foreground">所属分类</label>
            <input 
              v-model="newLink.category" 
              placeholder="例如: 开发协作, 运维部署" 
              class="w-full h-11 bg-background border border-border text-sm rounded-lg px-3.5 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
            />
          </div>
        </div>

        <div class="flex justify-end pt-2">
          <button 
            type="submit" 
            class="bg-foreground text-background hover:opacity-90 text-sm px-6 py-2.5 rounded-lg font-medium transition-opacity cursor-pointer shadow-sm"
          >
            保存并添加
          </button>
        </div>
      </form>

      <!-- 链接表格 -->
      <div class="border border-border/80 rounded-xl overflow-hidden bg-card shadow-sm">
        <table class="w-full text-left text-sm">
          <thead class="bg-secondary/70 text-muted-foreground border-b border-border text-xs uppercase tracking-wider">
            <tr>
              <th class="py-3.5 px-6 font-semibold">名称</th>
              <th class="py-3.5 px-6 font-semibold">跳转地址</th>
              <th class="py-3.5 px-6 font-semibold">分类</th>
              <th class="py-3.5 px-6 font-semibold text-right">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border/60">
            <tr v-for="link in links" :key="link.id" class="hover:bg-muted/30 transition-colors">
              <td class="py-4 px-6 font-medium text-foreground">{{ link.title }}</td>
              <td class="py-4 px-6 text-muted-foreground font-mono text-xs truncate max-w-[280px]">{{ link.url }}</td>
              <td class="py-4 px-6">
                <span class="inline-block px-2.5 py-1 rounded text-xs font-medium bg-secondary text-foreground">
                  {{ link.category || '默认' }}
                </span>
              </td>
              <td class="py-4 px-6 text-right">
                <button 
                  @click="handleDeleteLink(link.id)" 
                  class="text-sm font-medium text-red-500 hover:text-red-600 transition-colors cursor-pointer"
                >
                  删除
                </button>
              </td>
            </tr>
            <tr v-if="links.length === 0">
              <td colspan="4" class="py-12 text-center text-muted-foreground text-sm">暂无链接数据，请在上方添加</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 公告管理面板 -->
    <div v-else-if="activeTab === 'announcements'" class="space-y-8">
      <!-- 发布新公告表单 -->
      <form @submit.prevent="handleAddAnnouncement" class="p-6 rounded-xl border border-border/80 bg-card shadow-sm space-y-5">
        <h3 class="text-base font-semibold text-foreground tracking-tight">发布全站公告通知</h3>
        <textarea 
          v-model="newAnnouncement.content" 
          rows="3" 
          placeholder="请输入公告内容，如系统维护时间、团队重要通知等..." 
          class="w-full bg-background border border-border text-sm rounded-lg p-3.5 focus:outline-none focus:ring-2 focus:ring-foreground/20 focus:border-foreground"
          required
        ></textarea>
        <div class="flex items-center justify-between">
          <label class="flex items-center space-x-2 text-sm text-muted-foreground cursor-pointer select-none">
            <input type="checkbox" v-model="newAnnouncement.is_active" class="rounded border-border text-foreground w-4 h-4" />
            <span>立即对所有用户生效</span>
          </label>
          <button 
            type="submit" 
            class="bg-foreground text-background hover:opacity-90 text-sm px-6 py-2.5 rounded-lg font-medium transition-opacity cursor-pointer shadow-sm"
          >
            发布公告
          </button>
        </div>
      </form>

      <!-- 公告列表 -->
      <div class="border border-border/80 rounded-xl overflow-hidden bg-card shadow-sm">
        <table class="w-full text-left text-sm">
          <thead class="bg-secondary/70 text-muted-foreground border-b border-border text-xs uppercase tracking-wider">
            <tr>
              <th class="py-3.5 px-6 font-semibold">公告内容</th>
              <th class="py-3.5 px-6 font-semibold w-28">状态</th>
              <th class="py-3.5 px-6 font-semibold text-right w-44">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border/60">
            <tr v-for="item in announcements" :key="item.id" class="hover:bg-muted/30 transition-colors">
              <td class="py-4 px-6 text-foreground leading-relaxed">{{ item.content }}</td>
              <td class="py-4 px-6">
                <span 
                  :class="[
                    'inline-block px-2.5 py-1 rounded-full text-xs font-semibold',
                    item.is_active 
                      ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20' 
                      : 'bg-muted text-muted-foreground'
                  ]"
                >
                  {{ item.is_active ? '● 生效中' : '已停用' }}
                </span>
              </td>
              <td class="py-4 px-6 text-right space-x-3">
                <button 
                  @click="handleToggleAnnouncement(item.id)" 
                  class="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                >
                  {{ item.is_active ? '停用' : '启用' }}
                </button>
                <button 
                  @click="handleDeleteAnnouncement(item.id)" 
                  class="text-sm font-medium text-red-500 hover:text-red-600 transition-colors cursor-pointer"
                >
                  删除
                </button>
              </td>
            </tr>
            <tr v-if="announcements.length === 0">
              <td colspan="3" class="py-12 text-center text-muted-foreground text-sm">暂无公告数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
