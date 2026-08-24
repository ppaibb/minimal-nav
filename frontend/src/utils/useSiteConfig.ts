import { ref } from 'vue'

export interface SiteConfig {
  site_name: string
  site_desc: string
}

const siteConfig = ref<SiteConfig>({
  site_name: 'Minimal Nav',
  site_desc: '统一汇聚团队核心工具、部署控制台、设计协作及文档中心，即时检索快速直达。'
})

const isLoaded = ref(false)

export const useSiteConfig = () => {
  const loadSiteConfig = async () => {
    try {
      const res = await fetch('/api/settings')
      if (res.ok) {
        const data = await res.json()
        if (data.code === 0 && data.data) {
          siteConfig.value = {
            site_name: data.data.site_name || 'Minimal Nav',
            site_desc: data.data.site_desc || ''
          }
          document.title = siteConfig.value.site_name
          isLoaded.value = true
        }
      }
    } catch {
      // 保持默认
    }
  }

  const updateSiteConfig = async (newConfig: Partial<SiteConfig>, token: string) => {
    const res = await fetch('/api/settings', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-Admin-Token': token
      },
      body: JSON.stringify(newConfig)
    })
    const data = await res.json()
    if (data.code === 0 && data.data) {
      siteConfig.value = {
        site_name: data.data.site_name || siteConfig.value.site_name,
        site_desc: data.data.site_desc || siteConfig.value.site_desc
      }
      document.title = siteConfig.value.site_name
      return { success: true }
    }
    return { success: false, msg: data.msg || '保存失败' }
  }

  return {
    siteConfig,
    isLoaded,
    loadSiteConfig,
    updateSiteConfig
  }
}
