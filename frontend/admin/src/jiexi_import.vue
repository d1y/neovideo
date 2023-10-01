<template>
  <el-input v-model="val" type="textarea" :rows="12" placeholder="导入信息" />
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  data: string,
}>()

const val = ref(props.data)

defineExpose({
  getValue() {
    return val.value
  }
})
</script>