<script setup lang="ts">
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import { useSiteConfig } from '../utils/useSiteConfig'

const route = useRoute()
const { siteConfig, loadSiteConfig } = useSiteConfig()
const loading = ref(true)
const announcement = ref<{
  id?: number
  content?: string
  detail_md?: string
  created_at?: string
} | null>(null)

// 内置 AI 指南 Tab 状态 (当无自定义 detail_md 时降级展示)
const activeAiTab = ref<'claude' | 'codex' | 'workbuddy'>('claude')
const copiedIndex = ref<string | null>(null)

// 复制文本
const copyCode = (text: string, id: string) => {
  navigator.clipboard.writeText(text).then(() => {
    copiedIndex.value = id
    setTimeout(() => {
      if (copiedIndex.value === id) copiedIndex.value = null
    }, 2000)
  })
}

// 获取详情数据
const fetchDetail = async () => {
  loading.value = true
  const id = route.params.id

  if (id) {
    try {
      const res = await fetch(`/api/announcements/${id}`)
      if (res.ok) {
        const data = await res.json()
        if (data.code === 0 && data.data) {
          announcement.value = data.data
        }
      }
    } catch {
      // ignore
    }
  } else {
    // 若为 /docs/ai 或未传参，尝试获取最新的公告详情
    try {
      const res = await fetch('/api/announcements/active')
      if (res.ok) {
        const data = await res.json()
        if (data.code === 0 && data.data && data.data.length > 0) {
          // 优先找带有 detail_md 的或 AI 相关的公告
          const target = data.data.find((item: any) => item.detail_md || item.content?.includes('AI')) || data.data[0]
          announcement.value = target
        }
      }
    } catch {
      // ignore
    }
  }
  loading.value = false
}

// 解析 Markdown 为 HTML (自动消除首行重复的一级标题，避免双重标题与双重下划线)
const renderedHtml = computed(() => {
  if (!announcement.value?.detail_md) return ''
  let md = announcement.value.detail_md.trim()
  
  // 检查首行是否包含一级标题
  const firstLineMatch = md.match(/^#\s+(.+)$/m)
  if (firstLineMatch) {
    const mdTitle = firstLineMatch[1].trim()
    const contentTitle = (announcement.value.content || '').trim()
    // 若 Markdown 首行标题与公告标题含义相同，剥离首行避免重复渲染
    if (contentTitle && (
      contentTitle.includes(mdTitle) || 
      mdTitle.includes(contentTitle) ||
      contentTitle.replace(/^[^\w\u4e00-\u9fa5]+/, '').trim() === mdTitle.replace(/^[^\w\u4e00-\u9fa5]+/, '').trim()
    )) {
      md = md.replace(/^#\s+.+(\r?\n|$)/, '').trim()
    }
  }

  return marked.parse(md) as string
})

// 为动态 Markdown 中的所有代码块增强一键复制按钮
const enhanceCodeBlocks = () => {
  nextTick(() => {
    const codeContainers = document.querySelectorAll('.markdown-body pre')
    codeContainers.forEach((pre) => {
      if (pre.querySelector('.copy-btn-dynamic')) return

      pre.classList.add('relative', 'group')
      const btn = document.createElement('button')
      btn.className = 'copy-btn-dynamic absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer opacity-80 hover:opacity-100'
      btn.innerText = '复制'
      btn.onclick = () => {
        const codeText = (pre.querySelector('code') as HTMLElement)?.innerText || (pre as HTMLElement).innerText || ''
        navigator.clipboard.writeText(codeText).then(() => {
          btn.innerText = '✓ 已复制'
          btn.classList.add('text-emerald-600', 'dark:text-emerald-400')
          setTimeout(() => {
            btn.innerText = '复制'
            btn.classList.remove('text-emerald-600', 'dark:text-emerald-400')
          }, 1800)
        })
      }
      pre.appendChild(btn)
    })
  })
}

watch(renderedHtml, () => {
  if (renderedHtml.value) {
    enhanceCodeBlocks()
  }
})

onMounted(() => {
  loadSiteConfig()
  fetchDetail().then(() => {
    if (renderedHtml.value) {
      enhanceCodeBlocks()
    }
  })
})
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-7 pb-16 antialiased">
    <!-- 面包屑导航 (极简无衬线，字迹饱满无毛刺) -->
    <nav class="flex items-center space-x-2 text-xs text-muted-foreground select-none">
      <router-link to="/" class="hover:text-foreground hover:underline transition-colors">
        {{ siteConfig.site_name }}
      </router-link>
      <span class="text-muted-foreground/40 font-light">/</span>
      <span class="text-foreground font-medium">公告详情</span>
    </nav>

    <!-- 加载中骨架 -->
    <div v-if="loading" class="space-y-4 animate-pulse">
      <div class="h-9 w-2/3 bg-muted rounded"></div>
      <div class="h-4 w-1/3 bg-muted rounded"></div>
      <div class="h-64 w-full bg-muted/40 rounded-lg"></div>
    </div>

    <!-- 1. 动态渲染后台配置的 Markdown 内容 -->
    <div v-else-if="announcement?.detail_md" class="space-y-6">
      <!-- 统一出版物级大标题 (单一纯粹，避免多重重复下划线) -->
      <div class="space-y-2.5 border-b border-border/70 pb-4">
        <h1 class="text-2xl sm:text-3xl font-semibold tracking-tight text-foreground leading-snug">
          {{ announcement.content }}
        </h1>
        <div class="flex items-center space-x-3 text-xs text-muted-foreground">
          <span class="inline-flex items-center px-1.5 py-0.5 rounded border border-border/80 text-muted-foreground text-[10px] font-medium">
            公告详情
          </span>
          <span v-if="announcement.created_at" class="font-mono text-[11px]">
            {{ new Date(announcement.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }) }}
          </span>
        </div>
      </div>

      <!-- 纯正出版物风格 Markdown 渲染容器 -->
      <div 
        class="markdown-body text-foreground text-sm leading-relaxed space-y-4"
        v-html="renderedHtml"
      ></div>
    </div>

    <!-- 2. 若无后台自定义 Markdown，则降级展示内置标准 AI 接入指南 -->
    <div v-else class="space-y-8 animate-in fade-in-0 duration-150">
      <!-- 页面大标题与导语 -->
      <div class="space-y-2 border-b border-border/70 pb-4">
        <h1 class="text-2xl sm:text-3xl font-semibold tracking-tight text-foreground leading-snug">
          {{ announcement?.content || 'AI 编程助手接入指南' }}
        </h1>
        <p class="text-xs sm:text-sm text-muted-foreground leading-relaxed">
          团队统一接入 724AI 镜像中转网络，为 Claude Code、Codex 及 WorkBuddy 客户端提供稳定免翻墙的极速编码服务。
        </p>
      </div>

      <!-- 🔥 准备工作与注意事项提示框 -->
      <div class="p-4 rounded-lg border border-blue-500/30 bg-blue-500/10 text-xs sm:text-sm text-foreground space-y-1.5 leading-relaxed">
        <div class="font-semibold text-blue-600 dark:text-blue-400 flex items-center space-x-1.5">
          <span>🔥 准备工作与注意事项：</span>
        </div>
        <p>
          1. <strong>独立专属 Key：</strong>每个人使用独立的专属 API Key，请勿混用或共享。开通权限或获取 Key 请在飞书联系 <strong>@闫东</strong>。
        </p>
        <p>
          2. <strong>分组说明：</strong>当前 Key 区分 <strong>Claude</strong> 与 <strong>Codex</strong> 独立路由分组，如需跨组使用请联系管理员调整。
        </p>
      </div>

      <!-- 🌟 直线下划线 Tab 切换栏 -->
      <div class="flex items-center space-x-6 border-b border-border text-sm">
        <button
          @click="activeAiTab = 'claude'"
          :class="[
            'pb-3 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeAiTab === 'claude'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>Claude Code</span>
          <span class="text-[10px] font-mono text-muted-foreground">CLI / 终端</span>
        </button>

        <button
          @click="activeAiTab = 'codex'"
          :class="[
            'pb-3 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeAiTab === 'codex'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>Codex</span>
          <span class="text-[10px] font-mono text-muted-foreground">CLI & App</span>
        </button>

        <button
          @click="activeAiTab = 'workbuddy'"
          :class="[
            'pb-3 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
            activeAiTab === 'workbuddy'
              ? 'border-foreground text-foreground font-semibold'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          ]"
        >
          <span>WorkBuddy</span>
          <span class="text-[10px] font-mono text-muted-foreground">桌面客户端</span>
        </button>
      </div>

      <!-- Tab 1: Claude Code -->
      <div v-if="activeAiTab === 'claude'" class="space-y-8 animate-in fade-in-0 duration-150">
        <p class="text-sm text-muted-foreground">
          将 Claude Code 的默认请求地址替换为 <strong>724AI</strong> 中转站（<code class="font-mono bg-muted px-1.5 py-0.5 rounded text-foreground">api.724ai.org</code>），直连稳定无需代理。
        </p>

        <!-- 1. 安装 Claude Code -->
        <section class="space-y-4">
          <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
            1. 安装 Claude Code
          </h2>
          <p class="text-xs sm:text-sm text-muted-foreground">
            确保已安装 Node.js 环境（建议 Node 18+）。根据网络环境选择一种安装方式：
          </p>

          <div class="space-y-2">
            <h3 class="text-sm font-medium text-foreground">方式一：npm 安装（推荐，无需代理）</h3>
            <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
              <button
                @click="copyCode('npm install -g @anthropic-ai/claude-code', 'claude-npm')"
                class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
              >
                {{ copiedIndex === 'claude-npm' ? '✓ 已复制' : '复制' }}
              </button>
              <code>npm install -g @anthropic-ai/claude-code</code>
            </div>
            <p class="text-xs text-muted-foreground">
              安装完成后，即可在终端使用 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">claude</code> 命令。VS Code 用户可在插件市场搜索安装「Claude Code for VS Code」。
            </p>
          </div>

          <div class="space-y-2 pt-2">
            <h3 class="text-sm font-medium text-foreground">方式二：官方脚本安装（需代理）</h3>
            <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
              <button
                @click="copyCode('curl -fsSL https://claude.ai/install.sh | bash', 'claude-curl')"
                class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
              >
                {{ copiedIndex === 'claude-curl' ? '✓ 已复制' : '复制' }}
              </button>
              <code>curl -fsSL https://claude.ai/install.sh | bash</code>
            </div>
            <div class="p-3 rounded-md border border-amber-500/30 bg-amber-500/10 text-xs text-amber-800 dark:text-amber-300">
              若无代理直连会返回网页错误页，导致 bash 报 <code class="font-mono bg-amber-500/20 px-1 rounded">syntax error</code> 乱码。遇到报错请改用方式一（npm 安装）。
            </div>
          </div>
        </section>

        <!-- 2. 配置 724AI 接入 -->
        <section class="space-y-4">
          <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
            2. 配置 724AI 接入
          </h2>
          <p class="text-xs sm:text-sm text-muted-foreground">
            编辑配置文件 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">~/.claude/settings.json</code>（Windows 对应 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">C:\Users\用户名\.claude\settings.json</code>）：
          </p>

          <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
            <button
              @click="copyCode(`{\n  &quot;env&quot;: {\n    &quot;ANTHROPIC_API_KEY&quot;: &quot;sk-xxxxxxxxxxxxxxxx&quot;,\n    &quot;ANTHROPIC_BASE_URL&quot;: &quot;https://api.724ai.org&quot;,\n    &quot;CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS&quot;: &quot;1&quot;\n  }\n}`, 'claude-config')"
              class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
            >
              {{ copiedIndex === 'claude-config' ? '✓ 已复制' : '复制' }}
            </button>
            <pre class="leading-relaxed"><code>{
  "env": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.724ai.org",
    "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS": "1"
  }
}</code></pre>
          </div>

          <div class="p-3 rounded-md border border-rose-500/30 bg-rose-500/10 text-xs text-rose-700 dark:text-rose-300 leading-relaxed">
            <strong>注意：</strong><code class="font-mono bg-rose-500/20 px-1 rounded">ANTHROPIC_BASE_URL</code> 必须填写 <code class="font-mono bg-rose-500/20 px-1 rounded">https://api.724ai.org</code>，末尾<strong>严禁加斜杠 <code>/</code></strong>。请将 API_KEY 替换为您自己的专属 Key。
          </div>
        </section>

        <!-- 3. 免登录配置 -->
        <section class="space-y-4">
          <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
            3. 免登录配置（跳过官方鉴权）
          </h2>
          <p class="text-xs sm:text-sm text-muted-foreground">
            为了跳过强制的官方网页登录，在 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">~/.claude.json</code> 中添加跳过标识：
          </p>

          <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
            <button
              @click="copyCode(`{\n  &quot;hasCompletedOnboarding&quot;: true\n}`, 'claude-onboard')"
              class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
            >
              {{ copiedIndex === 'claude-onboard' ? '✓ 已复制' : '复制' }}
            </button>
            <pre class="leading-relaxed"><code>{
  "hasCompletedOnboarding": true
}</code></pre>
          </div>
        </section>

        <!-- 4. 开始使用与模型推荐 -->
        <section class="space-y-4">
          <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
            4. 开始使用与模型推荐
          </h2>
          <p class="text-xs sm:text-sm text-muted-foreground">
            在项目目录下输入命令启动：
          </p>

          <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
            <button
              @click="copyCode('claude', 'claude-run')"
              class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
            >
              {{ copiedIndex === 'claude-run' ? '✓ 已复制' : '复制' }}
            </button>
            <code>claude</code>
          </div>

          <div class="p-3 rounded-md border border-amber-500/30 bg-amber-500/10 text-xs text-amber-800 dark:text-amber-300">
            启动后 Claude 若提示是否使用自定义 Key，必须输入 <strong>yes</strong>。如果选错，可删除 <code class="font-mono bg-amber-500/20 px-1 rounded">~/.claude.json</code> 重新配置。
          </div>
        </section>
      </div>

      <!-- Tab 2: Codex -->
      <div v-if="activeAiTab === 'codex'" class="space-y-8 animate-in fade-in-0 duration-150">
        <p class="text-sm text-muted-foreground">
          Codex 提供桌面 App（图形界面，新手首选）和 CLI 命令行两种使用方式，共享同一套 724AI 中转配置。
        </p>

        <section class="space-y-4">
          <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
            1. 安装 Codex
          </h2>
          <div class="space-y-2">
            <h3 class="text-sm font-medium text-foreground">桌面客户端 App（图形界面，新手首选）</h3>
            <p class="text-xs sm:text-sm text-muted-foreground">
              前往 <a href="https://openai.com/zh-Hans-CN/codex/" target="_blank" class="text-primary underline">Codex 官方页面</a> 下载对应平台版本。
            </p>
          </div>

          <div class="space-y-2 pt-2">
            <h3 class="text-sm font-medium text-foreground">CLI 命令行安装</h3>
            <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
              <button
                @click="copyCode('npm install -g @openai/codex', 'codex-cli')"
                class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
              >
                {{ copiedIndex === 'codex-cli' ? '✓ 已复制' : '复制' }}
              </button>
              <code>npm install -g @openai/codex</code>
            </div>
          </div>
        </section>
      </div>

      <!-- Tab 3: WorkBuddy -->
      <div v-if="activeAiTab === 'workbuddy'" class="space-y-8 animate-in fade-in-0 duration-150">
        <p class="text-sm text-muted-foreground">
          WorkBuddy 支持兼容 OpenAI 协议的自定义模型，直接使用 724AI 通用接口地址即可。
        </p>
        <section class="space-y-4">
          <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
            配置填写
          </h2>
          <div class="p-4 rounded-lg border border-border bg-card text-xs sm:text-sm space-y-2 font-mono">
            <div><strong>接口地址:</strong> <code>https://api.724ai.org/v1</code></div>
            <div><strong>模型名称:</strong> <code>gpt-5.6-sol</code> / <code>gpt-5.6-luna</code></div>
          </div>
        </section>
      </div>
    </div>

    <!-- 底部返回与团队标识 -->
    <div class="pt-8 border-t border-border flex items-center justify-between text-xs text-muted-foreground">
      <span>© {{ new Date().getFullYear() }} {{ siteConfig.site_name }}</span>
      <router-link to="/" class="hover:text-foreground hover:underline transition-colors">
        ← 返回首页导航
      </router-link>
    </div>
  </div>
</template>

<style>
/* 针对 Markdown 解析内容的极简高质感排版 (纯粹无毛刺抗锯齿) */
.markdown-body {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-rendering: optimizeLegibility;
  color: var(--foreground);
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4 {
  font-weight: 600;
  color: var(--foreground);
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  letter-spacing: -0.015em;
  line-height: 1.35;
}
.markdown-body h1 { font-size: 1.35rem; border-bottom: 1px solid var(--border); padding-bottom: 0.3em; }
.markdown-body h2 { font-size: 1.15rem; border-bottom: 1px solid var(--border); padding-bottom: 0.3em; }
.markdown-body h3 { font-size: 1rem; }

.markdown-body p {
  margin-top: 0.5em;
  margin-bottom: 0.8em;
  color: var(--muted-foreground);
  line-height: 1.7;
}

.markdown-body pre {
  background-color: var(--muted);
  border: 1px solid var(--border);
  border-radius: 0.375rem;
  padding: 0.875rem;
  overflow-x: auto;
  margin: 1em 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.8125rem;
}

.markdown-body code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.8125rem;
  padding: 0.15em 0.35em;
  border-radius: 0.25rem;
  background-color: var(--muted);
  color: var(--foreground);
}

.markdown-body pre code {
  padding: 0;
  background: transparent;
}

.markdown-body ul,
.markdown-body ol {
  padding-left: 1.5em;
  margin: 0.5em 0 1em;
  color: var(--muted-foreground);
}

.markdown-body li {
  margin-bottom: 0.3em;
  line-height: 1.6;
}

.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.2em 0;
  font-size: 0.85rem;
}

.markdown-body th,
.markdown-body td {
  border: 1px solid var(--border);
  padding: 0.5rem 0.75rem;
  text-align: left;
}

.markdown-body th {
  background-color: var(--muted);
  font-weight: 600;
  color: var(--foreground);
}

.markdown-body blockquote {
  border-left: 3px solid var(--border);
  padding-left: 1rem;
  margin: 1em 0;
  color: var(--muted-foreground);
}

.markdown-body a {
  color: var(--primary);
  text-decoration: underline;
}
</style>
