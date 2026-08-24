import { Hono } from 'hono'
import { handle } from 'hono/cloudflare-pages'
import { cors } from 'hono/cors'

type Bindings = {
  DB: D1Database
  ADMIN_PASSWORD?: string
  JWT_SECRET?: string
}

const app = new Hono<{ Bindings: Bindings }>().basePath('/api')

// 1. CORS 跨域配置
app.use('*', cors({
  origin: '*',
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowHeaders: ['Content-Type', 'Authorization', 'X-Admin-Token', 'Accept', 'Origin'],
  exposeHeaders: ['Content-Length'],
  maxAge: 86400,
}))

// 统一成功/错误响应助手
const success = (c: any, data: any) => {
  return c.json({
    code: 0,
    msg: 'success',
    data,
  })
}

const error = (c: any, code: number, msg: string) => {
  const status = code >= 400 && code < 600 ? code : 200
  return c.json({
    code,
    msg,
    data: null,
  }, status as any)
}

// 2. 简易安全的 Web Crypto 签名 Token 工具 (无状态，完美支持多边缘节点)
async function createToken(secret: string, expireDays = 7): Promise<string> {
  const exp = Date.now() + expireDays * 24 * 60 * 60 * 1000
  const payload = JSON.stringify({ role: 'admin', exp })
  const enc = new TextEncoder()
  const key = await crypto.subtle.importKey(
    'raw',
    enc.encode(secret),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign']
  )
  const signature = await crypto.subtle.sign('HMAC', key, enc.encode(payload))
  const sigHex = Array.from(new Uint8Array(signature)).map(b => b.toString(16).padStart(2, '0')).join('')
  const b64Payload = btoa(payload)
  return `${b64Payload}.${sigHex}`
}

async function verifyToken(token: string, secret: string): Promise<boolean> {
  try {
    const parts = token.split('.')
    if (parts.length !== 2) return false
    const [b64Payload, sigHex] = parts
    const payloadStr = atob(b64Payload)
    const payload = JSON.parse(payloadStr)
    if (!payload.exp || Date.now() > payload.exp) return false

    const enc = new TextEncoder()
    const key = await crypto.subtle.importKey(
      'raw',
      enc.encode(secret),
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['verify']
    )
    const sigBytes = new Uint8Array(sigHex.match(/.{1,2}/g)!.map(byte => parseInt(byte, 16)))
    return await crypto.subtle.verify('HMAC', key, sigBytes, enc.encode(payloadStr))
  } catch {
    return false
  }
}

// 3. 管理员鉴权中间件
const authMiddleware = async (c: any, next: any) => {
  const token = c.req.header('X-Admin-Token') || c.req.query('token')
  if (!token) {
    return error(c, 401, '请先解锁管理权限')
  }

  const secret = c.env.ADMIN_PASSWORD || 'admin123'
  const isValid = await verifyToken(token, secret)
  if (!isValid) {
    return error(c, 401, '登录会话已失效，请重新输入口令')
  }
  await next()
}

// --- 路由模块 1: 健康检查 ---
app.get('/health', (c) => {
  return success(c, { status: 'up', platform: 'cloudflare-pages-d1' })
})

// --- 路由模块 2: 认证管理 ---
app.post('/auth/login', async (c) => {
  try {
    const body = await c.req.json()
    const password = body.password
    const adminPassword = c.env.ADMIN_PASSWORD || 'admin123'

    if (!password) {
      return error(c, 400, '请输入管理员口令')
    }

    if (password !== adminPassword) {
      return error(c, 401, '管理员口令错误')
    }

    const token = await createToken(adminPassword, 7)
    return success(c, {
      token,
      expires_in: 7 * 24 * 3600,
    })
  } catch (err: any) {
    return error(c, 400, '无效的请求参数: ' + err.message)
  }
})

app.get('/auth/check', authMiddleware, (c) => {
  return success(c, { authenticated: true })
})

// --- 路由模块 3: 导航链接管理 ---
app.get('/links', async (c) => {
  try {
    const { results } = await c.env.DB.prepare('SELECT * FROM links ORDER BY sort_order ASC, id ASC').all()
    return success(c, results || [])
  } catch (err: any) {
    return error(c, 500, '获取链接失败: ' + err.message)
  }
})

app.post('/links', authMiddleware, async (c) => {
  try {
    const body = await c.req.json()
    let { title, url, icon, category, sort_order, is_pinned } = body

    if (!title || !url) {
      return error(c, 400, '链接标题与 URL 不能为空')
    }

    if (!category) {
      category = 'Default'
    }

    // 若未指定图标，自动生成 Google 64px 高清 Favicon 抓取源
    if (!icon && url) {
      let targetUrl = url
      if (!targetUrl.startsWith('http://') && !targetUrl.startsWith('https://')) {
        targetUrl = 'https://' + targetUrl
      }
      icon = `https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=${encodeURIComponent(targetUrl)}&size=64`
    }

    const result = await c.env.DB.prepare(
      'INSERT INTO links (title, url, icon, category, sort_order, is_pinned) VALUES (?, ?, ?, ?, ?, ?)'
    ).bind(title, url, icon || '', category, sort_order || 0, is_pinned ? 1 : 0).run()

    const created = await c.env.DB.prepare('SELECT * FROM links WHERE id = ?').bind(result.meta.last_row_id).first()
    return success(c, created)
  } catch (err: any) {
    return error(c, 500, '创建链接失败: ' + err.message)
  }
})

app.put('/links/:id', authMiddleware, async (c) => {
  try {
    const id = parseInt(c.req.param('id'))
    if (isNaN(id)) return error(c, 400, '无效的 ID 参数')

    const body = await c.req.json()
    const { title, url, icon, category, sort_order, is_pinned } = body

    const existing = await c.env.DB.prepare('SELECT * FROM links WHERE id = ?').bind(id).first()
    if (!existing) return error(c, 404, '未找到该链接')

    await c.env.DB.prepare(
      'UPDATE links SET title = ?, url = ?, icon = ?, category = ?, sort_order = ?, is_pinned = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?'
    ).bind(
      title ?? existing.title,
      url ?? existing.url,
      icon ?? existing.icon,
      category ?? existing.category,
      sort_order ?? existing.sort_order,
      is_pinned !== undefined ? (is_pinned ? 1 : 0) : existing.is_pinned,
      id
    ).run()

    const updated = await c.env.DB.prepare('SELECT * FROM links WHERE id = ?').bind(id).first()
    return success(c, updated)
  } catch (err: any) {
    return error(c, 500, '更新链接失败: ' + err.message)
  }
})

app.delete('/links/:id', authMiddleware, async (c) => {
  try {
    const id = parseInt(c.req.param('id'))
    if (isNaN(id)) return error(c, 400, '无效的 ID 参数')

    const result = await c.env.DB.prepare('DELETE FROM links WHERE id = ?').bind(id).run()
    if (result.meta.changes === 0) {
      return error(c, 404, '未找到该链接')
    }

    return success(c, { deleted_id: id })
  } catch (err: any) {
    return error(c, 500, '删除链接失败: ' + err.message)
  }
})

// --- 路由模块 4: 公告管理 ---
app.get('/announcements', async (c) => {
  try {
    const { results } = await c.env.DB.prepare('SELECT * FROM announcements ORDER BY created_at DESC').all()
    // 兼容布尔类型
    const formatted = (results || []).map((item: any) => ({
      ...item,
      is_active: Boolean(item.is_active),
    }))
    return success(c, formatted)
  } catch (err: any) {
    return error(c, 500, '获取公告列表失败: ' + err.message)
  }
})

app.get('/announcements/active', async (c) => {
  try {
    const { results } = await c.env.DB.prepare('SELECT * FROM announcements WHERE is_active = 1 ORDER BY sort_order ASC, created_at DESC').all()
    const formatted = (results || []).map((item: any) => ({
      ...item,
      is_active: true,
    }))
    return success(c, formatted)
  } catch (err: any) {
    return error(c, 500, '获取生效公告失败: ' + err.message)
  }
})

app.get('/announcements/:id', async (c) => {
  try {
    const id = parseInt(c.req.param('id'))
    if (isNaN(id)) return error(c, 400, '无效的 ID 参数')

    const item: any = await c.env.DB.prepare('SELECT * FROM announcements WHERE id = ?').bind(id).first()
    if (!item) return error(c, 404, '未找到该公告详情')

    return success(c, {
      ...item,
      is_active: Boolean(item.is_active),
    })
  } catch (err: any) {
    return error(c, 500, '获取公告详情失败: ' + err.message)
  }
})

app.post('/announcements', authMiddleware, async (c) => {
  try {
    const body = await c.req.json()
    const { content, detail_md, is_active, sort_order } = body

    if (!content) return error(c, 400, '公告标题内容不能为空')

    const result = await c.env.DB.prepare(
      'INSERT INTO announcements (content, detail_md, is_active, sort_order) VALUES (?, ?, ?, ?)'
    ).bind(
      content,
      detail_md || '',
      is_active === false ? 0 : 1,
      sort_order || 0
    ).run()

    const created: any = await c.env.DB.prepare('SELECT * FROM announcements WHERE id = ?').bind(result.meta.last_row_id).first()
    return success(c, {
      ...created,
      is_active: Boolean(created.is_active),
    })
  } catch (err: any) {
    return error(c, 500, '发布公告失败: ' + err.message)
  }
})

app.put('/announcements/:id', authMiddleware, async (c) => {
  try {
    const id = parseInt(c.req.param('id'))
    if (isNaN(id)) return error(c, 400, '无效的 ID 参数')

    const body = await c.req.json()
    const { content, detail_md, is_active, sort_order } = body

    const existing: any = await c.env.DB.prepare('SELECT * FROM announcements WHERE id = ?').bind(id).first()
    if (!existing) return error(c, 404, '未找到该公告')

    await c.env.DB.prepare(
      'UPDATE announcements SET content = ?, detail_md = ?, is_active = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?'
    ).bind(
      content ?? existing.content,
      detail_md ?? existing.detail_md,
      is_active !== undefined ? (is_active ? 1 : 0) : existing.is_active,
      sort_order ?? existing.sort_order,
      id
    ).run()

    const updated: any = await c.env.DB.prepare('SELECT * FROM announcements WHERE id = ?').bind(id).first()
    return success(c, {
      ...updated,
      is_active: Boolean(updated.is_active),
    })
  } catch (err: any) {
    return error(c, 500, '更新公告失败: ' + err.message)
  }
})

app.put('/announcements/:id/toggle', authMiddleware, async (c) => {
  try {
    const id = parseInt(c.req.param('id'))
    if (isNaN(id)) return error(c, 400, '无效的 ID 参数')

    const existing: any = await c.env.DB.prepare('SELECT * FROM announcements WHERE id = ?').bind(id).first()
    if (!existing) return error(c, 404, '未找到该公告')

    const newActive = existing.is_active ? 0 : 1
    await c.env.DB.prepare('UPDATE announcements SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?').bind(newActive, id).run()

    return success(c, {
      ...existing,
      is_active: Boolean(newActive),
    })
  } catch (err: any) {
    return error(c, 500, '切换状态失败: ' + err.message)
  }
})

app.delete('/announcements/:id', authMiddleware, async (c) => {
  try {
    const id = parseInt(c.req.param('id'))
    if (isNaN(id)) return error(c, 400, '无效的 ID 参数')

    const result = await c.env.DB.prepare('DELETE FROM announcements WHERE id = ?').bind(id).run()
    if (result.meta.changes === 0) return error(c, 404, '未找到该公告')

    return success(c, { deleted_id: id })
  } catch (err: any) {
    return error(c, 500, '删除公告失败: ' + err.message)
  }
})

// --- 路由模块 5: 系统设置 ---
app.get('/settings', async (c) => {
  try {
    const { results } = await c.env.DB.prepare('SELECT * FROM settings').all()
    const settingsMap: Record<string, string> = {
      site_name: 'Minimal Nav',
      site_desc: '企业极简导航系统 · Cloudflare Serverless 版',
    }

    if (results) {
      for (const item of results as any[]) {
        settingsMap[item.key] = item.value
      }
    }

    return success(c, settingsMap)
  } catch (err: any) {
    return error(c, 500, '获取系统设置失败: ' + err.message)
  }
})

app.put('/settings', authMiddleware, async (c) => {
  try {
    const input = await c.req.json()
    if (typeof input !== 'object' || input === null) {
      return error(c, 400, '请求参数格式错误')
    }

    const statements = []
    for (const [key, value] of Object.entries(input)) {
      statements.push(
        c.env.DB.prepare(
          'INSERT INTO settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP'
        ).bind(key, String(value ?? ''))
      )
    }

    if (statements.length > 0) {
      await c.env.DB.batch(statements)
    }

    // 重新获取全部设置返回
    const { results } = await c.env.DB.prepare('SELECT * FROM settings').all()
    const settingsMap: Record<string, string> = {
      site_name: 'Minimal Nav',
      site_desc: '',
    }
    if (results) {
      for (const item of results as any[]) {
        settingsMap[item.key] = item.value
      }
    }

    return success(c, settingsMap)
  } catch (err: any) {
    return error(c, 500, '更新设置失败: ' + err.message)
  }
})

// --- 路由模块 6: 边缘网络探测工具 ---
app.get('/tools/favicon', (c) => {
  const targetUrl = c.req.query('url')
  if (!targetUrl) return error(c, 400, 'URL 参数不能为空')

  let fullUrl = targetUrl
  if (!fullUrl.startsWith('http://') && !fullUrl.startsWith('https://')) {
    fullUrl = 'https://' + fullUrl
  }

  try {
    const u = new URL(fullUrl)
    const domain = u.hostname
    const googleFavicon = `https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=${encodeURIComponent(fullUrl)}&size=64`
    const directFavicon = `${u.protocol}//${u.host}/favicon.ico`

    return success(c, {
      domain,
      favicon: googleFavicon,
      direct: directFavicon,
    })
  } catch {
    return error(c, 400, '无效的 URL 格式')
  }
})

app.post('/tools/ping', async (c) => {
  try {
    const body = await c.req.json()
    const reqUrl = body.url
    if (!reqUrl) return error(c, 400, '请提供要探测的 URL')

    let targetUrl = reqUrl
    if (!targetUrl.startsWith('http://') && !targetUrl.startsWith('https://')) {
      targetUrl = 'https://' + targetUrl
    }

    const startTime = Date.now()
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 4000)

    try {
      const resp = await fetch(targetUrl, {
        method: 'GET',
        signal: controller.signal,
        headers: {
          'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 MinimalNavPing/1.0',
        },
      })
      clearTimeout(timeoutId)
      const latency = Date.now() - startTime
      const isHealthy = resp.status >= 200 && resp.status < 400

      return success(c, {
        url: reqUrl,
        healthy: isHealthy,
        status_code: resp.status,
        latency_ms: latency,
        error: '',
      })
    } catch (fetchErr: any) {
      clearTimeout(timeoutId)
      const latency = Date.now() - startTime
      return success(c, {
        url: reqUrl,
        healthy: false,
        status_code: 0,
        latency_ms: latency,
        error: fetchErr.name === 'AbortError' ? '请求超时 (4s)' : fetchErr.message,
      })
    }
  } catch (err: any) {
    return error(c, 400, '请求错误: ' + err.message)
  }
})

// --- 路由模块 7: 备份与书签导入 ---
app.get('/backup/export', authMiddleware, async (c) => {
  try {
    const { results: links } = await c.env.DB.prepare('SELECT * FROM links').all()
    const { results: announcements } = await c.env.DB.prepare('SELECT * FROM announcements').all()

    const backupData = {
      version: '1.0',
      exported_at: new Date().toISOString(),
      links: links || [],
      announcements: (announcements || []).map((a: any) => ({
        ...a,
        is_active: Boolean(a.is_active),
      })),
    }

    c.header('Content-Disposition', 'attachment; filename=minimal-nav-backup.json')
    c.header('Content-Type', 'application/json')
    return c.body(JSON.stringify(backupData, null, 2))
  } catch (err: any) {
    return error(c, 500, '导出备份失败: ' + err.message)
  }
})

app.post('/backup/import', authMiddleware, async (c) => {
  try {
    const body = await c.req.json()
    const { mode, links, announcements } = body

    const statements = []

    if (mode === 'overwrite') {
      statements.push(c.env.DB.prepare('DELETE FROM links'))
      statements.push(c.env.DB.prepare('DELETE FROM announcements'))
    }

    if (Array.isArray(links)) {
      for (const l of links) {
        statements.push(
          c.env.DB.prepare(
            'INSERT INTO links (title, url, category, icon, sort_order, is_pinned) VALUES (?, ?, ?, ?, ?, ?)'
          ).bind(
            l.title || '',
            l.url || '',
            l.category || 'Default',
            l.icon || '',
            l.sort_order || 0,
            l.is_pinned ? 1 : 0
          )
        )
      }
    }

    if (Array.isArray(announcements)) {
      for (const a of announcements) {
        statements.push(
          c.env.DB.prepare(
            'INSERT INTO announcements (content, detail_md, is_active, sort_order) VALUES (?, ?, ?, ?)'
          ).bind(
            a.content || '',
            a.detail_md || '',
            a.is_active ? 1 : 0,
            a.sort_order || 0
          )
        )
      }
    }

    if (statements.length > 0) {
      await c.env.DB.batch(statements)
    }

    return success(c, {
      imported_links: links ? links.length : 0,
      imported_announcements: announcements ? announcements.length : 0,
    })
  } catch (err: any) {
    return error(c, 500, '导入备份失败: ' + err.message)
  }
})

// 解析 Netscape Bookmark HTML 工具函数
function parseNetscapeBookmarks(htmlContent: string) {
  const items: Array<{ title: string; url: string; category: string }> = []
  const lines = htmlContent.split('\n')

  let currentCategory = '常用书签'
  const categoryStack = [currentCategory]

  const h3Regex = /<H3[^>]*>(.*?)<\/H3>/i
  const aRegex = /<A\s+[^>]*HREF="([^"]+)"[^>]*>(.*?)<\/A>/i

  for (const line of lines) {
    const trimmed = line.trim()

    if (trimmed.toUpperCase().includes('<DL>')) {
      // 进入下一层
    } else if (trimmed.toUpperCase().includes('</DL>')) {
      if (categoryStack.length > 1) {
        categoryStack.pop()
        currentCategory = categoryStack[categoryStack.length - 1]
      }
    }

    const h3Match = trimmed.match(h3Regex)
    if (h3Match && h3Match[1]) {
      const catName = h3Match[1].trim()
      if (catName && catName.toLowerCase() !== 'bookmarks bar' && catName !== '书签栏') {
        currentCategory = catName
        categoryStack.push(currentCategory)
      }
    }

    const aMatch = trimmed.match(aRegex)
    if (aMatch && aMatch[1]) {
      const linkUrl = aMatch[1].trim()
      let linkTitle = (aMatch[2] || '').trim()

      if (linkUrl.startsWith('http://') || linkUrl.startsWith('https://')) {
        if (!linkTitle) linkTitle = linkUrl
        items.push({
          title: linkTitle,
          url: linkUrl,
          category: currentCategory,
        })
      }
    }
  }

  return items
}

app.post('/backup/import-bookmarks', authMiddleware, async (c) => {
  try {
    const formData = await c.req.formData()
    const file = formData.get('file') as File | null
    if (!file) {
      return error(c, 400, '请上传有效的 HTML 书签文件')
    }

    const content = await file.text()
    const bookmarks = parseNetscapeBookmarks(content)

    if (bookmarks.length === 0) {
      return error(c, 400, '未能从文件中解析出有效的网址书签')
    }

    // 批量分块插入 D1 (D1 batch 一次支持百条语句)
    const chunkSize = 50
    for (let i = 0; i < bookmarks.length; i += chunkSize) {
      const chunk = bookmarks.slice(i, i + chunkSize)
      const stmts = chunk.map((bm) => {
        const favicon = `https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=${encodeURIComponent(bm.url)}&size=64`
        return c.env.DB.prepare(
          'INSERT INTO links (title, url, category, icon) VALUES (?, ?, ?, ?)'
        ).bind(bm.title, bm.url, bm.category, favicon)
      })
      await c.env.DB.batch(stmts)
    }

    return success(c, {
      imported_count: bookmarks.length,
    })
  } catch (err: any) {
    return error(c, 500, '导入书签失败: ' + err.message)
  }
})

// 导出 Cloudflare Pages 处理器
export const onRequest = handle(app)
export default app
