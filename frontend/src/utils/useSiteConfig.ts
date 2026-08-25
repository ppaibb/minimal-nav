import { ref } from 'vue'

export interface SiteConfig {
  site_name: string
  site_desc: string
  icp_beian?: string
}

const siteConfig = ref<SiteConfig>({
  site_name: '团队常用工具推荐',
  site_desc: '',
  icp_beian: ''
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
            site_name: data.data.site_name || '团队常用工具推荐',
            site_desc: data.data.site_desc || '',
            icp_beian: data.data.icp_beian || ''
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
        site_desc: data.data.site_desc !== undefined ? data.data.site_desc : '',
        icp_beian: data.data.icp_beian !== undefined ? data.data.icp_beian : ''
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
