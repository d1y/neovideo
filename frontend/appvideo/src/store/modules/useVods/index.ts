import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { Category, VodItem } from '@/api/types'

export default defineStore("vods", () => {
  const map = ref(new Map<number, VodItem>)
  const currentCategory = ref<number>(-1)

  function setVodData(items: VodItem[]) {
    map.value.clear()
    items.forEach(item => {
      map.value.set(item.id, item)
    })
  }

  async function loadVodHomeDataWithApi(beforeCheck = true) {
    if (map.value.size) return
    // FIXME: remove this
  }

  const menus = computed<{
    id: number
    label: string
  }[]>(() => {
    return Array.from(map.value).map(item => {
      const data = item[1]
      return {
        id: data.id,
        label: data.name,
      }
    })
  })

  const category = computed<Map<number, Category[]>>(() => {
    const m = new Map<number, Category[]>()
    Array.from(map.value).forEach(item => {
      const id = item[0]
      const data = item[1]
      m.set(id, data.data.category)
    })
    return m
  })

  function getCategoryByID(id: number): Category[] {
    return category.value.get(id) || []
  }

  function setCurrentCategory(id: number) {
    currentCategory.value = id
  }

  return {
    map,
    menus,
    category,
    currentCategory,
    setCurrentCategory,
    getCategoryByID,
    loadVodHomeDataWithApi,
    setVodData,
  }
})