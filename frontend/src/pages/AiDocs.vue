<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref<'claude' | 'codex' | 'workbuddy'>('claude')
const copiedIndex = ref<string | null>(null)

// 复制文本到剪贴板
const copyCode = (text: string, id: string) => {
  navigator.clipboard.writeText(text).then(() => {
    copiedIndex.value = id
    setTimeout(() => {
      if (copiedIndex.value === id) {
        copiedIndex.value = null
      }
    }, 2000)
  })
}
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-8 pb-16">
    <!-- 面包屑导航 -->
    <div class="flex items-center space-x-2 text-xs font-mono text-muted-foreground">
      <router-link to="/" class="hover:text-foreground hover:underline transition-colors">
        内部工作台
      </router-link>
      <span>/</span>
      <span class="text-foreground">AI 编程助手接入指南</span>
    </div>

    <!-- 页面大标题与导语 -->
    <div class="space-y-2">
      <h1 class="text-3xl sm:text-4xl font-bold tracking-tight text-foreground">
        AI 编程助手接入指南
      </h1>
      <p class="text-sm text-muted-foreground leading-relaxed">
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

    <!-- 🌟 直线下划线 Tab 切换栏 (无多余卡片包裹) -->
    <div class="flex items-center space-x-6 border-b border-border text-sm">
      <button
        @click="activeTab = 'claude'"
        :class="[
          'pb-3 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
          activeTab === 'claude'
            ? 'border-foreground text-foreground font-semibold'
            : 'border-transparent text-muted-foreground hover:text-foreground'
        ]"
      >
        <span>Claude Code</span>
        <span class="text-[10px] font-mono text-muted-foreground">CLI / 终端</span>
      </button>

      <button
        @click="activeTab = 'codex'"
        :class="[
          'pb-3 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
          activeTab === 'codex'
            ? 'border-foreground text-foreground font-semibold'
            : 'border-transparent text-muted-foreground hover:text-foreground'
        ]"
      >
        <span>Codex</span>
        <span class="text-[10px] font-mono text-muted-foreground">CLI & App</span>
      </button>

      <button
        @click="activeTab = 'workbuddy'"
        :class="[
          'pb-3 font-medium transition-all cursor-pointer border-b-2 -mb-px flex items-center space-x-2 text-xs sm:text-sm',
          activeTab === 'workbuddy'
            ? 'border-foreground text-foreground font-semibold'
            : 'border-transparent text-muted-foreground hover:text-foreground'
        ]"
      >
        <span>WorkBuddy</span>
        <span class="text-[10px] font-mono text-muted-foreground">桌面客户端</span>
      </button>
    </div>

    <!-- ============================================== -->
    <!-- Tab 1: Claude Code                             -->
    <!-- ============================================== -->
    <div v-if="activeTab === 'claude'" class="space-y-8 animate-in fade-in-0 duration-150">
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

        <div class="space-y-2 pt-2">
          <h3 class="text-sm font-medium text-foreground">常用模型推荐</h3>
          <ul class="list-disc list-inside text-xs sm:text-sm text-muted-foreground space-y-1">
            <li><code class="font-mono text-foreground font-semibold">claude-sonnet-5</code>：日常主力编码推荐（响应快、成本适中）</li>
            <li><code class="font-mono text-foreground font-semibold">claude-fable-5</code>：复杂攻坚任务模型</li>
          </ul>

          <p class="text-xs text-muted-foreground pt-1">指定模型单行启动：</p>
          <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
            <button
              @click="copyCode('claude --model claude-sonnet-5', 'claude-model-cmd')"
              class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
            >
              {{ copiedIndex === 'claude-model-cmd' ? '✓ 已复制' : '复制' }}
            </button>
            <code>claude --model claude-sonnet-5</code>
          </div>
        </div>
      </section>
    </div>

    <!-- ============================================== -->
    <!-- Tab 2: Codex                                   -->
    <!-- ============================================== -->
    <div v-if="activeTab === 'codex'" class="space-y-8 animate-in fade-in-0 duration-150">
      <p class="text-sm text-muted-foreground">
        Codex 提供桌面 App（图形界面，新手首选）和 CLI 命令行两种使用方式，共享同一套 724AI 中转配置。
      </p>

      <!-- 1. 安装 Codex -->
      <section class="space-y-4">
        <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
          1. 安装 Codex
        </h2>

        <div class="space-y-2">
          <h3 class="text-sm font-medium text-foreground">方式 A：桌面客户端 App（图形界面，新手首选）</h3>
          <p class="text-xs sm:text-sm text-muted-foreground">
            前往 <a href="https://openai.com/zh-Hans-CN/codex/" target="_blank" class="text-primary underline">Codex 官方页面</a> 下载对应平台版本：
          </p>
          <ul class="list-disc list-inside text-xs sm:text-sm text-muted-foreground space-y-1 pl-1">
            <li><strong>Windows：</strong>点击下载入口，跳转至 Microsoft Store 安装。</li>
            <li><strong>macOS（Apple Silicon）：</strong>下载 Apple Silicon 版本（M1/M2/M3/M4）。</li>
            <li><strong>macOS（Intel）：</strong>下载 Intel 版本（Intel Core i5/i7/i9）。</li>
          </ul>
        </div>

        <div class="space-y-3 pt-2">
          <h3 class="text-sm font-medium text-foreground">方式 B：CLI 命令行工具</h3>
          
          <div class="space-y-1">
            <p class="text-xs text-muted-foreground font-medium">macOS / Linux (npm 或 Homebrew 安装)：</p>
            <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
              <button
                @click="copyCode('npm install -g @openai/codex', 'codex-npm')"
                class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
              >
                {{ copiedIndex === 'codex-npm' ? '✓ 已复制' : '复制' }}
              </button>
              <code>npm install -g @openai/codex</code>
            </div>
          </div>

          <div class="space-y-1 pt-1">
            <p class="text-xs text-muted-foreground font-medium">Windows (PowerShell 安装)：</p>
            <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
              <button
                @click="copyCode('npm install -g @openai/codex', 'codex-win-npm')"
                class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
              >
                {{ copiedIndex === 'codex-win-npm' ? '✓ 已复制' : '复制' }}
              </button>
              <code>npm install -g @openai/codex</code>
            </div>
          </div>
        </div>
      </section>

      <!-- 2. 核心配置文件 -->
      <section class="space-y-4">
        <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
          2. 核心配置文件
        </h2>
        <p class="text-xs sm:text-sm text-muted-foreground">
          新建或编辑 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">~/.codex/config.toml</code>（Windows 对应 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">C:\Users\用户名\.codex\config.toml</code>）：
        </p>

        <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
          <button
            @click="copyCode(`model_provider = &quot;724AI&quot;\nmodel_reasoning_effort = &quot;xhigh&quot;\ndisable_response_storage = true\nmodel_context_window = 1000000\nmodel_auto_compact_token_limit = 900000\n\n[model_providers.724AI]\nname = &quot;724AI&quot;\nwire_api = &quot;responses&quot;\nbase_url = &quot;https://api.724ai.org/v1&quot;\nrequires_openai_auth = true\nsupports_websockets = false`, 'codex-config')"
            class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
          >
            {{ copiedIndex === 'codex-config' ? '✓ 已复制' : '复制' }}
          </button>
          <pre class="leading-relaxed"><code>model_provider = "724AI"
model_reasoning_effort = "xhigh"
disable_response_storage = true
model_context_window = 1000000
model_auto_compact_token_limit = 900000

[model_providers.724AI]
name = "724AI"
wire_api = "responses"
base_url = "https://api.724ai.org/v1"
requires_openai_auth = true
supports_websockets = false</code></pre>
        </div>

        <div class="p-3 rounded-md border border-amber-500/30 bg-amber-500/10 text-xs text-amber-800 dark:text-amber-300">
          <strong>注意：</strong><code class="font-mono bg-amber-500/20 px-1 rounded">base_url</code> 必须以 <code class="font-mono bg-amber-500/20 px-1 rounded">/v1</code> 结尾，即 <code class="font-mono bg-amber-500/20 px-1 rounded">https://api.724ai.org/v1</code>。
        </div>

        <p class="text-xs sm:text-sm text-muted-foreground pt-1">
          在同一目录新建 <code class="font-mono bg-muted px-1 py-0.5 rounded text-foreground">~/.codex/auth.json</code>，写入专属 Key：
        </p>

        <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
          <button
            @click="copyCode(`{\n  &quot;OPENAI_API_KEY&quot;: &quot;sk-xxxxxxxxxxxxxxx&quot;\n}`, 'codex-auth')"
            class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
          >
            {{ copiedIndex === 'codex-auth' ? '✓ 已复制' : '复制' }}
          </button>
          <pre class="leading-relaxed"><code>{
  "OPENAI_API_KEY": "sk-xxxxxxxxxxxxxxx"
}</code></pre>
        </div>
      </section>

      <!-- 3. 启动与模型 -->
      <section class="space-y-4">
        <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
          3. 启动与模型说明
        </h2>
        <div class="relative group rounded-md border border-border bg-muted/40 font-mono text-xs p-3.5 overflow-x-auto">
          <button
            @click="copyCode('codex', 'codex-run')"
            class="absolute top-2 right-2 border border-border bg-background hover:bg-accent px-2 py-1 rounded text-[11px] font-sans transition-colors cursor-pointer"
          >
            {{ copiedIndex === 'codex-run' ? '✓ 已复制' : '复制' }}
          </button>
          <code>codex</code>
        </div>

        <div class="space-y-2 pt-1">
          <h3 class="text-sm font-medium text-foreground">常用模型说明</h3>
          <ul class="list-disc list-inside text-xs sm:text-sm text-muted-foreground space-y-1">
            <li><code class="font-mono text-foreground font-semibold">gpt-5.6-sol</code>：复杂攻坚任务模型（推理能力强，消耗较大）</li>
            <li><code class="font-mono text-foreground font-semibold">gpt-5.6-luna</code>：日常主力编码推荐（响应快、成本适中）</li>
          </ul>
        </div>
      </section>
    </div>

    <!-- ============================================== -->
    <!-- Tab 3: WorkBuddy                               -->
    <!-- ============================================== -->
    <div v-if="activeTab === 'workbuddy'" class="space-y-8 animate-in fade-in-0 duration-150">
      <p class="text-sm text-muted-foreground">
        WorkBuddy 支持兼容 OpenAI 协议的自定义模型，直接使用 724AI 通用接口地址即可。
      </p>

      <!-- 1. 下载与安装 -->
      <section class="space-y-4">
        <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
          1. 下载与安装
        </h2>
        <p class="text-xs sm:text-sm text-muted-foreground">
          前往官网获取最新客户端：<a href="https://www.workbuddy.cn/" target="_blank" class="text-primary underline font-mono">https://www.workbuddy.cn/</a>
        </p>
        <p class="text-xs sm:text-sm text-muted-foreground">
          安装后打开客户端，点击左下角头像 &rarr; <strong>设置</strong> &rarr; <strong>模型</strong>，点击「添加模型」。
        </p>
      </section>

      <!-- 2. 填写模型配置 -->
      <section class="space-y-4">
        <h2 class="text-lg font-semibold tracking-tight text-foreground border-b border-border/80 pb-2">
          2. 填写模型配置
        </h2>

        <div class="rounded-lg border border-border bg-card overflow-hidden">
          <table class="w-full text-xs sm:text-sm">
            <thead class="bg-muted/40 border-b border-border text-muted-foreground">
              <tr>
                <th class="p-3 text-left font-medium w-1/3">配置项</th>
                <th class="p-3 text-left font-medium">推荐填写值 / 说明</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <tr>
                <td class="p-3 font-semibold text-foreground">提供商</td>
                <td class="p-3 text-muted-foreground">选择 <strong>自定义 / Custom</strong></td>
              </tr>
              <tr>
                <td class="p-3 font-semibold text-foreground">接口地址</td>
                <td class="p-3 font-mono text-foreground"><code>https://api.724ai.org/v1</code></td>
              </tr>
              <tr>
                <td class="p-3 font-semibold text-foreground">API Key</td>
                <td class="p-3 text-muted-foreground">填写从 724AI 获取的专属密钥（<code class="font-mono">sk-...</code>）</td>
              </tr>
              <tr>
                <td class="p-3 font-semibold text-foreground">模型名称</td>
                <td class="p-3 font-mono text-foreground">例如 <code>gpt-5.6-sol</code> 或 <code>gpt-5.6-luna</code></td>
              </tr>
              <tr>
                <td class="p-3 font-semibold text-foreground">高级配置</td>
                <td class="p-3 text-muted-foreground">建议勾选「工具调用」与「图片输入」</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="p-4 rounded-lg border border-blue-500/30 bg-blue-500/10 text-xs sm:text-sm text-foreground">
          保存配置后，在聊天窗口顶部的模型下拉菜单中切换至您添加的自定义模型即可正常使用。
        </div>
      </section>
    </div>

    <!-- 底部返回与团队标识 -->
    <div class="pt-8 border-t border-border flex items-center justify-between text-xs text-muted-foreground">
      <span>© 2026 内部工作台 · 研发效能团队</span>
      <router-link to="/" class="hover:text-foreground hover:underline transition-colors">
        ← 返回首页导航
      </router-link>
    </div>
  </div>
</template>
