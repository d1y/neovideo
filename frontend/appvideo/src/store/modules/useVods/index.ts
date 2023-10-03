import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { Category, VodItem } from '@/api/types'

export default defineStore("vods", () => {
  const map = ref(new Map<number, VodItem>)

  function setVodData(items: VodItem[]) {
    map.value.clear()
    items.forEach(item => {
      map.value.set(item.id, item)
    })
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

  return {
    map,
    menus,
    setVodData,
  }
})