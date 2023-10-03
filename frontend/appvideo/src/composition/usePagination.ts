import { ListHeader } from '@/api/types'
import { Ref, computed, ComputedRef } from 'vue'
export default (page: Ref<number>, header: ComputedRef<ListHeader | undefined>) => {
  const isPrev = computed(() => {
    return page.value > 1
  })
  const isNext = computed(() => {
    const h = header.value
    if (!h) return false
    return page.value < h.page_count
  })

  const show = computed(() => {
    const h = header.value
    if (!h) return false
    return h.page_count > 1
  })

  const text = computed(() => {
    const h = header.value
    if (!h) return ""
    return `第${page.value}/${h.page_count}页(${h.record_count}条)`
  })

  return {
    isPrev,
    isNext,
    text,
    show,
  }
}